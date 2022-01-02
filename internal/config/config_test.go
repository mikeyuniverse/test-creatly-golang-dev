package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
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

func unsetEnv(envMap map[string]string) error {
	for key := range envMap {
		err := os.Unsetenv(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func Test_newServer(t *testing.T) {
	testTable := []struct {
		name      string
		prefix    string
		envMap    map[string]string
		wantError bool
		expect    *Server
		received  *Server
	}{
		{
			name:      "OK: without env vars",
			prefix:    "SERVER",
			wantError: false,
			envMap: map[string]string{
				"SERVER_HOST": "",
				"SERVER_PORT": "",
			},
			expect: &Server{
				Port: "",
				Host: "",
			},
		},
		{
			name:      "Fail: port not int",
			prefix:    "SERVER",
			wantError: true,
			envMap: map[string]string{
				"SERVER_HOST": "",
				"SERVER_PORT": "port",
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
				t.FailNow()
			}

			config, err := newServer(test.prefix)
			if err != nil && !test.wantError {
				t.FailNow()
			}

			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				t.FailNow()
			}

			_, err = strconv.Atoi(test.expect.Port)
			if err != nil && config.Port != "" && !test.wantError {
				fmt.Println(err)
				t.FailNow()
			}

			err = unsetEnv(test.envMap)
			if err != nil {
				t.FailNow()
			}

		})
	}
}

func Test_newRepo(t *testing.T) {
	testTable := []struct {
		name      string
		prefix    string
		expect    *Repo
		envMap    map[string]string
		wantError bool
	}{
		{
			name:   "OK",
			prefix: "REPO",
			expect: &Repo{
				Host:             "localhost",
				Port:             "80908",
				DatabaseName:     "storage",
				UsersCollection:  "users",
				FilesCollection:  "files",
				TokensCollection: "tokens",
			},
			envMap: map[string]string{
				"REPO_HOST":             "localhost",
				"REPO_PORT":             "80908",
				"REPO_DATABASENAME":     "storage",
				"REPO_USERSCOLLECTION":  "users",
				"REPO_FILESCOLLECTION":  "files",
				"REPO_TOKENSCOLLECTION": "tokens",
			},
			wantError: false,
		},
		{
			name:   "FAIL: invalid values",
			prefix: "REPO",
			expect: &Repo{
				Host:             "localhost",
				Port:             "port", // Here error
				DatabaseName:     "storage",
				UsersCollection:  "users",
				FilesCollection:  "files",
				TokensCollection: "tokens",
			},
			envMap: map[string]string{
				"REPO_HOST":             "localhost",
				"REPO_PORT":             "80908",
				"REPO_DATABASENAME":     "storage",
				"REPO_USERSCOLLECTION":  "users",
				"REPO_FILESCOLLECTION":  "files",
				"REPO_TOKENSCOLLECTION": "tokens",
			},
			wantError: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := setEnv(test.envMap)
			if err != nil {
				t.Fatalf("setEnv error - %s\n", err.Error())
			}

			config, err := newRepo(test.prefix)
			if err != nil && !test.wantError {
				t.Fatalf("config init error - %s\n", err.Error())
			}

			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				t.Fatalf("configs not equals\nReceived - %+v\nWant - %+v\n", config, test.expect)
			}

			err = unsetEnv(test.envMap)
			if err != nil {
				t.Fatalf("unset env error - %s\n", err.Error())
			}
		})
	}
}

func Test_newFileConfig(t *testing.T) {
	testTable := []struct {
		name      string
		prefix    string
		expect    *File
		envMap    map[string]string
		wantError bool
	}{
		{
			name:   "OK",
			prefix: "FILE",
			expect: &File{
				Limit: 123352350,
			},
			envMap: map[string]string{
				"FILE_LIMIT": "123352350",
			},
			wantError: false,
		},
		{
			name:   "FAIL: wrong type of variable",
			prefix: "FILE",
			expect: &File{
				Limit: 0,
			},
			envMap: map[string]string{
				"FILE_LIMIT": "limit",
			},
			wantError: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := setEnv(test.envMap)
			if err != nil {
				t.FailNow()
			}

			config, err := newFileConfig(test.prefix)
			if err != nil && !test.wantError {
				t.FailNow()
			}

			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				t.FailNow()
			}

			err = unsetEnv(test.envMap)
			if err != nil {
				t.FailNow()
			}
		})
	}
}

func Test_newStorageConfig(t *testing.T) {
	testTable := []struct {
		name      string
		envMap    map[string]string
		wantError bool
		expect    *Storage
		prefix    string
	}{
		{
			name:   "OK",
			prefix: "STORAGE",
			envMap: map[string]string{
				"STORAGE_ACCESSKEY":  "179g381vdyo",
				"STORAGE_SECRETKEY":  "18e721gf2fg01g711378gfjksog",
				"STORAGE_REGION":     "eu-west",
				"STORAGE_BUCKETNAME": "my-bucket",
				"STORAGE_TIMEOUT":    "60s",
			},
			expect: &Storage{
				AccessKey:  "179g381vdyo",
				SecretKey:  "18e721gf2fg01g711378gfjksog",
				Region:     "eu-west",
				BucketName: "my-bucket",
				Timeout:    time.Second * 60,
			},
			wantError: false,
		},
		{
			name:   "FAIL: accessKey not initialize",
			prefix: "STORAGE",
			envMap: map[string]string{
				"STORAGE_ACCESSKEY":  "",
				"STORAGE_SECRETKEY":  "18e721gf2fg01g711378gfjksog",
				"STORAGE_REGION":     "eu-west",
				"STORAGE_BUCKETNAME": "my-bucket",
				"STORAGE_TIMEOUT":    "60s",
			},
			expect: &Storage{
				AccessKey:  "179g381vdyo",
				SecretKey:  "18e721gf2fg01g711378gfjksog",
				Region:     "eu-west",
				BucketName: "my-bucket",
				Timeout:    time.Second * 60,
			},
			wantError: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := setEnv(test.envMap)
			if err != nil {
				fmt.Println("set env error")
				t.FailNow()
			}

			config, err := newStorageConfig(test.prefix)
			if err != nil && !test.wantError {
				fmt.Println("error init config")
				t.FailNow()
			}

			// Compare received config and expect config
			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				fmt.Println("config not equals")
				fmt.Printf("%+v\n", config)
				t.FailNow()
			}

			// Check timeout
			if test.expect.Timeout == 0 && !test.wantError {
				fmt.Println("invalid timeout error")
				fmt.Printf("%+v\n", config)
				t.FailNow()
			}

			// Check required params
			if (config.AccessKey == "" ||
				config.SecretKey == "" ||
				config.Region == "") &&
				!test.wantError {
				fmt.Println("required params must be not empty")
				fmt.Printf("%+v\n", config)
				t.FailNow()
			}

			err = unsetEnv(test.envMap)
			if err != nil {
				t.FailNow()
			}
		})
	}
}

func Test_newJWTConfig(t *testing.T) {
	testTable := []struct {
		name      string
		envMap    map[string]string
		wantError bool
		expect    *JWT
		prefix    string
	}{
		{
			name:   "OK",
			prefix: "JWT",
			envMap: map[string]string{
				"JWT_SIGNINGKEY":      "190fh[9iqn",
				"JWT_TOKENTTL":        "6000",
				"JWT_TOKENHEADERNAME": "Authorization",
			},
			expect: &JWT{
				SigningKey:      "190fh[9iqn",
				TokenTTL:        6000,
				TokenHeaderName: "Authorization",
			},
			wantError: false,
		},
		{
			name:   "FAIL: empty header name",
			prefix: "JWT",
			envMap: map[string]string{
				"JWT_SIGNINGKEY":      "190fh[9iqn",
				"JWT_TOKENTTL":        "6000",
				"JWT_TOKENHEADERNAME": "",
			},
			expect: &JWT{
				SigningKey:      "190fh[9iqn",
				TokenTTL:        6000,
				TokenHeaderName: "Authorization",
			},
			wantError: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := setEnv(test.envMap)
			if err != nil {
				t.Fatalf("setEnv error - %s\n", err.Error())
			}

			config, err := newJWTConfig(test.prefix)
			if err != nil && !test.wantError {
				t.Fatalf("config init error - %s\n", err.Error())
			}

			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				t.Fatalf("configs not equals\nReceived - %+v\nWant - %+v\n", config, test.expect)
			}

			if (config.TokenHeaderName == "" || config.TokenTTL == 0) && !test.wantError {
				t.Fatalf("empty value in required params\nReceived - %+v\nWant - %+v\n", config, test.expect)
			}

			err = unsetEnv(test.envMap)
			if err != nil {
				t.Fatalf("unsetEnv error - %s\n", err.Error())
			}
		})
	}
}

func Test_newAuthConfig(t *testing.T) {
	testTable := []struct {
		name      string
		envMap    map[string]string
		wantError bool
		expect    *Auth
		prefix    string
	}{
		{
			name:   "OK",
			prefix: "AUTH",
			envMap: map[string]string{
				"AUTH_SALT":         "23ffiuvbsa",
				"AUTH_HEADERUSERID": "user",
			},
			expect: &Auth{
				Salt:         "23ffiuvbsa",
				HeaderUserId: "user",
			},
			wantError: false,
		},
		{
			name:   "FAIL: empty HeaderUserId",
			prefix: "AUTH",
			envMap: map[string]string{
				"AUTH_SALT":         "23ffiuvbsa",
				"AUTH_HEADERUSERID": "user",
			},
			expect: &Auth{
				Salt:         "23ffiuvbsa",
				HeaderUserId: "",
			},
			wantError: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := setEnv(test.envMap)
			if err != nil {
				t.Fatalf("setEnv error - %s\n", err.Error())
			}

			config, err := newAuthConfig(test.prefix)
			if err != nil && !test.wantError {
				t.Fatalf("error init config - %s\n", err.Error())
			}

			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				t.Fatalf("configs not equals\nReceived - %+v\nWant - %+v\n", config, test.expect)
			}

			if config.HeaderUserId == "" && !test.wantError {
				t.Fatalf("empty header user id\nConfig - %+v\n", config)
			}

			err = unsetEnv(test.envMap)
			if err != nil {
				t.Fatalf("error unset env - %s\n", err.Error())
			}
		})
	}
}

func Test_New(t *testing.T) {
	testTable := []struct {
		name      string
		filepath  string
		envMap    map[string]string
		expect    *Config
		wantError bool
	}{
		{
			name:     "OK",
			filepath: "test files/ok.env",
			expect: &Config{
				Server: &Server{
					Host: "localhost",
					Port: "8000",
				},
				Repo: &Repo{
					Host:             "mongodb",
					Port:             "27017",
					DatabaseName:     "database_name",
					UsersCollection:  "users",
					FilesCollection:  "files",
					TokensCollection: "tokens",
				},
				Files: &File{
					Limit: 60001,
				},
				Storage: &Storage{
					AccessKey:  "AIOYFOSUDIFBSIYF",
					SecretKey:  "UOYVivOUYVOYVVPIVouyvp878P7Cyouv",
					Region:     "eu-north-1",
					BucketName: "mybucket",
					Timeout:    time.Second * 60,
				},
				JWT: &JWT{
					SigningKey:      "aisdbup872d3bib28d3",
					TokenTTL:        3600,
					TokenHeaderName: "Authorization",
				},
				Auth: &Auth{
					Salt:         "923undwpinpwq3bp",
					HeaderUserId: "userID",
				},
			},
			wantError: false,
		},
		{
			name:     "FAIL: ",
			filepath: "test files/fail.env",
			expect: &Config{
				Server: &Server{
					Host: "localhost",
					Port: "8000",
				},
				Repo: &Repo{
					Host:             "mongodb",
					Port:             "27017",
					DatabaseName:     "database_name",
					UsersCollection:  "users",
					FilesCollection:  "files",
					TokensCollection: "tokens",
				},
				Files: &File{
					Limit: 60001,
				},
				Storage: &Storage{
					AccessKey:  "AIOYFOSUDIFBSIYF",
					SecretKey:  "UOYVivOUYVOYVVPIVouyvp878P7Cyouv",
					Region:     "eu-north-1",
					BucketName: "mybucket",
					Timeout:    time.Second * 60,
				},
				JWT: &JWT{
					SigningKey:      "aisdbup872d3bib28d3",
					TokenTTL:        3600,
					TokenHeaderName: "Authorization",
				},
				Auth: &Auth{
					Salt:         "923undwpinpwq3bp",
					HeaderUserId: "userID",
				},
			},
			wantError: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			config, err := New(test.filepath)
			if err != nil {
				t.Fatalf("init config error - %s\n", err.Error())
			}

			// Compare received and expected configs
			if !reflect.DeepEqual(config, test.expect) && !test.wantError {
				fmt.Println("configs not equals")
				fmt.Printf("Server - %+v\n", config.Server)
				fmt.Printf("Repo - %+v\n", config.Repo)
				fmt.Printf("Auth - %+v\n", config.Auth)
				fmt.Printf("Files - %+v\n", config.Files)
				fmt.Printf("Storage - %+v\n", config.Storage)
				fmt.Printf("JWT - %+v\n", config.JWT)
				t.FailNow()
			}

			// Check required fields

			err = unsetEnv(test.envMap)
			if err != nil {
				t.Fatalf("unsetEnv error - %s\n", err.Error())
			}
		})
	}
}
