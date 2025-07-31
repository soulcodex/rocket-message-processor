package configs

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppServiceName string `env:"SERVICE_NAME"`
	AppEnv         string `env:"ENV"`
	AppVersion     string `env:"VERSION"`
}

type PostgresConfig struct {
	PostgresHost   string `env:"HOST" envDefault:"localhost"`
	PostgresPort   int    `env:"PORT" envDefault:"5432"`
	PostgresUser   string `env:"USER" envDefault:"rockets_message_processor_role"`
	PostgresPass   string `env:"PASSWORD" envDefault:"rockets_message_processor"`
	PostgresSchema string `env:"SCHEMA" envDefault:"rockets_message_processor"`
	PostgresDB     string `env:"DATABASE" envDefault:"rockets_message_processor"`
	PostgresSSL    string `env:"SSL_MODE" envDefault:"disable"`
}

type RedisConfig struct {
	RedisHosts    string `env:"HOSTS" envDefault:"localhost:6379"`
	RedisPassword string `env:"PASSWORD" envDefault:""`
}

type HTTPConfig struct {
	HTTPHost         string `env:"HOST" envDefault:"0.0.0.0"`
	HTTPPort         int    `env:"PORT" envDefault:"8080"`
	HTTPReadTimeout  int    `env:"READ_TIMEOUT" envDefault:"30"`
	HTTPWriteTimeout int    `env:"WRITE_TIMEOUT" envDefault:"30"`
}

type UncategorizedConfig struct {
	JSONSchemaBasePath string `env:"JSON_SCHEMA_BASE_PATH" envDefault:"./schemas"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"debug"`
}
type Config struct {
	AppConfig           `envPrefix:"APP_"`
	HTTPConfig          `envPrefix:"HTTP_"`
	PostgresConfig      `envPrefix:"POSTGRES_"`
	RedisConfig         `envPrefix:"REDIS_"`
	UncategorizedConfig `envPrefix:""`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	config, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
