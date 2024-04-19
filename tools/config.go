package tools

import (
	"github.com/spf13/viper"
	"log"
)

// Config Какие переменные среды вытаскиваем из .env
type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBConnect     string `mapstructure:"DB_CONNECT"`
}

// LoadConfig Загружаем переменные среды (надо подавать путь относительно текущей папки)
func LoadConfig(path string) (config Config) {
	// Указываем файл
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Загружаем переменные
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Не получилось загрузить конфигурации:", err)
	}

	return
}
