package server

import (
	"creatly-task/internal/models"
	"creatly-task/internal/services"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	MaxSizeLimit    int64
	hasher          Hasher
	tokenHeaderName string
}

func NewHandlers(services *services.Services, FileSizeLimit int64, hasher Hasher, tokenHeaderName string) *Handlers {
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
	}

	h.services.ParseToken(token)

	// err := h.services.ParseToken()

	// token, err := h.checkTokenFromHeader(c)
	// if err != nil {
	// 	// TODO Change Abort type (delete error text)
	// 	c.AbortWithError(http.StatusUnauthorized, err)
	// 	return
	// }

	// 	userId, err := h.services.GetUserIdByToken(token)
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, textToMap("unauthorized"))
	// 		return
	// 	}

	// 	c.Set("userID", userId) // Pass userId in context

	c.Next()
}

func (h *Handlers) Files(c *gin.Context) {
	// files, err := h.services.Files()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, textToMap("error getting file data"))
	// 	return
	// }

	// c.JSON(http.StatusOK, files)
	c.JSON(http.StatusOK, textToMap("access to closed functionality"))
}

func (h *Handlers) UploadFile(c *gin.Context) {
	userID := c.GetInt("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, textToMap("unknown userID"))
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, textToMap("file not found"))
		return
	}

	filename := header.Filename
	// TODO Check filename and format of file

	size := header.Size
	if size >= h.MaxSizeLimit {
		fmt.Printf("Want filesize - %d\nAccept filesize - %d\n", h.MaxSizeLimit, size)
		c.JSON(http.StatusBadRequest, textToMap("file too large"))
		return
	}

	err = h.services.UploadFile(&models.FileUploadInput{
		Filename: filename,
		Size:     size,
		UserId:   userID,
		FileData: file,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error while saving the file"))
		return
	}

	c.JSON(http.StatusOK, textToMap("upload success"))
}

func textToMap(text string) map[string]string {
	return map[string]string{"message": text}
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
