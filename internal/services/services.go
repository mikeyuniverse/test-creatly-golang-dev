package services

import (
	"creatly-task/internal/models"
	"creatly-task/internal/repo"
	"errors"
	"fmt"
	"time"
)

type Tokener interface {
	GenerateToken(userId string) (string, error)
	ParseToken(token string) (string, error)
}

type CloudStorage interface {
	UploadFile(file []byte, filesize int64, filename string) (string, error)
}

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

func (s *Services) UploadFile(file *models.FileUploadInput) error {
	url, err := s.cloud.UploadFile(file.FileData, file.Size, file.Filename)
	if err != nil {
		return err
	}

	err = s.db.Files.AddLog(&models.FileUploadLogInput{
		Size:       file.Size,
		UploadDate: time.Now().Unix(),
		Filename:   file.Filename,
		UserId:     file.UserId,
		Url:        url,
	})
	if err != nil {
		return fmt.Errorf("error with log uploaded file - %s", err.Error())
	}

	return nil
}

func (s *Services) GetUserIdByToken(token string) (int64, error) {
	userId, err := s.db.Tokens.GetUserIDByToken(token)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *Services) ParseToken(token string) (string, error) {
	userID, err := s.tokener.ParseToken(token)
	if err != nil {
		return "", err
	}
	return userID, nil
}
