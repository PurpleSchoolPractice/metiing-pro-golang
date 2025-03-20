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
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запуск сервера приложения",
	Long:  `Запуск HTTP-сервера приложения Meeting Pro.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetConfig()
		logging := logger.NewLogger(cfg)
		application := app.NewApp()

		srv := server.NewServer(logging, application)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := srv.Start(ctx); err != nil {
			log.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(ServeCmd)
}
