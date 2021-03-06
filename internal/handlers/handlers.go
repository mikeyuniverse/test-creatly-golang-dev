package handlers

import (
	"creatly-task/internal/models"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mock.go

type Hasher interface {
	Hash(password string) (string, error)
}

type Services interface {
	SignUp(user *models.UserSignUpInput) error
	SignIn(user *models.UserSignInInput) (string, error)
	Files() ([]models.FileOut, error)
	UploadFile(file *models.FileUploadInput) error
	ParseToken(token string) (string, error)
}

type Handlers struct {
	services        Services
	MaxSizeLimit    int // Bytes count
	hasher          Hasher
	tokenHeaderName string
	userHeaderName  string
}

func New(services Services, FileSizeLimit int, hasher Hasher, tokenHeaderName, userHeaderName string) *Handlers {
	return &Handlers{
		services:        services,
		MaxSizeLimit:    FileSizeLimit,
		hasher:          hasher,
		tokenHeaderName: tokenHeaderName,
		userHeaderName:  userHeaderName,
	}
}

func (h *Handlers) SignUp(c *gin.Context) {
	var input models.UserSignUpInput

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, textToMap("invalid input"))
		return
	}

	input.Password, err = h.hasher.Hash(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error while encrypting password"))
		return
	}

	err = h.services.SignUp(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error while creating an account"))
		return
	}

	c.JSON(http.StatusOK, textToMap("success"))
}

func (h *Handlers) SignIn(c *gin.Context) {
	var user models.UserSignInInput

	fmt.Printf("%+v\n", user)

	err := c.BindJSON(&user)
	if err != nil || (user.Email == "" || user.PasswordHash == "") {
		fmt.Printf("\n\n%+v\n", user)
		c.JSON(http.StatusBadRequest, textToMap("invalid credentials"))
		return
	}

	user.PasswordHash, err = h.hasher.Hash(user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error while encrypting password"))
		return
	}

	token, err := h.services.SignIn(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, textToMap("invalid creds"))
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, map[string]string{"token": token}) // Additional return token in JSON response
}

func (h *Handlers) AuthMiddleware(c *gin.Context) {

	token, err := h.getTokenFromHeader(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := h.services.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, textToMap(err.Error()))
		return
	}

	c.Set(h.userHeaderName, userID)
}

func (h *Handlers) Files(c *gin.Context) {
	files, err := h.services.Files()
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error getting file data"))
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *Handlers) UploadFile(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, textToMap("invalid content-type"))
		return
	}

	userIdValue := c.Keys[h.userHeaderName]
	if userIdValue == nil {
		c.JSON(http.StatusUnauthorized, textToMap("userID not found"))
		return
	}

	userID := c.Keys[h.userHeaderName].(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, textToMap("invalid userID"))
		return
	}

	body, err := readBody(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error while read body"))
		return
	}

	filename := fmt.Sprintf("%s-%d.png", userID, time.Now().Unix())
	filesize := c.Request.ContentLength

	if filesize >= int64(h.MaxSizeLimit) {
		c.JSON(http.StatusBadRequest, textToMap("file to large"))
		return
	}

	err = h.services.UploadFile(&models.FileUploadInput{
		Filename: filename,
		Size:     filesize,
		UserId:   userID,
		FileData: body,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error with upload file"))
		return
	}

	c.JSON(http.StatusOK, textToMap("upload success"))
}

func (h *Handlers) getTokenFromHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(h.tokenHeaderName)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid auth header value")
	}

	if headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header subject")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("invalid auth header token")
	}

	return headerParts[1], nil
}

func textToMap(text string) map[string]string {
	return map[string]string{"message": text}
}

func readBody(body io.ReadCloser) ([]byte, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
