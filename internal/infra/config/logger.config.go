package config

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	l "github.com/mauriciorobertodev/whappy-go/internal/app/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/logger"
)

func LoadLoggers(level l.Level) {
	app.RegisterLogger(app.LogKeyService, logger.NewCuteLogger("SERVICE", level))
	app.RegisterLogger(app.LogKeyInstanceService, logger.NewCuteLogger("INSTANCE SERVICE", level))
	app.RegisterLogger(app.LogKeyFileService, logger.NewCuteLogger("FILE SERVICE", level))
	app.RegisterLogger(app.LogKeyMessageService, logger.NewCuteLogger("MESSAGE SERVICE", level))
	app.RegisterLogger(app.LogKeyChatService, logger.NewCuteLogger("CHAT SERVICE", level))
	app.RegisterLogger(app.LogKeyContactService, logger.NewCuteLogger("CONTACT SERVICE", level))
	app.RegisterLogger(app.LogKeyGroupService, logger.NewCuteLogger("GROUP SERVICE", level))
	app.RegisterLogger(app.LogKeyPictureService, logger.NewCuteLogger("PICTURE SERVICE", level))
	app.RegisterLogger(app.LogKeyUploadService, logger.NewCuteLogger("UPLOAD SERVICE", level))
	app.RegisterLogger(app.LogKeyBlocklistService, logger.NewCuteLogger("BLOCKLIST SERVICE", level))
	app.RegisterLogger(app.LogKeyTokenService, logger.NewCuteLogger("TOKEN SERVICE", level))
	app.RegisterLogger(app.LogKeyWhatsapp, logger.NewCuteLogger("WHATSAPP", level))
	app.RegisterLogger(app.LogKeyDatabase, logger.NewCuteLogger("DATABASE", level))
	app.RegisterLogger(app.LogKeyCache, logger.NewCuteLogger("CACHE", level))
	app.RegisterLogger(app.LogKeyMiddleware, logger.NewCuteLogger("MIDDLEWARE", level))
	app.RegisterLogger(app.LogKeyEventBus, logger.NewCuteLogger("EVENT BUS", level))
	app.RegisterLogger(app.LogKeyMigrator, logger.NewCuteLogger("MIGRATOR", level))
}
