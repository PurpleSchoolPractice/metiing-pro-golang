package main

import (
	"context"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/app"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/server"
)

func main() {
	cfg := configs.LoadConfig()
	logging := logger.NewLogger(cfg)
	application := app.NewApp()
	database := db.NewDB(cfg)

	srv := server.NewServer(logging, application)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		log.Println(err.Error())
	}
}
