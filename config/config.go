package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type (
	Project struct {
		Server            `yaml:"server"`
		Postgres          `yaml:"postgresql"`
		Minio             `yaml:"minio"`
		Redis             `yaml:"redis"`
		RestGrpcServer    `yaml:"rest-grpc-server"`
		AuthGrpcServer    `yaml:"auth-grpc-server"`
		CommentGrpcServer `yaml:"comment-grpc-server"`
		UserGrpcServer    `yaml:"user-grpc-server"`
		OrderGrpcServer   `yaml:"order-grpc-server"`
		SessionGrpcServer `yaml:"session-grpc-server"`
		Payment           `yaml:"payment"`
	}

	RestGrpcServer struct {
		Host string `yaml:"host" env:"REST_GRPC_HOST"`
		Port uint16 `yaml:"port" env:"REST_GRPC_PORT"`
	}

	AuthGrpcServer struct {
		Host string `yaml:"host" env:"AUTH_GRPC_HOST"`
		Port uint16 `yaml:"port" env:"AUTH_GRPC_PORT"`
	}

	CommentGrpcServer struct {
		Host string `yaml:"host" env:"CMNT_GRPC_HOST"`
		Port uint16 `yaml:"port" env:"CMNT_GRPC_PORT"`
	}

	UserGrpcServer struct {
		Host string `yaml:"host" env:"USER_GRPC_HOST"`
		Port uint16 `yaml:"port" env:"USER_GRPC_PORT"`
	}

	OrderGrpcServer struct {
		Host string `yaml:"host" env:"ORDER_GRPC_HOST"`
		Port uint16 `yaml:"port" env:"ORDER_GRPC_PORT"`
	}

	SessionGrpcServer struct {
		Host string `yaml:"host" env:"SESSION_GRPC_HOST"`
		Port uint16 `yaml:"port" env:"SESSION_GRPC_PORT"`
	}

	Server struct {
		Host             string `yaml:"host"    env:"SRV_HOST"`
		Port             uint16 `yaml:"port"    env:"SRV_PORT"`
		WriteTimeout     uint16 `yaml:"write-timeout"    env:"SRV_WRITE_TM"`
		ReadTimeout      uint16 `yaml:"read-timeout"    env:"SRV_READ_TM"`
		IdleTimeout      uint16 `yaml:"idle-timeout"    env:"SRV_IDLE_TM"`
		ShutdownDuration uint16 `yaml:"shutdown-duration"    env:"SRV_SHUTDOWN_DUR"`
		CsrfSecretKey    string `yaml:"csrf-secret-key" env:"CSRF_SECRET_KEY"`
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
		Host                 string `yaml:"host"    env:"R_HOST"`
		Port                 uint16 `yaml:"port" env:"R_PORT"`
		DatabaseSession      int    `yaml:"database-session" env:"R_DB_SESSION"`
		DatabaseCsrf         int    `yaml:"database-csrf" env:"R_DB_CSRF"`
		User                 string `yaml:"user" env:"R_USER"`
		Password             string `yaml:"password" env:"R_PASSWORD"`
	}

	Payment struct {
		SecretKey string `yaml:"secret-key" env:"P_SECRET_KEY"`
		StoreId   string `yaml:"store-id" env:"P_STORE_ID"`
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
