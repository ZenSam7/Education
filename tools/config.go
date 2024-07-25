package tools

import (
	"github.com/spf13/viper"
	"time"
)

// Config Какие переменные среды вытаскиваем из .env
type Config struct {
	HttpServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	DBUserName           string        `mapstructure:"POSTGRES_USER"`
	DBPassword           string        `mapstructure:"POSTGRES_PASSWORD"`
	DBHost               string        `mapstructure:"DB_HOST"`
	DBSSLMode            string        `mapstructure:"DB_SSL_MODE"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig Загружаем переменные среды (надо подавать путь к .env относительно
// файла из текущей папки (откуда эта функция вызывается))
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
		Log.Fatal().Err(err).Msg("Не получилось загрузить конфигурации")
	}

	return
}
