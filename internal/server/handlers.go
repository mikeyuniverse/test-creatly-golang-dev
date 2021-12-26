package server

import (
	"creatly-task/internal/models"
	"creatly-task/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const COOKIE_NAME = "token"

type Handlers struct {
	services     *services.Services
	MaxSizeLimit int64
}

func NewHandlers(services *services.Services, FileSizeLimit int64) *Handlers {
	return &Handlers{services: services, MaxSizeLimit: FileSizeLimit}
}

func (h *Handlers) SignUp(c *gin.Context) {
	var input models.UserSignUpInput

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, textToMap("invalid input"))
		return
	}

	err = h.services.SignUp(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, textToMap("error while creating an account"))
		return
	}

	c.JSON(http.StatusOK, textToMap("ok"))
}

func (h *Handlers) SignIn(c *gin.Context) {
	// TODO How auth?
	var user models.UserSignInInput

	c.Bind(&user)

	token, err := h.services.SignIn(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, textToMap("invalid creds"))
		return
	}

	fmt.Printf("SUCCESS: SignIn\nToken - %s\n", token)

}

func (h *Handlers) AuthMiddleware(c *gin.Context) {
	token, err := c.Cookie(COOKIE_NAME)
	if err != nil {
		c.JSON(http.StatusUnauthorized, textToMap("unauthorized"))
		return
	}

	// TODO Check: token is valid?

	// TODO Get userId by token
	userId, err := h.services.GetUserIdByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, textToMap("unauthorized"))
		return
	}

	// TODO Pass userId in context
	c.Set("userID", userId)

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
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, textToMap("file not found"))
		return
	}

	filename := header.Filename
	size := header.Size
	if size >= h.MaxSizeLimit {
		c.JSON(http.StatusBadRequest, textToMap("file too large"))
		return
	}

	err = h.services.UploadFile(&models.FileUploadInput{
		Filename: filename,
		Volume:   size,
		UserId:   123, // TODO Accept userID from context or by token
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
