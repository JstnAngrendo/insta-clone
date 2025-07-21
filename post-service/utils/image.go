package utils

import (
	"strings"

	"github.com/disintegration/imaging"
)

func GenerateThumbnail(imagePath string) (string, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", err
	}
	thumb := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	thumbPath := strings.Replace(imagePath, ".jpg", "_thumb.jpg", 1)
	err = imaging.Save(thumb, thumbPath)
	if err != nil {
		return "", err
	}
	return thumbPath, nil
}
