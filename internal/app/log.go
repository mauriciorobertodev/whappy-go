package app

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
)

var loggers = make(map[string]logger.Logger)

func RegisterLogger(name string, l logger.Logger) {
	loggers[name] = l
}

func GetLogger(name string) logger.Logger {
	if l, ok := loggers[name]; ok {
		return l
	}
	return logger.NewDefaultLogger("", logger.LevelNone)
}

const (
	LogKeyCache            = "cache"
	LogKeyService          = "service"
	LogKeyInstanceService  = "instance_service"
	LogKeyFileService      = "file_service"
	LogKeyMessageService   = "message_service"
	LogKeyChatService      = "chat_service"
	LogKeyContactService   = "contact_service"
	LogKeyGroupService     = "group_service"
	LogKeyPictureService   = "picture_service"
	LogKeyUploadService    = "upload_service"
	LogKeyBlocklistService = "blocklist_service"
	LogKeyTokenService     = "token_service"
	LogKeyWhatsapp         = "whatsapp"
	LogKeyDatabase         = "database"
	LogKeyMiddleware       = "middleware"
	LogKeyEventBus         = "eventbus"
	LogKeyMigrator         = "migrator"
	LogKeyCacheService     = "cache_service"
)

func GetCacheLogger() logger.Logger {
	return GetLogger(LogKeyCache)
}

func GetServiceLogger() logger.Logger {
	return GetLogger(LogKeyService)
}

func GetWhatsappLogger() logger.Logger {
	return GetLogger(LogKeyWhatsapp)
}

func GetDatabaseLogger() logger.Logger {
	return GetLogger(LogKeyDatabase)
}

func GetMiddlewareLogger() logger.Logger {
	return GetLogger(LogKeyMiddleware)
}

func GetEventBusLogger() logger.Logger {
	return GetLogger(LogKeyEventBus)
}

func GetFileServiceLogger() logger.Logger {
	return GetLogger(LogKeyFileService)
}

func GetMessageServiceLogger() logger.Logger {
	return GetLogger(LogKeyMessageService)
}

func GetChatServiceLogger() logger.Logger {
	return GetLogger(LogKeyChatService)
}

func GetContactServiceLogger() logger.Logger {
	return GetLogger(LogKeyContactService)
}

func GetGroupServiceLogger() logger.Logger {
	return GetLogger(LogKeyGroupService)
}

func GetInstanceServiceLogger() logger.Logger {
	return GetLogger(LogKeyInstanceService)
}

func GetPictureServiceLogger() logger.Logger {
	return GetLogger(LogKeyPictureService)
}

func GetUploadServiceLogger() logger.Logger {
	return GetLogger(LogKeyUploadService)
}

func GetBlocklistServiceLogger() logger.Logger {
	return GetLogger(LogKeyBlocklistService)
}

func GetTokenServiceLogger() logger.Logger {
	return GetLogger(LogKeyTokenService)
}

func GetMigratorLogger() logger.Logger {
	return GetLogger(LogKeyMigrator)
}
