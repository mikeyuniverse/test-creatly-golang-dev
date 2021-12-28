package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	SERVER_PREFIX     = "SERVER"
	REPOSITORY_PREFIX = "MONGO"
	FILE_PREFIX       = "FILE"
	STORAGE_PREFIX    = "STORAGE"
)

type Server struct {
	Port string
}

func newServer(prefix string) (*Server, error) {
	var s Server
	err := envconfig.Process(prefix, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

type Repo struct {
	Host             string
	Port             string
	DatabaseName     string
	UsersCollection  string
	FilesCollection  string
	TokensCollection string
}

func newRepo(prefix string) (*Repo, error) {
	var r Repo
	err := envconfig.Process(prefix, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

type File struct {
	Limit int
}

func newFileConfig(prefix string) (*File, error) {
	var f File
	err := envconfig.Process(prefix, &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

type Storage struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
	Timeout    time.Duration
}

func newStorageConfig(prefix string) (*Storage, error) {
	var s Storage
	err := envconfig.Process(prefix, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

type Config struct {
	Server  *Server
	Repo    *Repo
	Files   *File
	Storage *Storage
}

func New(filename string) (*Config, error) {
	err := godotenv.Load(filename)
	if err != nil {
		return nil, err
	}

	server, err := newServer(SERVER_PREFIX)
	if err != nil {
		return nil, err
	}

	repo, err := newRepo(REPOSITORY_PREFIX)
	if err != nil {
		return nil, err
	}

	file, err := newFileConfig(FILE_PREFIX)
	if err != nil {
		return nil, err
	}

	storage, err := newStorageConfig(STORAGE_PREFIX)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server:  server,
		Repo:    repo,
		Files:   file,
		Storage: storage,
	}, nil
}
