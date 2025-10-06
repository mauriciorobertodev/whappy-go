package main

import (
	"context"

	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/consumer"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/database"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/eventbus"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/registry"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/token"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/whatsapp/meow"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/handler"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env")
	}

	// Startup
	l := logger.NewCuteLogger("MAIN", "info")
	l.Info("🚀 Starting Whappy Go API v" + config.APP_VERSION)

	// Bootstrapping
	l.Info("⚙️  Loading app configuration...")
	appConfig := config.LoadAppConfig()

	// Logger
	config.LoadLoggers(appConfig.LOG_LEVEL)

	// Database
	l.Info("🗄️  Connecting to database...")
	whappyDB := database.New(config.LoadWhappyDatabaseConfig())

	// Migrations
	l.Info("🛠️  Running database migrations...")
	database.NewMigrator(whappyDB, whappyDB.DriverName()).Up()

	// Registries
	l.Info("📋 Setting up registries...")
	instRegistry := registry.NewInMemoryInstanceRegistry()

	// Events
	l.Info("📦 Setting up event bus...")
	bus := eventbus.New(config.LoadEventBusConfig())

	// Consumers
	l.Info("🍿 Setting up event consumers...")
	if appConfig.IsDevelopment() {
		debugConsumer := consumer.NewDevConsumer()
		bus.SubscribeAll(debugConsumer.Handler)
	}

	// Cache
	l.Info("🧠 Setting up cache...")
	cache := cache.New(config.LoadCacheConfig())

	// Storage
	l.Info("🗃️  Setting up storage...")
	storageConfig := config.LoadStorageConfig()
	storage := storage.New(storageConfig)

	// Whatsmeow
	l.Info("😻 Setting up whatsmeow...")
	whatsapp := meow.New(ctx, config.LoadWhatsmeowDatabaseConfig(), storage, bus, cache)

	// Token
	l.Info("🔑 Setting up token...")
	generator := token.NewGenerator()
	hasher := token.NewHasher(config.LoadTokenConfig())

	// Repositories
	l.Info("📚 Setting up repositories...")
	instRepo := repository.NewInstanceRepository(whappyDB)
	tokenRepo := repository.NewTokenRepository(whappyDB)
	fileRepo := repository.NewFileRepository(whappyDB)
	webhookRepo := repository.NewWebhookRepository(whappyDB)

	// Services / Use Cases
	l.Info("🔧 Setting up services...")
	tokenService := service.NewTokenService(tokenRepo, hasher, generator, bus, cache)
	webhookService := service.NewWebhookService(webhookRepo, bus, appConfig.MAX_WEBHOOKS)
	instService := service.NewInstanceService(tokenService, instRepo, instRegistry, bus)
	sessionService := service.NewSessionService(instRepo, whatsapp, bus)
	fileService := service.NewFileService(storage, fileRepo)
	messageService := service.NewMessageService(whatsapp, storage, fileService, cache, appConfig.CACHE_FILE_UPLOAD_TTL)
	chatService := service.NewChatService(whatsapp)
	contactService := service.NewContactService(whatsapp)
	groupService := service.NewGroupService(whatsapp, bus, fileService)
	pictureService := service.NewPictureService(whatsapp)
	uploadService := service.NewUploadService(fileService, fileRepo, storage, bus)
	blocklistService := service.NewBlocklistService(whatsapp, bus)

	// Middleware
	l.Info("🛡️  Setting up middleware...")
	authMiddleware := middleware.NewAuthMiddleware(appConfig.ADMIN_TOKEN, tokenService)
	instMiddleware := middleware.NewInstanceMiddleware(instRegistry, instRepo, sessionService)

	// Handlers
	l.Info("🖥️  Setting up HTTP handlers...")
	healthHandler := handler.NewHealthHandler()
	instHandler := handler.NewInstanceHandler(instService, instRegistry)
	sessionHandler := handler.NewSessionHandler(sessionService)
	messageHandler := handler.NewMessageHandler(messageService)
	chatHandler := handler.NewChatHandler(chatService)
	contactHandler := handler.NewContactHandler(contactService)
	groupHandler := handler.NewGroupHandler(groupService, bus)
	pictureHandler := handler.NewPictureHandler(pictureService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	blocklistHandler := handler.NewBlocklistHandler(blocklistService)
	webhookHandler := handler.NewWebhookHandler(webhookService)

	// Router
	l.Info("🛣️  Setting up HTTP routes...")
	r := http.NewFiberApp(config.APP_NAME, config.APP_VERSION, appConfig.IsProduction())

	healthHandler.RegisterRoutes(r)
	instHandler.RegisterRoutes(r, authMiddleware)
	sessionHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	messageHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	chatHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	contactHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	groupHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	pictureHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	uploadHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	blocklistHandler.RegisterRoutes(r, authMiddleware, instMiddleware)
	webhookHandler.RegisterRoutes(r, authMiddleware, instMiddleware)

	if storageConfig.IsLocal() {
		r.Get("/storage/*", static.New(storageConfig.Path))
	}

	l.Info("🔥 Server is running on port " + appConfig.APP_PORT)
	panic(r.Listen(":" + appConfig.APP_PORT))
}
