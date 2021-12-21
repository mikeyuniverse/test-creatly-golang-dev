package main

import (
	"creatly-task/internal/config"
	"creatly-task/internal/mongodb"
	"creatly-task/internal/repo"
	"creatly-task/internal/server"
	"creatly-task/internal/services"
	"log"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalf(" - - - - - - - CONFIG NOT INIT.\n%s", err)
	}

	db, err := mongodb.New(config.Repo)
	if err != nil {
		log.Fatalf(" - - - - - - - DATABASE NOT INIT.\n%s", err)
	}

	repo := repo.New(db, config.Repo)
	if err != nil {
		log.Fatalf(" - - - - - - - REPOSITORY NOT INIT.\n%s", err)
	}

	services := services.New(repo)

	server := server.New(config.Server, services)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
