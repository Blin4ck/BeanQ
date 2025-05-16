package main

import (
	"log"
	"os"
	"test/repository"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка загрузки env файла")
	}
}

func main() {
	var db *gorm.DB
	var err error
	DNS := os.Getenv("DNS")
	var counter int
	delay := 1 * time.Second

	for counter < 5 {
		db, err = repository.ConectDB(DNS)

		if err == nil {
			log.Println("Успешное подключение к базе данных!")
			break
		}

		log.Printf("Ошибка подключения: %v\nНовая попытка через %v (попытка %d/5)", err, delay, counter+1)
		time.Sleep(delay)
		delay *= 2
		counter++
	}
	repo := &repository.Repository{DB: db}
	if err := repo.Migrate(); err != nil {
		log.Printf("Ошибка миграции: %v", err)
	}

}
