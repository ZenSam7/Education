package tools

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
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
	DBReplicaHost        string        `mapstructure:"DB_REPLICA"`
	DBSSLMode            string        `mapstructure:"DB_SSL_MODE"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	CacheDuration        time.Duration `mapstructure:"CHACHE_DURATUON"`
	EmailHost            string        `mapstructure:"EMAIL_HOST"`
	EmailPort            int           `mapstructure:"EMAIL_PORT"`
	EmailPassword        string        `mapstructure:"EMAIL_PASSWORD"`
	EmailSender          string        `mapstructure:"EMAIL_SENDER"`
}

// absolutePath Абсолютный путь к корневому каталогу проекта
var absolutePath string

// PathSeparator "/" или "\"
var PathSeparator string

// findAbsolutePath Вне зависимости откуда вызываем функцию (хоть из теста,
// хоть из пакета), находим путь к корневому проекту
func findAbsolutePath() {
	// Если у нас пока нету полного пути к корню, то значит, мы в первый раз вызвали эту функцию
	if len(absolutePath) != 0 {
		return
	}

	currentPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось узнать абсолютный путь")
	}

	// В зависимости от ОС разные разделители
	if strings.Contains(currentPath, "/") {
		PathSeparator = "/"
	} else {
		PathSeparator = "\\"
	}

	// Постепенно движимся от текущего пакета всё ближе к корневому каталогу, где находится файл main.go
	// (я специально сделал посложнее, чтобы можно было проект переименовать проект как угодно и этот алгоритм
	// не зависел от названия проекта (названия папки корневой директории))
	tempPath := strings.Split(currentPath, PathSeparator)
	for i := len(tempPath); i > 0; i-- {
		checkDir := strings.Join(tempPath[:i], PathSeparator)
		filesInCheckDir, err := os.ReadDir(checkDir)
		if err != nil {
			log.Fatal().Err(err).Msg("не получилось узнать нахождение корневого файла")
		}

		// Содержит ли checkDir main.go
		for _, file := range filesInCheckDir {
			// "mail.go" для запуска с компа, "main" для докер контейнера
			if file.Name() == "main.go" || file.Name() == "main" {
				absolutePath = checkDir
				return
			}
		}
	}

	log.Fatal().Msgf("не получилось узнать корневую директорию каталога, проблема в алгоритме")
}

// LoadConfig Загружаем переменные среды
func LoadConfig() (config Config) {
	findAbsolutePath()

	// Указываем файл
	viper.AddConfigPath(absolutePath)
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
		log.Fatal().Err(err).Msg("Не получилось загрузить конфигурации")
	}

	return
}
