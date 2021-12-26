package server

import (
	"creatly-task/internal/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *gin.Engine
}

func NewServer(config *config.Server, handlers Handlers) *Server {
	server := gin.Default()

	auth := server.Group("/auth")
	{
		auth.POST("/sign-up", handlers.SignUp)
		auth.POST("/sign-in", handlers.SignUp)
	}

	files := server.Group("/files")
	{
		files.Use(handlers.AuthMiddleware)
		files.GET("", handlers.Files)
		files.POST("/upload", handlers.UploadFile)
	}

	return &Server{
		httpServer: server,
	}
}

func (s *Server) Start() error {
	err := s.httpServer.Run()
	if err != nil {
		return err
	}

	return nil
}
