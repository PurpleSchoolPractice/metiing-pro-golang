package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/app"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/auth"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/eventParticipant"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/server"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type AppComponents struct {
	Config *configs.Config
	Logger *logger.Logger
	App    *app.App
	Router *chi.Mux
	Server *server.Server
}

func setupApplication() *AppComponents {

	// Инициализация конфигурации и логгера
	cfg := configs.LoadConfig()
	log := logger.NewLogger(cfg)

	// Создание основных компонентов
	application := app.NewApp()
	database := db.NewDB(cfg)
	router := chi.NewRouter()
	srv := server.NewServer(log, application, router)

	// Инициализация слоя данных
	userRepo := user.NewUserRepository(database)
	secretRepo := secret.NewSecretRepository(database, log)

	// Создаем JWT сервис с настройками
	jwtService := jwt.NewJWT(cfg.Auth.Secret)
	// Устанавливаем время жизни токенов
	jwtService.AccessTokenTTL = time.Minute * 4
	jwtService.RefreshTokenTTL = time.Minute * 5

	// Инициализация сервисов
	authService := auth.NewAuthService(userRepo, secretRepo, jwtService)

	// Регистрация обработчиков
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      cfg,
		AuthService: authService,
	})
	user.NewUserHandler(router, user.UserHandlerDeps{
		Config:         cfg,
		UserRepository: userRepo,
		JWTService:     jwtService,
	})

	// Инициализация репозитория событий
	eventRepo := event.NewEventRepository(database)

	// Регистрация обработчиков событий
	event.NewEventHandler(router, event.EventHandlerDeps{
		EventRepository: eventRepo,
		JWTService:      jwtService,
	})

	// Инициализация репозитория участников событий
	eventParticipantRepo := eventParticipant.NewEventParticipantRepository(database)

	// Регистрация обработчиков участников событий
	eventParticipant.NewEventParticipantHandler(router, eventParticipant.EventParticipantDepsHandler{
		EventParticipantRepository: eventParticipantRepo,
		JWTService:                 jwtService,
	})

	return &AppComponents{
		Config: cfg,
		Logger: log,
		App:    application,
		Router: router,
		Server: srv,
	}
}

func main() {
	// Запускаем Cobra
	//Execute()
	components := setupApplication()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := components.Server.Start(ctx); err != nil {
		components.Logger.Info(err.Error())
	}
}
