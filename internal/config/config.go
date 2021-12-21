package config

import "github.com/kelseyhightower/envconfig"

const (
	SERVER_PREFIX     = "SERVER"
	REPOSITORY_PREFIX = "POSTGRES"
)

type Server struct {
	Port string
}

func newServer(prefix string) (*Server, error) {
	var s *Server
	err := envconfig.Process(prefix, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type Repo struct{}

func newRepo(prefix string) (*Repo, error) {
	var r *Repo
	err := envconfig.Process(prefix, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type Config struct {
	Server *Server
	Repo   *Repo
}

func New() (*Config, error) {
	server, err := newServer(SERVER_PREFIX)
	if err != nil {
		return nil, err
	}

	repo, err := newRepo(REPOSITORY_PREFIX)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: server,
		Repo:   repo,
	}, nil
}
