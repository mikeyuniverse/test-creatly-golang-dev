package server

import (
	"creatly-task/internal/models"
	"creatly-task/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const COOKIE_NAME = "token"

type Hasher interface {
	Hash(password string) (string, error)
}

type Handlers struct {
	services     *services.Services
	MaxSizeLimit int64
	hasher       Hasher
}

func NewHandlers(services *services.Services, FileSizeLimit int64, hasher Hasher) *Handlers {
	return &Handlers{services: services, MaxSizeLimit: FileSizeLimit, hasher: hasher}
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

	err := c.Bind(&user)
	if err != nil {
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

	c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handlers) AuthMiddleware(c *gin.Context) {
	token, err := c.Cookie(COOKIE_NAME)
	if err != nil {
		c.JSON(http.StatusUnauthorized, textToMap("unauthorized"))
		return
	}

	// TODO Check: token is valid?

	userId, err := h.services.GetUserIdByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, textToMap("unauthorized"))
		return
	}

	c.Set("userID", userId) // Pass userId in context

	c.Next()
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
