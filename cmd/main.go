package main

import (
	"creatly-task/internal/config"
	"creatly-task/internal/mongodb"
	"creatly-task/internal/repo"
	"creatly-task/internal/server"
	"creatly-task/internal/services"
	"creatly-task/pkg/storage"
	"log"
)

func main() {
	config, err := config.New(".env")
	if err != nil {
		log.Fatalf(" - - - - - - - CONFIG NOT INIT.\n%s", err)
	}

	db, err := mongodb.New(config.Repo)
	if err != nil {
		log.Fatalf(" - - - - - - - DATABASE NOT INIT.\n%s", err)
	}

	storage, err := storage.New(config.Storage)
	if err != nil {
		log.Fatalf(" - - - - - - - STORAGE NOT INIT.\n%s", err)
	}

	repo := repo.New(db, storage, config.Repo)
	if err != nil {
		log.Fatalf(" - - - - - - - REPOSITORY NOT INIT.\n%s", err)
	}

	services := services.New(repo)

	handlers := server.NewHandlers(services, config.Files.Limit)

	server := server.NewServer(config.Server, *handlers)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
