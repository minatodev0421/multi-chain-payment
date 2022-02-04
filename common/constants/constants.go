package constants

const (
	DEFAULT_SELECT_LIMIT    = "100"
	PAGE_SIZE_DEFAULT_VALUE = "10"

	URL_EVENT_PREFIX   = "events"
	URL_BILLING_PREFIX = "billing"
	URL_STORAGE_PREFIX = "storage"

	HTTP_STATUS_SUCCESS = "success"
	HTTP_STATUS_FAIL    = "fail"
	HTTP_STATUS_ERROR   = "error"

	HTTP_CODE_200_OK                    = "200" //http.StatusOk
	HTTP_CODE_400_BAD_REQUEST           = "400" //http.StatusBadRequest
	HTTP_CODE_401_UNAUTHORIZED          = "401" //http.StatusUnauthorized
	HTTP_CODE_500_INTERNAL_SERVER_ERROR = "500" //http.StatusInternalServerError
	HTTP_REQUEST_HEADER_AUTHRORIZATION  = "Authorization"

	NETWORK_NAME_POLYGON = "polygon"
	COIN_NAME_USDC       = "USDC"
	COIN_ADDRESS_USDC    = "0xe11A86849d99F524cAC3E7A0Ec1241828e332C62"

	TRANSACTION_STATUS_SUCCESS = "success"
	TRANSACTION_STATUS_FAIL    = "fail"

	URL_HOST_GET_COMMON      = "/common"
	URL_HOST_GET_HOST_INFO   = "/host/info"
	URL_SYSTEM_CONFIG_PARAMS = "/system/params"

	SOURCE_ID_OF_PAYMENT = 4

	PROCESS_STATUS_WAITING_PAYMENT     = "Pending" //wait for payment
	PROCESS_STATUS_PAID                = "Paid"
	PROCESS_STATUS_PROCESSING          = "Processing"
	PROCESS_STATUS_TASK_CREATED        = "TaskCreated"
	PROCESS_STATUS_DEAL_SENT           = "DealSent"
	PROCESS_STATUS_DEAL_SENT_FAILED    = "DealSentFailed"
	PROCESS_STATUS_DEAL_SEND_CANCELLED = "DealSendCancelled"
	PROCESS_STATUS_UNLOCK_REFUNDED     = "UnlockRefundSucceeded"
	PROCESS_STATUS_UNLOCK_REFUNDFAILED = "UnlockRefundFailed"

	DEAL_STATUS_ACTIVE = "StorageDealActive"
	DEAL_STATUS_ERROR  = "StorageDealError"

	IPFS_URL_PREFIX_BEFORE_HASH = "/ipfs/"
	IPFS_File_PINNED_STATUS     = "Pinned"

	LOTUS_TASK_TYPE_VERIFIED = "verified"
	LOTUS_TASK_TYPE_REGULAR  = "regular"

	LOTUS_PRICE_MULTIPLE_1E18 = 1e18 // 10^18
	FILE_BLOCK_NUMBER_MAX     = 999999999999999
	TIME_HALF_AN_HOUR         = 30 * 60 * 1000

	SOURCE_FILE_STATUS_CREATED      = "Created"
	SOURCE_FILE_STATUS_TASK_CREATED = "TaskCreated"

	OFFLINE_DEAL_UNLOCK_STATUS_NOT_UNLOCKED  = "NotUnlocked"
	OFFLINE_DEAL_UNLOCK_STATUS_UNLOCKED      = "Unlocked"
	OFFLINE_DEAL_UNLOCK_STATUS_UNLOCK_FAILED = "UnlockFailed"

	SIGNATURE_DEFAULT_VALUE = "0" //init value,no unlock operation has been performed
	SIGNATURE_SUCCESS_VALUE = "1" //init value,no unlock operation has been performed
	SIGNATURE_FAILED_VALUE  = "2" //init value,no unlock operation has been performed

	DURATION_DAYS_DEFAULT = 525

	SOURCE_FILE_TYPE_NORMAL = 0

	SOURCE_FILE_UPLOAD_HISTORY_STATUS_CREATED = "Created"
	SOURCE_FILE_UPLOAD_HISTORY_STATUS_DELETED = "Deleted"

	BYTES_1GB     = 1024 * 1024 * 1024
	EPOCH_PER_DAY = 24 * 60 * 2
)
