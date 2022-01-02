package server

import (
	"creatly-task/internal/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *gin.Engine
	port       string
	host       string
}

type Handlers interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	AuthMiddleware(c *gin.Context)
	Files(c *gin.Context)
	UploadFile(c *gin.Context)
}

func New(config *config.Server, handlers Handlers) *Server {
	server := gin.Default()
	server.MaxMultipartMemory = 8 << 20 // 8 MiB

	auth := server.Group("/")
	{
		auth.POST("/sign-up", handlers.SignUp)
		auth.POST("/sign-in", handlers.SignIn)
	}

	files := server.Group("/")
	{
		files.Use(handlers.AuthMiddleware)
		files.GET("/files", handlers.Files)
		files.POST("/upload", handlers.UploadFile)
	}

	return &Server{
		httpServer: server,
		port:       config.Port,
		host:       config.Host,
	}
}

func (s *Server) Start() error {
	err := s.httpServer.Run(fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}

	return nil
}
