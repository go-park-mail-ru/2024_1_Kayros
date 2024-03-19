package app

import (
	"log"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %s", err)
	}

	app.Run(cfg)
	log.Printf("Сервер завершил работу")
}
