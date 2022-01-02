package config

import (
	"os"
	"reflect"
	"testing"
)

func setEnv(envMap map[string]string) error {
	for key, value := range envMap {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Test_newServer(t *testing.T) {
	testTable := []struct {
		name     string
		prefix   string
		envMap   map[string]string
		expect   *Server
		received *Server
	}{
		{
			name:   "OK: without env vars",
			prefix: "SERVER",
			envMap: map[string]string{
				"SERVER_HOST": "",
				"SERVER_PORT": "",
			},
			expect: &Server{
				Port: "",
				Host: "",
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := setEnv(test.envMap)
			if err != nil {
				t.Fail()
			}
			config, err := newServer(test.prefix)
			if err != nil {
				t.Fail()
			}

			if !reflect.DeepEqual(config, test.expect) {
				t.Fail()
			}
		})
	}
}
