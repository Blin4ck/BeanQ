package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SessionSecret string
	DBDSN         string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки env файла")
	}
	SessionSecret = os.Getenv("SESSION_SECRET")
	DBDSN = os.Getenv("DNS")
}
