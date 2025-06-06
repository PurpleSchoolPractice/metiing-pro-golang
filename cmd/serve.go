package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/go-chi/chi/v5"

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
		cfg := configs.LoadConfig()
		logging := logger.NewLogger(cfg)
		router := chi.NewRouter()
		application := app.NewApp()

		srv := server.NewServer(logging, application, router)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := srv.Start(ctx); err != nil {
			logging.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(ServeCmd)
}
