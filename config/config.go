package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

var Config ProjectConfiguration

type (
	ProjectConfiguration struct {
		Server   `yaml:"server"`
		Postgres `yaml:"postgres"`
		Minio    `yaml:"minio"`
		Redis    `yaml:"redis"`

		RestGrpcServer         `yaml:"rest-grpc-server"`
		RestGrpcServerExporter `yaml:"rest-grpc-server-exporter"`

		AuthGrpcServer         `yaml:"auth-grpc-server"`
		AuthGrpcServerExporter `yaml:"auth-grpc-server-exporter"`

		CommentGrpcServer         `yaml:"comment-grpc-server"`
		CommentGrpcServerExporter `yaml:"comment-grpc-server-exporter"`

		UserGrpcServer         `yaml:"user-grpc-server"`
		UserGrpcServerExporter `yaml:"user-grpc-server-exporter"`

		SessionGrpcServer         `yaml:"session-grpc-server"`
		SessionGrpcServerExporter `yaml:"session-grpc-server-exporter"`

		Payment `yaml:"payment"`
		Oauth   `yaml:"oauth"`
	}

	RestGrpcServer struct {
		Host string `yaml:"host" env:"REST_GRPC_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"REST_GRPC_PORT" env-default:"8001"`
	}
	RestGrpcServerExporter struct {
		Host string `yaml:"host" env:"REST_GRPC_EXPORTER_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"REST_GRPC_EXPORTER_PORT" env-default:"8011"`
	}

	AuthGrpcServer struct {
		Host string `yaml:"host" env:"AUTH_GRPC_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"AUTH_GRPC_PORT" env-default:"8002"`
	}
	AuthGrpcServerExporter struct {
		Host string `yaml:"host" env:"AUTH_GRPC_EXPORTER_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"AUTH_GRPC_EXPORTER_PORT" env-default:"8012"`
	}

	CommentGrpcServer struct {
		Host string `yaml:"host" env:"CMNT_GRPC_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"CMNT_GRPC_PORT" env-default:"8003"`
	}
	CommentGrpcServerExporter struct {
		Host string `yaml:"host" env:"COMMENT_GRPC_EXPORTER_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"COMMENT_GRPC_EXPORTER_PORT" env-default:"8013"`
	}

	UserGrpcServer struct {
		Host string `yaml:"host" env:"USER_GRPC_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"USER_GRPC_PORT" env-default:"8004"`
	}
	UserGrpcServerExporter struct {
		Host string `yaml:"host" env:"USER_GRPC_EXPORTER_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"USER_GRPC_EXPORTER_PORT" env-default:"8014"`
	}

	SessionGrpcServer struct {
		Host string `yaml:"host" env:"SESSION_GRPC_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"SESSION_GRPC_PORT" env-default:"8005"`
	}
	SessionGrpcServerExporter struct {
		Host string `yaml:"host" env:"SESSION_GRPC_EXPORTER_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"SESSION_GRPC_EXPORTER_PORT" env-default:"8015"`
	}

	Server struct {
		Host             string        `yaml:"host"    env:"SRV_HOST" env-default:"localhost"`
		Port             int           `yaml:"port"    env:"SRV_PORT" env-default:"8000"`
		WriteTimeout     time.Duration `yaml:"write-timeout"    env:"SRV_WRITE_TM" env-default:"5s"`
		ReadTimeout      time.Duration `yaml:"read-timeout"    env:"SRV_READ_TM" env-default:"5s"`
		IdleTimeout      time.Duration `yaml:"idle-timeout"    env:"SRV_IDLE_TM" env-default:"20s"`
		ShutdownDuration time.Duration `yaml:"shutdown-duration"    env:"SRV_SHUTDOWN_DUR" env-default:"5s"`
		CsrfSecretKey    string        `yaml:"csrf-secret-key" env:"CSRF_SECRET_KEY"`
	}

	Postgres struct {
		Host            string        `yaml:"host"    env:"PG_HOST" env-default:"localhost"`
		Port            int           `yaml:"port" env:"PG_PORT" env-default:"5432"`
		User            string        `yaml:"user" env:"PG_USER" env-default:"kayrosteam"`
		Password        string        `yaml:"password" env:"PG_PASSWORD" env-default:"resto123"`
		Database        string        `yaml:"database" env:"PG_DB" env-default:"main"`
		SslMode         string        `yaml:"sslmode" env:"PG_SSL_MODE" env-default:"disable"`
		MaxOpenConns    int           `yaml:"max-open-connections" env:"PG_MAX_OPEN_CONN" env-default:"10"`
		ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime" env:"PG_CONN_MAX_LIFETIME" env-default:"30s"`
		MaxIdleConns    int           `yaml:"max-idle-conns" env:"PG_IDLE_CONNS" env-default:"50"`
		ConnMaxIdleTime time.Duration `yaml:"conn-max-idle-time" env:"PG_MAX_IDLE_TIME" env-default:"15s"`
	}

	Minio struct {
		PortServer    int    `yaml:"port-server" env:"M_SERVER_PORT" env-default:"9001"`
		ConsoleServer int    `yaml:"port-console" env:"M_CONSOLE_PORT" env-default:"9000"`
		SecretKey     string `yaml:"secret-key" env:"M_SECRET_KEY" env-default:"resto123"`
		AccessKey     string `yaml:"access-key" env:"M_ACCESS_KEY" env-default:"kayrosteam"`
		SslMode       bool   `yaml:"sslmode" env:"M_SSL_MODE" env-default:"false"`
		Endpoint      string `yaml:"endpoint" env:"M_ENDPOINT" env-default:"minio:9000"`
	}

	Redis struct {
		Host            string `yaml:"host"    env:"R_HOST" env-default:"redis"`
		Port            int    `yaml:"port" env:"R_PORT" env-default:"6379"`
		DatabaseSession int    `yaml:"database-session" env:"R_DB_SESSION" env-default:"0"`
		DatabaseCsrf    int    `yaml:"database-csrf" env:"R_DB_CSRF" env-default:"1"`
		User            string `yaml:"user" env:"R_USER" env-default:"kayrosteam"`
		Password        string `yaml:"password" env:"R_PASSWORD" env-default:"resto123"`
	}

	Payment struct {
		SecretKey string `yaml:"secret-key" env:"P_SECRET_KEY"`
		StoreId   string `yaml:"store-id" env:"P_STORE_ID"`
	}

	Oauth struct {
		AccessToken string `yaml:"access-token" env:"OAUTH_ACCESS_TOKEN"`
	}
)

func Read(logger *zap.Logger) {
	if err := cleanenv.ReadConfig("config/config.yaml", &Config); err != nil {
		logger.Fatal(fmt.Sprintf("error while reading application configuration: %v", err))
	}

	if err := cleanenv.ReadEnv(&Config); err != nil {
		logger.Fatal(fmt.Sprintf("error creating configuration object: %v", err))
	}
	logger.Info("reading configuration is successful")
}
