package main

import (
	"context"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/auth"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/go-chi/chi/v5"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"os"
	"os/signal"
	"syscall"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/app"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/server"
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
	secret.NewSecretRepository(database, log)

	// Инициализация сервисов
	authService := auth.NewAuthService(userRepo)

	// Регистрация обработчиков
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      cfg,
		AuthService: authService,
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
	Execute()
	components := setupApplication()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := components.Server.Start(ctx); err != nil {
		components.Logger.Info(err.Error())
	}
}
