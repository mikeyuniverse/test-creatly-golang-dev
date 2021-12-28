package server

import (
	"creatly-task/internal/models"
	"creatly-task/internal/services"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const COOKIE_NAME = "token"
const HEADER_WITH_TOKEN = "Authorization"
const SIGNING_KEY = "aldnupp12nd2no"

type Hasher interface {
	Hash(password string) (string, error)
}

type Handlers struct {
	services        *services.Services
	MaxSizeLimit    int // Bytes
	hasher          Hasher
	tokenHeaderName string
	userHeaderName  string
}

func NewHandlers(services *services.Services, FileSizeLimit int, hasher Hasher, tokenHeaderName, userHeaderName string) *Handlers {
	return &Handlers{
		services:        services,
		MaxSizeLimit:    FileSizeLimit,
		hasher:          hasher,
		tokenHeaderName: tokenHeaderName,
	}
}

func (h *Handlers) SignUp(c *gin.Context) {
	var input models.UserSignUpInput

	err := c.BindJSON(&input)
	if err != nil {
		fmt.Printf("\n\n%+v\n", input)
		c.JSON(http.StatusBadRequest, textToMap("invalid input"))
		return
	}

	input.Password, err = h.hasher.Hash(input.Password)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		c.JSON(http.StatusInternalServerError, textToMap("error while encrypting password"))
		return
	}

	err = h.services.SignUp(&input)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		c.JSON(http.StatusInternalServerError, textToMap("error while creating an account"))
		return
	}

	c.JSON(http.StatusOK, textToMap("success"))
}

func (h *Handlers) SignIn(c *gin.Context) {
	var user models.UserSignInInput

	err := c.Bind(&user)
	if err != nil {
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
	// Check header Content-Type. Must be "image/jpeg" or "image/png"
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, textToMap("invalid content-type"))
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
	bytes, err := io.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
