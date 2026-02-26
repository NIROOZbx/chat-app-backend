package bootstrap

import (
	"chat-app/internal/handlers"
	"chat-app/internal/repositories"
	"chat-app/internal/routes"
	"chat-app/internal/services"
	"chat-app/internal/shared/cache"
	"chat-app/internal/shared/config"
	"chat-app/internal/shared/database"
	"chat-app/internal/shared/hub"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/middleware"
	"chat-app/internal/shared/pubsub"
	"chat-app/internal/shared/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type App struct {
	DB     *sqlx.DB
	Server *http.Server
	Redis  *redis.Client
	Log    *logger.Logger
}

func StartUp(cfg *config.Config) *App {

	db := database.ConnectSupabaseDB(cfg.DB)

	client := cache.SetupRedis(&cfg.Redis)

	l := logger.NewLogger(0, "logs/", 7)

	r := gin.Default()

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	publisher := pubsub.NewPubSub(client)
	manager := hub.NewManager(client)

	roomRepo := repositories.NewRoomRepo(db)
	roomService := services.NewRoomService(roomRepo, manager)
	roomHandler := handlers.NewRoomHandler(roomService, l, cfg.Cloudinary)

	messageRepo := repositories.NewMessageRepo(db)
	messageService := services.NewMessageService(messageRepo, publisher, client)
	messageHandler := handlers.NewMessageHandler(messageService, l)

	wsService := services.NewWsService(publisher, manager, l)
	wsHandler := handlers.NewWSConn(wsService, l, messageService)

	joinRepo := repositories.NewJoinRoomRepository(db)
	joinService := services.NewJoinRoomService(joinRepo, publisher)
	sessionHandle := session.CreateStore(client)
	joinHandler := handlers.NewJoinRoomHandler(joinService, l, cfg.Cloudinary, sessionHandle)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, l, sessionHandle, cfg.Cloudinary)

	routes.SetupRoutes(r, roomHandler, joinHandler, userHandler, sessionHandle, wsHandler, messageHandler)

	return &App{
		DB:     db,
		Server: srv,
		Redis:  client,
		Log:    l,
	}

}
