package main

import (
	"creatly-task/internal/config"
	"creatly-task/internal/mongodb"
	"creatly-task/internal/repo"
	"creatly-task/internal/server"
	"creatly-task/internal/services"
	jwtauth "creatly-task/pkg/auth/jwt"
	"creatly-task/pkg/hasher"
	"creatly-task/pkg/storage"
	"log"
)

// TODO ADD THIS CONSTANTS IN CONFIGURATION
const SALT = "923undwpinpwq3bp" // Salt for password hashing
const JWT_SIGNING_KEY = "aisdbup872d3bib28d3"
const JWT_TOKEN_TTL = 3600 // Seconds
const JWT_TOKEN_HEADER_NAME = "Authorization"
const AUTH_HEADER_USER_ID = "userID"

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

	tokener := jwtauth.New(JWT_SIGNING_KEY, JWT_TOKEN_TTL)

	services := services.New(repo, tokener, storage)

	hasher := hasher.New(SALT)
	handlers := server.NewHandlers(services, config.Files.Limit, hasher, JWT_TOKEN_HEADER_NAME, AUTH_HEADER_USER_ID)

	server := server.NewServer(config.Server, *handlers)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
