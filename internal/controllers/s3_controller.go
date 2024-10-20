package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

type FileController struct {
	S3Client   *minio.Client
	BucketName string
}

func (fc *FileController) ListImages(c echo.Context) error {
	imageURLs := []string{}

	objectCh := fc.S3Client.ListObjects(context.Background(), fc.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error listing objects: %v", object.Err))
		}
		if strings.HasSuffix(object.Key, ".jpg") || strings.HasSuffix(object.Key, ".jpeg") || strings.HasSuffix(object.Key, ".png") {
			imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", fc.BucketName, object.Key)
			imageURLs = append(imageURLs, imageURL)
		}
	}

	return c.JSON(http.StatusOK, imageURLs)
}

func (fc *FileController) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error getting uploaded file: %v", err))
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error opening uploaded file: %v", err))
	}
	defer src.Close()

	objectName := file.Filename
	contentType := file.Header.Get("Content-Type")

	_, err = fc.S3Client.PutObject(context.Background(), fc.BucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error uploading file: %v", err))
	}

	return c.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", objectName))
}
