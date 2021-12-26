package services

import (
	"creatly-task/internal/repo"

	"github.com/gin-gonic/gin"
)

type Services struct {
	db *repo.Repo
}

func New(repo *repo.Repo) *Services {
	return &Services{db: repo}
}

func (s *Services) AuthMiddleware(c *gin.Context) {}

func (s *Services) SignUp(c *gin.Context) {}

func (s *Services) SignIn(c *gin.Context) {}

func (s *Services) Files(c *gin.Context) {}

func (s *Services) UploadFile(c *gin.Context) {}
