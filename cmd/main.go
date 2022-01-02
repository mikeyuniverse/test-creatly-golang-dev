package main

import (
	"creatly-task/internal/config"
	"creatly-task/internal/handlers"
	"creatly-task/internal/mongodb"
	"creatly-task/internal/repo"
	"creatly-task/internal/server"
	"creatly-task/internal/services"
	jwtauth "creatly-task/pkg/auth/jwt"
	"creatly-task/pkg/hasher"
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

	repo := repo.New(db, config.Repo)
	if err != nil {
		log.Fatalf(" - - - - - - - REPOSITORY NOT INIT.\n%s", err)
	}

	tokener := jwtauth.New(config.JWT)

	services := services.New(repo, tokener, storage)

	hasher := hasher.New(config.Auth.Salt)
	handlers := handlers.New(services, config.Files.Limit, hasher, config.JWT.TokenHeaderName, config.Auth.HeaderUserId)

	server := server.New(config.Server, handlers)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
