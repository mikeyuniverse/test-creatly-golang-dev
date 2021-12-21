package repo

import (
	"creatly-task/internal/config"
	"creatly-task/internal/mongodb"
)

type Users interface{}

type Tokens interface{}

type Files interface{}

type Repo struct {
	Users  Users
	Tokens Tokens
	Files  Files
}

func New(db *mongodb.Mongo, config *config.Repo) *Repo {
	return &Repo{
		Users:  newUsersRepo(db, config.UsersCollection),
		Tokens: newTokensRepo(db, config.TokensCollection),
		Files:  newFilesRepo(db, config.FilesCollection),
	}
}
