package repo

import (
	"creatly-task/internal/config"
	"creatly-task/internal/models"
	"creatly-task/internal/mongodb"
)

type Users interface {
	CreateUser(*models.UserSignUpInput) error
}

type Tokens interface{}

type Files interface {
	All() ([]models.FileOut, error)
}

type Repo struct {
	Users  Users
	Tokens Tokens
	Files  Files
}

// TODO Implementation of interfaces
func New(db *mongodb.Mongo, config *config.Repo) *Repo {
	return &Repo{
		Users:  newUsersRepo(db, config.UsersCollection),
		Tokens: newTokensRepo(db, config.TokensCollection),
		Files:  newFilesRepo(db, config.FilesCollection),
	}
}
