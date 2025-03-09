package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/app"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/server"
)

func main() {
	logging := logger.NewLogger()
	application := app.NewApp()

	srv := server.NewServer(logging, application)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		log.Println(err.Error())
	}
}
