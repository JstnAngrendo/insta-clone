package storage

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context, file *multipart.FileHeader) (string, error) {
	dst := "uploads/" + file.Filename
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		return "", err
	}
	return "/static/" + file.Filename, nil
}
