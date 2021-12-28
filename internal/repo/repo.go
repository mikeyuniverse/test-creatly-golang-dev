package repo

import (
	"creatly-task/internal/config"
	"creatly-task/internal/models"
	"creatly-task/internal/mongodb"
)

type Users interface {
	CreateUser(*models.UserSignUpInput) error
	GetUserByCreds(email string) (*models.UserSignInOutput, error)
}

type Tokens interface {
	GetUserIDByToken(token string) (int64, error) // Получение userID по токену
}

type Files interface {
	All() ([]models.FileOut, error)
	AddLog(log *models.FileUploadLogInput) error
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
