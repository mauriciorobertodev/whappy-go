package app

type AppCode string

const (
	CodeUnknown          AppCode = "UNKNOWN"
	CodeSuccess          AppCode = "SUCCESS"
	CodeInvalidJSON      AppCode = "INVALID_JSON"
	CodeValidationFailed AppCode = "VALIDATION_FAILED"
	CodeInternalError    AppCode = "INTERNAL_ERROR"
	CodeInstanceNotFound AppCode = "INSTANCE_NOT_FOUND"
	CodeMissingData      AppCode = "MISSING_DATA"
	CodeInvalidCursor    AppCode = "INVALID_CURSOR"
	CodeInvalidToken     AppCode = "INVALID_TOKEN"
	CodeNotAdmin         AppCode = "NOT_ADMIN"

	CodeInstanceBanned              AppCode = "INSTANCE_BANNED"
	CodeInstanceAlreadyLoggedIn     AppCode = "INSTANCE_ALREADY_LOGGED_IN"
	CodeInstanceIsPairing           AppCode = "INSTANCE_IS_PAIRING"
	CodeInstanceAlreadyConnected    AppCode = "INSTANCE_ALREADY_CONNECTED"
	CodeInstanceIsConnecting        AppCode = "INSTANCE_IS_CONNECTING"
	CodeInstanceNotLoggedIn         AppCode = "INSTANCE_NOT_LOGGED_IN"
	CodeInstanceNotConnected        AppCode = "INSTANCE_NOT_CONNECTED"
	CodeInstanceNotPairing          AppCode = "INSTANCE_NOT_PAIRING"
	CodeInstanceAlreadyLoggedOut    AppCode = "INSTANCE_ALREADY_LOGGED_OUT"
	CodeInstanceAlreadyDisconnected AppCode = "INSTANCE_ALREADY_DISCONNECTED"
	CodeInstanceNoQrCode            AppCode = "INSTANCE_NO_QR_CODE"
	CodeInstanceAlreadyPaired       AppCode = "INSTANCE_ALREADY_PAIRED"

	CodeFailPairingClientOutdated     AppCode = "FAIL_PAIRING_CLIENT_OUTDATED"
	CodeFailPairingWithoutMultidevice AppCode = "FAIL_PAIRING_WITHOUT_MULTIDEVICE"
	CodeFailPairingTimeout            AppCode = "FAIL_PAIRING_TIMEOUT"
	CodeFailWhatsappUpload            AppCode = "FAIL_WHATSAPP_UPLOAD"

	CodeInvalidThumbnail AppCode = "INVALID_THUMBNAIL"
	CodeInvalidImage     AppCode = "INVALID_IMAGE"
	CodeInvalidVideo     AppCode = "INVALID_VIDEO"
	CodeInvalidAudio     AppCode = "INVALID_AUDIO"
	CodeInvalidVoice     AppCode = "INVALID_VOICE"

	CodeMediaUnreachable AppCode = "MEDIA_UNREACHABLE"
	CodeMediaCorrupted   AppCode = "MEDIA_CORRUPTED"

	CodeDatabaseError AppCode = "DATABASE_ERROR"

	CodeFileNotFound        AppCode = "FILE_NOT_FOUND"
	CodeInvalidFile         AppCode = "INVALID_FILE"
	CodeFileTooLarge        AppCode = "FILE_TOO_LARGE"
	CodeUnsupportedFileType AppCode = "UNSUPPORTED_FILE_TYPE"
	CodeCorruptedFile       AppCode = "CORRUPTED_FILE"
	CodeFileUploadFailed    AppCode = "FILE_UPLOAD_FAILED"
	CodeFileDownloadFailed  AppCode = "FILE_DOWNLOAD_FAILED"
	CodeFileStorageFailed   AppCode = "FILE_STORAGE_FAILED"
	CodeFileUnreachable     AppCode = "FILE_UNREACHABLE"
	CodeFileCannotBeDeleted AppCode = "FILE_CANNOT_BE_DELETED"
	CodeFileSourceEmpty     AppCode = "FILE_SOURCE_EMPTY"

	GLOBAL_STORAGE_UNAVAILABLE AppCode = "GLOBAL_STORAGE_UNAVAILABLE"
)
