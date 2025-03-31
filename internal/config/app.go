package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var AppConfig *appConfig
var onceApp sync.Once

type appConfig struct {
	Port string
}

func InitAppConfig() {
	onceApp.Do(func () {
		err := godotenv.Load();

		if err != nil {
			log.Println("No .env file found, using system environment variables")
		}

		AppConfig = &appConfig{
			Port: getEnv("PORT", ":8080"),
		}
	})
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}