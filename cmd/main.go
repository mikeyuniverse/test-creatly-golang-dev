package main

import (
	"creatly-task/internal/config"
	"creatly-task/internal/mongodb"
	"creatly-task/internal/repo"
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

	repo.New(db, config.Repo)
	if err != nil {
		log.Fatalf(" - - - - - - - REPOSITORY NOT INIT.\n%s", err)
	}

	// services, err := services.New(repo)
	// if err != nil {
	// 	log.Fatalf(" - - - - - - - SERVICES NOT INIT.\n%s", err)
	// }

	// server := server.New(config.Server, services)

	// if err := server.Start(); err != nil {
	// 	log.Fatal(err)
	// }
}
