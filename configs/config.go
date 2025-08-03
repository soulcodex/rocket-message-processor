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
type RedisConfig struct {
	RedisURL string `env:"URL" envDefault:"redis://localhost:6379"`
}

type HTTPConfig struct {
	HTTPHost         string `env:"HOST" envDefault:"0.0.0.0"`
	HTTPPort         int    `env:"PORT" envDefault:"8080"`
	HTTPReadTimeout  int    `env:"READ_TIMEOUT" envDefault:"30"`
	HTTPWriteTimeout int    `env:"WRITE_TIMEOUT" envDefault:"30"`
}

type UncategorizedConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}
type Config struct {
	AppConfig           `envPrefix:"APP_"`
	HTTPConfig          `envPrefix:"HTTP_"`
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
