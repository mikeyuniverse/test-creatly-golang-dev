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
	JWT_PREFIX        = "JWT"
	AUTH_PREFIX       = "AUTH"
)

type Server struct {
	Port string
	Host string
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

type JWT struct {
	SigningKey      string
	TokenTTL        int64
	TokenHeaderName string
}

func newJWTConfig(prefix string) (*JWT, error) {
	var j JWT
	err := envconfig.Process(prefix, &j)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

type Auth struct {
	Salt         string
	HeaderUserId string
}

func newAuthConfig(prefix string) (*Auth, error) {
	var a Auth
	err := envconfig.Process(prefix, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type Config struct {
	Server  *Server
	Repo    *Repo
	Files   *File
	Storage *Storage
	JWT     *JWT
	Auth    *Auth
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

	jwtConfig, err := newJWTConfig(JWT_PREFIX)
	if err != nil {
		return nil, err
	}

	authConfig, err := newAuthConfig(AUTH_PREFIX)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server:  server,
		Repo:    repo,
		Files:   file,
		Storage: storage,
		JWT:     jwtConfig,
		Auth:    authConfig,
	}, nil
}
