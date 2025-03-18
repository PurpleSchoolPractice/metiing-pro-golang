package main

import (
	"fmt"
	"os"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cfg - структура конфигурации приложения
type Cfg struct {
	Logger   string
	Database configs.DatabaseConfig
	Server   configs.ServerConfig
}

var appConfig Cfg

// rootCmd представляет собой базовую команду при вызове без подкоманд
var rootCmd = &cobra.Command{
	Use:   "meeting-pro",
	Short: "Meeting Pro - приложение для управления встречами",
	Long:  `Meeting Pro - приложение для управления встречами и расписаниями встреч.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Загружаем конфигурацию с Viper
		initConfig()
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
	cobra.OnInitialize(initConfig)

	// Указываем путь к конфигурационному файлу
	rootCmd.PersistentFlags().String("config", "", "путь к конфигурационному файлу")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().String("port", "8080", "порт сервера")
	rootCmd.PersistentFlags().String("db-host", "localhost", "хост базы данных")
	rootCmd.PersistentFlags().String("db-port", "5432", "порт базы данных")
	rootCmd.PersistentFlags().String("db-user", "user", "пользователь базы данных")
	rootCmd.PersistentFlags().String("db-password", "password", "пароль базы данных")
	rootCmd.PersistentFlags().String("db-name", "dbname", "имя базы данных")
	rootCmd.PersistentFlags().String("log-level", "info", "уровень логирования (debug, info, warn, error)")

	viper.AutomaticEnv()
}

func initConfig() {
	if cfgFile := viper.GetString("config"); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Используется конфигурационный файл:", viper.ConfigFileUsed())
		//TODO Испрвить логирование
	}

	appConfig = Cfg{
		Logger: viper.GetString("log-level"),
		Database: configs.DatabaseConfig{
			Host:     viper.GetString("db-host"),
			Port:     viper.GetString("db-port"),
			Username: viper.GetString("db-user"),
			Password: viper.GetString("db-password"),
			Database: viper.GetString("db-name"),
		},
		Server: configs.ServerConfig{
			Port: viper.GetString("port"),
		},
	}
}
