package main

import (
	"log"
	"test/config"
	"test/repository"
	"test/router"
)

func main() {
	config.LoadEnv()

	db, err := repository.RetryConectDB(config.DBDSN, 5)

	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	repo := &repository.Repository{DB: db}
	if err := repo.Migrate(); err != nil {
		log.Printf("Ошибка миграции: %v", err)
	}

	r := router.SetupRouter(db)
	r.Run(":8080")
}
