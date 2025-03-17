package main

import (
	"fmt"
	"os"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/spf13/cobra"
)

var (
	// Флаги командной строки
	cfgFile    string
	serverPort string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	logLevel   string

	// Конфигурация приложения
	appConfig *configs.Config
)

// rootCmd представляет собой базовую команду при вызове без подкоманд
var rootCmd = &cobra.Command{
	Use:   "meeting-pro",
	Short: "Meeting Pro - приложение для управления встречами",
	Long:  `Meeting Pro - приложение для управления встречами и расписаниями встреч.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Загружаем конфигурацию
		appConfig = configs.LoadConfig()

		// Применяем флаги командной строки, если они заданы
		if cmd.Flags().Changed("port") {
			appConfig.Server.Port = serverPort
		}
		if cmd.Flags().Changed("db-host") {
			appConfig.Database.Host = dbHost
		}
		if cmd.Flags().Changed("db-port") {
			appConfig.Database.Port = dbPort
		}
		if cmd.Flags().Changed("db-user") {
			appConfig.Database.Username = dbUser
		}
		if cmd.Flags().Changed("db-password") {
			appConfig.Database.Password = dbPassword
		}
		if cmd.Flags().Changed("db-name") {
			appConfig.Database.Database = dbName
		}
		if cmd.Flags().Changed("log-level") {
			appConfig.Logging.Level = logLevel
		}
	},
}

// Execute добавляет все дочерние команды к корневой команде и устанавливает флаги.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	// Флаги, общие для всех команд
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "путь к конфигурационному файлу")
	rootCmd.PersistentFlags().StringVar(&serverPort, "port", "", "порт сервера")
	rootCmd.PersistentFlags().StringVar(&dbHost, "db-host", "", "хост базы данных")
	rootCmd.PersistentFlags().StringVar(&dbPort, "db-port", "", "порт базы данных")
	rootCmd.PersistentFlags().StringVar(&dbUser, "db-user", "", "пользователь базы данных")
	rootCmd.PersistentFlags().StringVar(&dbPassword, "db-password", "", "пароль базы данных")
	rootCmd.PersistentFlags().StringVar(&dbName, "db-name", "", "имя базы данных")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "уровень логирования (debug, info, warn, error)")
}

// GetConfig возвращает конфигурацию приложения
func GetConfig() *configs.Config {
	return appConfig
}
