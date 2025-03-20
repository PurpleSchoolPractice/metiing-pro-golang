package main

import (
	"context"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"

	"log"
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
	cfg := GetConfig()

	logging := logger.NewLogger(cfg)

	application := app.NewApp()

	srv := server.NewServer(logging, application)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		log.Println(err.Error())
	}
}
