package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type (
	Project struct {
		Server   `yaml:"server"`
		Postgres `yaml:"postgresql"`
		Minio    `yaml:"minio"`
		Redis    `yaml:"redis"`
	}

	Server struct {
		Host             string `yaml:"host"    env:"SRV_HOST"`
		Port             uint16 `yaml:"port"    env:"SRV_PORT"`
		WriteTimeout     uint16 `yaml:"write-timeout"    env:"SRV_WRITE_TM"`
		ReadTimeout      uint16 `yaml:"read-timeout"    env:"SRV_READ_TM"`
		IdleTimeout      uint16 `yaml:"idle-timeout"    env:"SRV_IDLE_TM"`
		ShutdownDuration uint16 `yaml:"shutdown-duration"    env:"SRV_SHUTDOWN_DUR"`
		CsrfSecretKey    string `yaml:"csrf_secret_key" env:"CSRF_SECRET KEY"`
	}

	Postgres struct {
		Host            string `yaml:"host"    env:"PG_HOST"`
		Port            uint16 `yaml:"port" env:"PG_PORT"`
		User            string `yaml:"user" env:"PG_USER"`
		Password        string `yaml:"password" env:"PG_PASSWORD"`
		Database        string `yaml:"database" env:"PG_DB"`
		SslMode         string `yaml:"sslmode" env:"PG_SSL_MODE"`
		MaxOpenConns    uint32 `yaml:"max-open-connections" env:"PG_MAX_OPEN_CONN"`
		ConnMaxLifetime uint16 `yaml:"conn-max-lifetime" env:"PG_CONN_MAX_LIFETIME"`
		MaxIdleConns    uint32 `yaml:"max-idle-conns" env:"PG_IDLE_CONNS"`
		ConnMaxIdleTime uint16 `yaml:"conn-max-idle-time" env:"PG_MAX_IDLE_TIME"`
	}

	Minio struct {
		PortServer    uint16 `yaml:"port-server" env:"M_SERVER_PORT"`
		ConsoleServer uint16 `yaml:"console-server" env:"M_CONSOLE_PORT"`
		SecretKey     string `yaml:"secret-key" env:"M_SECRET_KEY"`
		AccessKey     string `yaml:"access-key" env:"M_ACCESS_KEY"`
		SslMode       bool   `yaml:"sslmode" env:"M_SSL_MODE"`
		Endpoint      string `yaml:"endpoint" env:"M_ENDPOINT"`
	}

	Redis struct {
		Host            string `yaml:"host"    env:"R_HOST"`
		Port            uint16 `yaml:"port" env:"R_PORT"`
		DatabaseSession int    `yaml:"database-session" env:"R_DB"`
		DatabaseCsrf    int    `yaml:"database-csrf" env:"R_DB"`
		User            string `yaml:"user" env:"R_USER"`
		Password        string `yaml:"password" env:"R_PASSWORD"`
	}
)

func NewConfig(logger *zap.Logger) *Project {
	cfg := &Project{}

	err := cleanenv.ReadConfig("config/config.yaml", cfg)
	if err != nil {
		logger.Fatal("Error reading application configuration", zap.Error(err))
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		logger.Fatal("Error creating configuration object", zap.Error(err))
	}

	logger.Info("Reading configuration successful")
	return cfg
}
