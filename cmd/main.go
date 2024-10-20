package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")

	s3Client, err := minio.New("s3.amazonaws.com", &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKey, "Y9b4/t84BT+mvCcVgAGPcpzV53TxoRpzHHQGrfKH", ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()

	e.GET("/images", func(c echo.Context) error {
		bucketName := "aws-da-abhinav"
		imageURLs := []string{}

		objectCh := s3Client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error listing objects: %v", object.Err))
			}
			if strings.HasSuffix(object.Key, ".jpg") || strings.HasSuffix(object.Key, ".jpeg") || strings.HasSuffix(object.Key, ".png") {
				imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, object.Key)
				imageURLs = append(imageURLs, imageURL)
			}
		}

		return c.JSON(http.StatusOK, imageURLs)
	})

	e.POST("/upload", func(c echo.Context) error {
		bucketName := "aws-da-abhinav"

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

		_, err = s3Client.PutObject(context.Background(), bucketName, objectName, src, file.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error uploading file: %v", err))
		}

		return c.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", objectName))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
