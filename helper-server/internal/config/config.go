package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	AppVersion  string `env:"APP_VERSION" env-default:"0.0.0"`
	AppInstance string `env:"APP_INSTANCE" env-default:""`

	AnalyticsSecret string `env:"ANALYTICS_SECRET" env-default:""`
	FeedbackSecret  string `env:"FEEDBACK_SECRET" env-default:""`

	DatabaseDSN string `env:"DB_DSN"`
}

func LoadLaunchConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	if v, err := os.ReadFile("version"); err == nil {
		cfg.AppVersion = string(v)
	} else {
		log.Println("File 'version' not found, using version from env")
	}

	return cfg, nil
}
