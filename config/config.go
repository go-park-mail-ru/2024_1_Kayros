package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server   `yaml:"server"`
		Postgres `yaml:"postgresql"`
		Minio    `yaml:"minio"`
		Redis    `yaml:"redis"`
	}

	Server struct {
		WriteTimeout     time.Duration `env-required:"true" yaml:"write-timeout"    env:"SRV_WRITE_TM"`
		ReadTimeout      time.Duration `env-required:"true" yaml:"read-timeout"    env:"SRV_READ_TM"`
		IdleTimeout      time.Duration `env-required:"true" yaml:"idle-timeout"    env:"SRV_IDLE_TM"`
		ShutdownDuration time.Duration `env-required:"true" yaml:"shutdown-duration"    env:"SRV_SHUTDOWN_DUR"`
		Host             string        `env-required:"true" yaml:"host"    env:"SRV_HOST"`
		Port             int           `env-required:"true" yaml:"port"    env:"SRV_PORT"`
	}

	Postgres struct {
		Host     string `env-required:"true" yaml:"host"    env:"PG_HOST"`
		Port     int    `env-required:"true" yaml:"port" env:"PG_PORT"`
		User     string `env-required:"true" yaml:"user" env:"PG_USER"`
		Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		Database string `env-required:"true" yaml:"database" env:"PG_DB"`
	}

	Minio struct {
		PortServer    int    `env-required:"true" yaml:"port-server" env:"M_SERVER_PORT"`
		ConsoleServer int    `env-required:"true" yaml:"console-server" env:"M_CONSOLE_PORT"`
		SecretKey     string `env-required:"true" yaml:"secret-key" env:"M_SECRET_KEY"`
		AccessKey     string `env-required:"true" yaml:"access-key" env:"M_ACCESS_KEY"`
	}

	Redis struct {
		Port int `env-required:"true" yaml:"port" env:"R_PORT"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	// нужно будет переписать --> нет смысла подключать целый пакет для чтения конфигурации
	err := cleanenv.ReadConfig("./config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("Ошибка конфигурации приложения: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Чтение конфигурации выполнено успешно")
	return cfg, nil
}
