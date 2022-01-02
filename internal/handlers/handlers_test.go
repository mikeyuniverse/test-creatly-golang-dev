package handlers

import (
	"bytes"
	mock_handlers "creatly-task/internal/handlers/mocks"
	"creatly-task/internal/models"
	"errors"
	"net/http/httptest"
	"testing"

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

func Test_SignIn(t *testing.T) {}

func Test_AuthMiddleware(t *testing.T) {}

func Test_Files(t *testing.T) {}

func Test_UploadFile(t *testing.T) {}

func Test_getTokenFromHeader(t *testing.T) {}

func Test_textToMap(t *testing.T) {}

func Test_readBody(t *testing.T) {}
