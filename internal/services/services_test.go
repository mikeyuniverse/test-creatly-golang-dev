package services

import (
	"creatly-task/internal/models"
	"creatly-task/internal/repo"
	mock_repo "creatly-task/internal/repo/mocks"
	mock_services "creatly-task/internal/services/mocks"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func newServices(t *testing.T) *Services {
// 	return services
// }

func Test_SignUp(t *testing.T) {
	testTable := []struct {
		name      string
		expect    error
		input     models.UserSignUpInput
		behavior  func(*mock_repo.MockUsers)
		wantError bool
	}{
		{
			name:   "OK",
			expect: nil,
			input: models.UserSignUpInput{
				Email:    "some@mail.com",
				Password: "SuperStrongPassword",
			},
			behavior: func(mu *mock_repo.MockUsers) {
				mu.EXPECT().CreateUser(&models.UserSignUpInput{
					Email:    "some@mail.com",
					Password: "SuperStrongPassword",
				}).Return(nil)
			},
			wantError: false,
		},
		{
			name:   "ERROR: user not created, error in db",
			expect: errors.New("database error"),
			input: models.UserSignUpInput{
				Email:    "some@mail.com",
				Password: "SuperStrongPassword",
			},
			behavior: func(mu *mock_repo.MockUsers) {
				mu.EXPECT().CreateUser(&models.UserSignUpInput{
					Email:    "some@mail.com",
					Password: "SuperStrongPassword",
				}).Return(errors.New("database error"))
			},
			wantError: true,
		},
	}

	for _, test := range testTable {

		ctrl := gomock.NewController(t)
		usersRepo := mock_repo.NewMockUsers(ctrl)
		tokenRepo := mock_repo.NewMockTokens(ctrl)
		filesRepo := mock_repo.NewMockFiles(ctrl)
		repo := &repo.Repo{
			Users:  usersRepo,
			Tokens: tokenRepo,
			Files:  filesRepo,
		}
		tokens := mock_services.NewMockTokener(ctrl)
		cloud := mock_services.NewMockCloudStorage(ctrl)

		test.behavior(usersRepo)

		services := New(repo, tokens, cloud)

		err := services.SignUp(&test.input)
		if err != nil && err != test.expect && !test.wantError {
			t.Fatalf("error service SignUp - %s\n", err.Error())
		}
	}

}

func Test_SignIn(t *testing.T) {
	testTable := []struct {
		name      string
		input     models.UserSignInInput
		behavior  func(*mock_repo.MockUsers, *mock_services.MockTokener)
		wantError bool
		outToken  string
	}{
		{
			name: "OK",
			input: models.UserSignInInput{
				Email:        "some@mail.com",
				PasswordHash: "wd781bpi2du08237f82v",
			},
			behavior: func(mu *mock_repo.MockUsers, mt *mock_services.MockTokener) {
				mu.EXPECT().GetUserByCreds("some@mail.com").Return(&models.UserSignInOutput{
					UserID:   primitive.ObjectID{53, 50, 51, 52, 53, 54, 50, 56, 57, 58, 49},
					Email:    "some@mail.com",
					Password: "wd781bpi2du08237f82v",
				}, nil)
				mt.EXPECT().GenerateToken(primitive.ObjectID{53, 50, 51, 52, 53, 54, 50, 56, 57, 58, 49}.String()).Return("token", nil)
			},
			wantError: false,
			outToken:  "token",
		},
		{
			name: "ERROR: returned invalid token",
			input: models.UserSignInInput{
				Email:        "some@mail.com",
				PasswordHash: "wd781bpi2du08237f82v",
			},
			behavior: func(mu *mock_repo.MockUsers, mt *mock_services.MockTokener) {
				mu.EXPECT().GetUserByCreds("some@mail.com").Return(&models.UserSignInOutput{
					UserID:   primitive.ObjectID{53, 50, 51, 52, 53, 54, 50, 56, 57, 58, 49},
					Email:    "some@mail.com",
					Password: "wd781bpi2du08237f82v",
				}, nil)
				mt.EXPECT().GenerateToken(primitive.ObjectID{53, 50, 51, 52, 53, 54, 50, 56, 57, 58, 49}.String()).Return("", nil) // Here error
			},
			wantError: true,
			outToken:  "token",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersRepo := mock_repo.NewMockUsers(ctrl)
			tokenRepo := mock_repo.NewMockTokens(ctrl)
			filesRepo := mock_repo.NewMockFiles(ctrl)
			repo := &repo.Repo{
				Users:  usersRepo,
				Tokens: tokenRepo,
				Files:  filesRepo,
			}
			tokens := mock_services.NewMockTokener(ctrl)
			cloud := mock_services.NewMockCloudStorage(ctrl)

			test.behavior(usersRepo, tokens)

			services := New(repo, tokens, cloud)

			token, err := services.SignIn(&test.input)
			if err != nil && !test.wantError {
				t.Fatalf("SignIn error - %s\n", err.Error())
			}

			if test.outToken != token && !test.wantError {
				t.Fatal("unexpected token")
			}

		})
	}
}

func Test_Files(t *testing.T) {
	testTable := []struct {
		name      string
		behavior  func(*mock_repo.MockFiles)
		wantError bool
		outFiles  []models.FileOut
	}{
		{
			name: "OK",
			behavior: func(mf *mock_repo.MockFiles) {
				mf.EXPECT().All().Return([]models.FileOut{
					{
						Filename: "file 1",
						Size:     100,
						Date:     19236328,
						UserId:   "1",
						Url:      "https://s3.storage.com/123/1",
					},
				}, nil)
			},
			wantError: false,
			outFiles: []models.FileOut{
				{
					Filename: "file 1",
					Size:     100,
					Date:     19236328,
					UserId:   "1",
					Url:      "https://s3.storage.com/123/1",
				},
			},
		},
		{
			name: "ERROR: error in Files.All()",
			behavior: func(mf *mock_repo.MockFiles) {
				mf.EXPECT().All().Return([]models.FileOut{}, errors.New("some error"))
			},
			wantError: true,
			outFiles:  []models.FileOut{},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersRepo := mock_repo.NewMockUsers(ctrl)
			tokenRepo := mock_repo.NewMockTokens(ctrl)
			filesRepo := mock_repo.NewMockFiles(ctrl)
			repo := &repo.Repo{
				Users:  usersRepo,
				Tokens: tokenRepo,
				Files:  filesRepo,
			}
			tokens := mock_services.NewMockTokener(ctrl)
			cloud := mock_services.NewMockCloudStorage(ctrl)

			test.behavior(filesRepo)

			services := New(repo, tokens, cloud)

			files, err := services.Files()

			if err != nil && !test.wantError {
				t.Fatalf("Service Files error - %s\n", err.Error())
			}

			if !reflect.DeepEqual(test.outFiles, files) && !test.wantError {
				t.Fatalf("files not equals\nReceived - %+v\nWant - %+v\n", files, test.outFiles)
			}

		})
	}
}

func Test_UploadFile(t *testing.T) {
	testTable := []struct {
		name        string
		behavior    func(*mock_services.MockCloudStorage, *mock_repo.MockFiles)
		wantError   bool
		inputUpload models.FileUploadInput
	}{
		{
			name: "OK",
			behavior: func(mcs *mock_services.MockCloudStorage, mf *mock_repo.MockFiles) {
				mcs.EXPECT().UploadFile([]byte{}, int64(10000), "file1.png").Return("https://s3.storage.com/1", nil)
				mf.EXPECT().AddLog(&models.FileUploadLogInput{
					Size:       10000,
					UploadDate: time.Now().Unix(),
					Filename:   "file1.png",
					UserId:     "1",
					Url:        "https://s3.storage.com/1",
				}).Return(nil)
			},
			wantError: false,
			inputUpload: models.FileUploadInput{
				FileData: []byte{},
				Size:     10000,
				Filename: "file1.png",
				UserId:   "1",
			},
		},
		{
			name: "ERROR: upload error",
			behavior: func(mcs *mock_services.MockCloudStorage, mf *mock_repo.MockFiles) {
				mcs.EXPECT().UploadFile([]byte{}, int64(60000000), "file1.png").Return("", errors.New("uploading error"))
			},
			wantError: true,
			inputUpload: models.FileUploadInput{
				FileData: []byte{},
				Size:     60000000,
				Filename: "file1.png",
				UserId:   "1",
			},
		},
		{
			name: "ERROR: add log error",
			behavior: func(mcs *mock_services.MockCloudStorage, mf *mock_repo.MockFiles) {
				mcs.EXPECT().UploadFile([]byte{}, int64(60000000), "file1.png").Return("https://s3.storage.com/1", nil)
				mf.EXPECT().AddLog(&models.FileUploadLogInput{
					Size:       60000000,
					UploadDate: time.Now().Unix(),
					Filename:   "file1.png",
					UserId:     "1",
					Url:        "https://s3.storage.com/1",
				}).Return(errors.New("add log error"))
			},
			wantError: true,
			inputUpload: models.FileUploadInput{
				FileData: []byte{},
				Size:     60000000,
				Filename: "file1.png",
				UserId:   "1",
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersRepo := mock_repo.NewMockUsers(ctrl)
			tokenRepo := mock_repo.NewMockTokens(ctrl)
			filesRepo := mock_repo.NewMockFiles(ctrl)
			repo := &repo.Repo{
				Users:  usersRepo,
				Tokens: tokenRepo,
				Files:  filesRepo,
			}
			tokens := mock_services.NewMockTokener(ctrl)
			cloud := mock_services.NewMockCloudStorage(ctrl)

			test.behavior(cloud, filesRepo)

			services := New(repo, tokens, cloud)

			err := services.UploadFile(&test.inputUpload)

			if err != nil && !test.wantError {
				t.Fatalf("Service UploadFile error - %s\n", err.Error())
			}
		})
	}
}

func Test_ParseToken(t *testing.T) {
	testTable := []struct {
		name       string
		behavior   func(*mock_services.MockTokener)
		inputToken string
		outUserId  string
		wantError  bool
	}{
		{
			name: "OK",
			behavior: func(mt *mock_services.MockTokener) {
				mt.EXPECT().ParseToken("293o89bcuwp8yb0823peob2pf9u829p").Return("1", nil)
			},
			inputToken: "293o89bcuwp8yb0823peob2pf9u829p",
			outUserId:  "1",
			wantError:  false,
		},
		{
			name: "ERROR: parse error",
			behavior: func(mt *mock_services.MockTokener) {
				mt.EXPECT().ParseToken("whooohooo").Return("", errors.New("isn't token"))
			},
			inputToken: "whooohooo",
			outUserId:  "",
			wantError:  true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			usersRepo := mock_repo.NewMockUsers(ctrl)
			tokenRepo := mock_repo.NewMockTokens(ctrl)
			filesRepo := mock_repo.NewMockFiles(ctrl)

			repo := &repo.Repo{
				Users:  usersRepo,
				Tokens: tokenRepo,
				Files:  filesRepo,
			}

			tokens := mock_services.NewMockTokener(ctrl)
			cloud := mock_services.NewMockCloudStorage(ctrl)

			test.behavior(tokens)

			services := New(repo, tokens, cloud)

			userID, err := services.ParseToken(test.inputToken)
			if err != nil && !test.wantError {
				t.Fatalf("Service ParseToken error - %s\n", err.Error())
			}

			if userID != test.outUserId && !test.wantError {
				t.Fatalf("Invalid userID\nReceived - %s\nWant - %s\n", userID, test.outUserId)
			}
		})
	}
}
