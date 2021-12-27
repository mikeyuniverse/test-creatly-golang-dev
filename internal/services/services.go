package services

import (
	"creatly-task/internal/models"
	"creatly-task/internal/repo"
	"errors"
	"fmt"
)

// TODO Create interface for Tokener
type Tokener interface {
	GenerateToken(userId string) (string, error)
}

// TODO Create interface for Cloud Storage
type CloudStorage interface{}

type Services struct {
	db      *repo.Repo
	tokener Tokener
	cloud   CloudStorage
}

func New(repo *repo.Repo, tokener Tokener, cloud CloudStorage) *Services {
	return &Services{
		db:      repo,
		tokener: tokener,
		cloud:   cloud,
	}
}

func (s *Services) SignUp(user *models.UserSignUpInput) error {
	return s.db.Users.CreateUser(user)
}

// TODO Implementation SignIn
func (s *Services) SignIn(user *models.UserSignInInput) (string, error) {
	userFromDB, err := s.db.Users.GetUserByCreds(user.Email)
	if err != nil {
		return "", err
	}

	if userFromDB.Password != user.PasswordHash {
		fmt.Printf("want password - %s\naccepted - %s\n", userFromDB.Password, user.PasswordHash)
		return "", errors.New("wrong password")
	}

	token, err := s.tokener.GenerateToken(userFromDB.UserID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Services) Files() ([]models.FileOut, error) {
	files, err := s.db.Files.All()
	if err != nil {
		return []models.FileOut{}, err
	}
	return files, nil
}

// TODO Implementation UploadFile
func (s *Services) UploadFile(file *models.FileUploadInput) error {
	// Upload file to cloud server
	// err := s.db.Files.AddLogUploadedFile(file)  // Create log for file
	// if err != nil {
	// 	// TODO Что делать с ошибкой? Логгировать или возвращать пользователю?
	// 	return err
	// }
	return nil
}

func (s *Services) GetUserIdByToken(token string) (int64, error) {
	userId, err := s.db.Tokens.GetUserIDByToken(token)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
