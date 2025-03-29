package main

import (
	"context"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"os"
	"os/signal"
	"syscall"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/app"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/server"
)

func main() {
	// Запускаем Cobra для обработки командной строки
	Execute()
	// Получаем конфигурацию после обработки командной строки
	confg := configs.LoadConfig()
	logging := logger.NewLogger(confg)
	application := app.NewApp()
	database := db.NewDB(confg)

	//Repository
	user.NewUserRepository(database)
	secret.NewSecretRepository(database, logging)

	srv := server.NewServer(logging, application)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := srv.Start(ctx); err != nil {
		logging.Info(err.Error())
	}
}
