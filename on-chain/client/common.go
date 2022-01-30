package client

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"payment-bridge/blockchain/browsersync/scanlockpayment/polygon"
	"payment-bridge/common/constants"
	"payment-bridge/common/utils"
	"payment-bridge/config"
	"payment-bridge/database"
	"payment-bridge/models"
	"payment-bridge/on-chain/goBind"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/logs"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetEthClient() (*ethclient.Client, error) {
	polygonRpcUrl := config.GetConfig().PolygonRpcUrl
	ethClient, err := ethclient.Dial(polygonRpcUrl)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return ethClient, nil
}

func GetSwanPaymentTransactor(ethClient *ethclient.Client) (*common.Address, *goBind.SwanPaymentTransactor, error) {
	recipient := common.HexToAddress(polygon.GetConfig().PolygonMainnetNode.PaymentContractAddress)
	swanPaymentTransactor, err := goBind.NewSwanPaymentTransactor(recipient, ethClient)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	return &recipient, swanPaymentTransactor, nil
}

func GetFilswanOracleSession(ethClient *ethclient.Client) (*goBind.FilswanOracleSession, error) {
	daoAddress := common.HexToAddress(polygon.GetConfig().PolygonMainnetNode.DaoSwanOracleAddress)
	daoOracleContractInstance, err := goBind.NewFilswanOracle(daoAddress, ethClient)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	filswanOracleSession := &goBind.FilswanOracleSession{}
	filswanOracleSession.Contract = daoOracleContractInstance

	return filswanOracleSession, nil
}

func GetTransactOpts(ethClient *ethclient.Client) (*bind.TransactOpts, error) {
	privateKeyOnPolygon := os.Getenv("privateKeyOnPolygon")
	if len(privateKeyOnPolygon) <= 0 {
		err := fmt.Errorf("env variable privateKeyOnPolygon is not defined")
		logs.GetLogger().Fatal(err)
	}

	if strings.HasPrefix(strings.ToLower(privateKeyOnPolygon), "0x") {
		privateKeyOnPolygon = privateKeyOnPolygon[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyOnPolygon)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	adminAddress := common.HexToAddress(config.GetConfig().AdminWalletOnPolygon)
	nonce, err := ethClient.PendingNonceAt(context.Background(), adminAddress)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	chainId, err := ethClient.ChainID(context.Background())
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	transactOpts.Nonce = big.NewInt(int64(nonce))
	transactOpts.GasPrice = gasPrice
	transactOpts.GasLimit = uint64(polygon.GetConfig().PolygonMainnetNode.GasLimit)
	transactOpts.Context = context.Background()

	return transactOpts, nil
}

func GetPaymentSession() (*goBind.SwanPaymentSession, error) {
	ethClient, err := GetEthClient()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanPaymentSession := &goBind.SwanPaymentSession{}

	paymentGatewayAddress := common.HexToAddress(polygon.GetConfig().PolygonMainnetNode.PaymentContractAddress)

	paymentGatewayInstance, err := goBind.NewSwanPayment(paymentGatewayAddress, ethClient)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanPaymentSession.Contract = paymentGatewayInstance

	return swanPaymentSession, nil
}

func GetContractAbi() (*abi.ABI, error) {
	paymentAbiString := goBind.SwanPaymentABI

	contractAbi, err := abi.JSON(strings.NewReader(paymentAbiString))
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &contractAbi, nil
}

func GetPaymentInfo(srcFilePayloadCid string) (*models.EventLockPayment, error) {
	swanPaymentSession, err := GetPaymentSession()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	paymentInfo, err := swanPaymentSession.GetLockedPaymentInfo(srcFilePayloadCid)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	var event *models.EventLockPayment
	if paymentInfo.IsExisted {
		event = new(models.EventLockPayment)
		event.TokenAddress = paymentInfo.Token.Hex()
		event.AddressFrom = paymentInfo.Owner.String()
		event.AddressTo = paymentInfo.Recipient.String()
		event.LockedFee = paymentInfo.LockedFee.String()
		event.PayloadCid = srcFilePayloadCid
		event.LockPaymentTime = strconv.FormatInt(utils.GetCurrentUtcMilliSecond(), 10)
		usdcCoinId, err := models.FindCoinIdByUUID(constants.COIN_TYPE_USDC_ON_POLYGON_UUID)
		if err != nil {
			logs.GetLogger().Error(err)
		} else {
			event.CoinId = usdcCoinId
		}

		networkId, err := models.FindNetworkIdByUUID(constants.NETWORK_TYPE_POLYGON_UUID)
		if err != nil {
			logs.GetLogger().Error(err)
		} else {
			event.NetworkId = networkId
		}

		err = database.SaveOne(event)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}

		dealFiles, err := models.GetDealFileBySourceFilePayloadCid(srcFilePayloadCid)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}

		if len(dealFiles) > 0 {
			dealFileId := dealFiles[0].ID

			err := models.UpdateDealFileLockPaymentStatus(dealFileId, constants.LOCK_PAYMENT_STATUS_PROCESSING)
			if err != nil {
				logs.GetLogger().Error(err)
				return nil, err
			}
		}
	}

	return event, nil
}
