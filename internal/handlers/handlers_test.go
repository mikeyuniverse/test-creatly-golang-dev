package handlers

import (
	"bytes"
	mock_handlers "creatly-task/internal/handlers/mocks"
	"creatly-task/internal/models"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SignUp(t *testing.T) {
	testTable := []struct {
		name           string
		inputPassword  string
		signUpInput    models.UserSignUpInput
		bodyInput      string
		behavior       func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices, passwordInput, passwordOutput string, err error, signUp *models.UserSignUpInput, signUpError error)
		outStatusCode  int
		outMessage     string
		passwordOutput string
		errorOutput    error
		signUpError    error
	}{
		{
			name:          "OK",
			inputPassword: "password",
			bodyInput:     `{"email": "some@mail.com", "password": "password"}`,
			behavior: func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices, passwordInput, passwordOutput string, hashErr error, signUp *models.UserSignUpInput, signUpError error) {
				mh.EXPECT().Hash(passwordInput).Return(passwordOutput, hashErr)
				s.EXPECT().SignUp(signUp).Return(signUpError)
			},
			outStatusCode:  200,
			outMessage:     `{"message":"success"}`,
			passwordOutput: "drowssap",
			errorOutput:    nil,
			signUpInput: models.UserSignUpInput{
				Email:    "some@mail.com",
				Password: "drowssap",
			},
			signUpError: nil,
		},
		{
			name:          "ERROR: wrong input",
			inputPassword: "password",
			bodyInput:     ``, // Error is here
			behavior: func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices, passwordInput, passwordOutput string, hashErr error, signUp *models.UserSignUpInput, signUpError error) {
			},
			outStatusCode: 400,
			outMessage:    `{"message":"invalid input"}`,
		},
		{
			name:          "ERROR: hasher error",
			inputPassword: "password",
			bodyInput:     `{"email": "some@mail.com", "password": "password"}`,
			behavior: func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices, passwordInput, passwordOutput string, hashErr error, signUp *models.UserSignUpInput, signUpError error) {
				mh.EXPECT().Hash(passwordInput).Return(passwordOutput, hashErr)
				s.EXPECT().SignUp(signUp).Return(signUpError)
			},
			outStatusCode:  500,
			outMessage:     `{"message":"error while creating an account"}`,
			passwordOutput: "drowssap",
			errorOutput:    nil,
			signUpInput: models.UserSignUpInput{
				Email:    "some@mail.com",
				Password: "drowssap",
			},
			signUpError: errors.New("some error"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hasher := mock_handlers.NewMockHasher(ctrl)
			services := mock_handlers.NewMockServices(ctrl)

			test.behavior(hasher, services, test.inputPassword, test.passwordOutput, test.errorOutput, &test.signUpInput, test.signUpError)

			handlers := New(services, 100000, hasher, "Authorization", "userId")

			r := gin.New()
			r.POST("/sign-up", handlers.SignUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(test.bodyInput))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.outStatusCode, w.Code)
			assert.Equal(t, test.outMessage, w.Body.String())
		})
	}
}

func Test_SignIn(t *testing.T) {
	testTable := []struct {
		name           string
		behavior       func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices)
		bodyInput      string
		outStatusCode  int
		outMessage     string
		outHeaderValue string
	}{
		{
			name:          "OK",
			bodyInput:     `{"email": "some@mail.com", "password": "qwerty"}`,
			outStatusCode: 200,
			outMessage:    `{"token":"token"}`,
			behavior: func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices) {
				mh.EXPECT().Hash("qwerty").Return("ytrewq", nil)
				s.EXPECT().SignIn(&models.UserSignInInput{
					Email:        "some@mail.com",
					PasswordHash: "ytrewq",
				}).Return("token", nil)
			},
			outHeaderValue: "Bearer token",
		},
		{
			name:           "ERROR: invalid creds",
			bodyInput:      `{"email": "", "password": "qwerty"}`,
			outStatusCode:  400,
			outMessage:     `{"message":"invalid credentials"}`,
			behavior:       func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices) {},
			outHeaderValue: "",
		},
		{
			name:          "ERROR: hash password error",
			bodyInput:     `{"email": "some@mail.com", "password": "qwerty"}`,
			outStatusCode: 500,
			outMessage:    `{"message":"error while encrypting password"}`,
			behavior: func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices) {
				mh.EXPECT().Hash("qwerty").Return("", errors.New("internal error"))
			},
			outHeaderValue: "",
		},
		{
			name:          "ERROR: sign-in service error",
			bodyInput:     `{"email": "some@mail.com", "password": "qwerty"}`,
			outStatusCode: 400,
			outMessage:    `{"message":"invalid creds"}`,
			behavior: func(mh *mock_handlers.MockHasher, s *mock_handlers.MockServices) {
				mh.EXPECT().Hash("qwerty").Return("ytrewq", nil)
				s.EXPECT().SignIn(&models.UserSignInInput{
					Email:        "some@mail.com",
					PasswordHash: "ytrewq",
				}).Return("", errors.New("internal error"))
			},
			outHeaderValue: "",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hasher := mock_handlers.NewMockHasher(ctrl)
			services := mock_handlers.NewMockServices(ctrl)

			test.behavior(hasher, services)

			handlers := New(services, 100000, hasher, "Authorization", "userId")

			r := gin.New()
			r.POST("/sign-in", handlers.SignIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(test.bodyInput))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.outStatusCode, w.Code)
			assert.Equal(t, test.outMessage, w.Body.String())
			assert.Equal(t, w.Header().Get("Authorization"), test.outHeaderValue)
		})
	}
}

func Test_AuthMiddleware(t *testing.T) {
	testTable := []struct {
		name              string
		AuthHeaderName    string
		AuthHeaderValue   string
		wantError         bool
		userIdHeaderName  string
		userIdHeaderValue string
		behavior          func(s *mock_handlers.MockServices)
		statusCode        int
	}{
		{
			name:              "OK",
			AuthHeaderName:    "Authorization",
			AuthHeaderValue:   "Bearer token",
			wantError:         true,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "1",
			behavior: func(s *mock_handlers.MockServices) {
				s.EXPECT().ParseToken("token").Return("1", nil)
			},
			statusCode: 200,
		},
		{
			name:              "ERROR: invalid token",
			AuthHeaderName:    "Authorization",
			AuthHeaderValue:   "token", // Error here - "Bearer" not include
			wantError:         false,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "1",
			behavior:          func(s *mock_handlers.MockServices) {},
			statusCode:        401,
		},
		{
			name:              "ERROR: parse token error",
			AuthHeaderName:    "Authorization",
			AuthHeaderValue:   "Bearer token",
			wantError:         true,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "1",
			behavior: func(s *mock_handlers.MockServices) {
				s.EXPECT().ParseToken("token").Return("", errors.New("parse error"))
			},
			statusCode: 401,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hasher := mock_handlers.NewMockHasher(ctrl)
			services := mock_handlers.NewMockServices(ctrl)

			test.behavior(services)

			handlers := New(services, 100000, hasher, test.AuthHeaderName, test.userIdHeaderName)

			r := gin.Default()
			r.GET("/test", handlers.AuthMiddleware)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Add(test.AuthHeaderName, test.AuthHeaderValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Result().StatusCode, test.statusCode)
		})
	}
}

func Test_Files(t *testing.T) {
	testTable := []struct {
		name          string
		behavior      func(s *mock_handlers.MockServices)
		outBody       string
		outStatusCode int
	}{
		{
			name: "OK",
			behavior: func(s *mock_handlers.MockServices) {
				s.EXPECT().Files().Return([]models.FileOut{
					{
						Filename: "file_1.png",
						Size:     2000,
						Date:     19674823,
						UserId:   "1",
						Url:      "https://s3.storage.com/file_1.png",
					},
				}, nil)
			},
			outBody:       `[{"filename":"file_1.png","size":2000,"uploadDate":19674823,"userId":"1","url":"https://s3.storage.com/file_1.png"}]`,
			outStatusCode: 200,
		},
		{
			name: "ERROR: service files return error",
			behavior: func(s *mock_handlers.MockServices) {
				s.EXPECT().Files().Return([]models.FileOut{}, errors.New("error"))
			},
			outBody:       `{"message":"error getting file data"}`,
			outStatusCode: 500,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hasher := mock_handlers.NewMockHasher(ctrl)
			services := mock_handlers.NewMockServices(ctrl)

			test.behavior(services)

			handlers := New(services, 100000, hasher, "Authorization", "userId")

			r := gin.New()
			r.GET("/files", handlers.Files)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/files", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.outStatusCode, w.Code)    // Check status code
			assert.Equal(t, test.outBody, w.Body.String()) // Check body
		})
	}
}

func Test_UploadFile(t *testing.T) {
	testTable := []struct {
		name              string
		behavior          func(s *mock_handlers.MockServices)
		outStatusCode     int
		outBody           string
		wantError         bool
		userIdHeaderName  string
		userIdHeaderValue string
		contentType       string
	}{
		{
			name: "OK",
			behavior: func(s *mock_handlers.MockServices) {
				s.EXPECT().UploadFile(&models.FileUploadInput{
					Filename: fmt.Sprintf("%s-%d.png", "1", time.Now().Unix()),
					Size:     7,
					UserId:   "1",
					FileData: []byte{49, 50, 51, 52, 53, 54, 55},
				}).Return(nil)
			},
			outStatusCode:     200,
			outBody:           `{"message":"upload success"}`,
			wantError:         false,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "1",
			contentType:       "image/png",
		},
		{
			name:              "ERROR: wrong content-type",
			behavior:          func(s *mock_handlers.MockServices) {},
			outStatusCode:     400,
			outBody:           `{"message":"invalid content-type"}`,
			wantError:         true,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "1",
			contentType:       "csv/text",
		},
		{
			name:              "ERROR: invalid userId",
			behavior:          func(s *mock_handlers.MockServices) {},
			outStatusCode:     401,
			outBody:           `{"message":"invalid userID"}`,
			wantError:         true,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "",
			contentType:       "image/png",
		},
		{
			name: "ERROR: file uploading error",
			behavior: func(s *mock_handlers.MockServices) {
				s.EXPECT().UploadFile(&models.FileUploadInput{
					Filename: fmt.Sprintf("%s-%d.png", "1", time.Now().Unix()),
					Size:     7,
					UserId:   "1",
					FileData: []byte{49, 50, 51, 52, 53, 54, 55},
				}).Return(errors.New("upload err"))
			},
			outStatusCode:     500,
			outBody:           `{"message":"error with upload file"}`,
			wantError:         true,
			userIdHeaderName:  "userId",
			userIdHeaderValue: "1",
			contentType:       "image/png",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hasher := mock_handlers.NewMockHasher(ctrl)
			services := mock_handlers.NewMockServices(ctrl)

			test.behavior(services)

			handlers := New(services, 100000, hasher, "Authorization", test.userIdHeaderName)

			r := gin.New()
			r.POST("/upload", handlers.UploadFile)

			// Create Request
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{49, 50, 51, 52, 53, 54, 55}))
			c.Request.Header.Add("Content-Type", test.contentType)

			r.Use(func(c *gin.Context) {
				c.Set(test.userIdHeaderName, test.userIdHeaderValue)
			})
			r.POST("/upload", handlers.UploadFile)

			// Make Request
			r.ServeHTTP(w, c.Request)

			// Assert
			if !assert.Equal(t, test.outStatusCode, w.Code) && !test.wantError {
				t.Fatalf("status code not equals")
			}

			if !assert.Equal(t, test.outBody, w.Body.String()) && !test.wantError {
				t.Fatalf("repsonse body not equals")
			}
		})
	}
}
