package tools

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

// Config Какие переменные среды вытаскиваем из .env
type Config struct {
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	DBUserName           string        `mapstructure:"POSTGRES_USER"`
	DBPassword           string        `mapstructure:"POSTGRES_PASSWORD"`
	DBHost               string        `mapstructure:"DB_HOST"`
	DBSSLMode            string        `mapstructure:"DB_SSL_MODE"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig Загружаем переменные среды (надо подавать путь к .env относительно текущей папки)
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
