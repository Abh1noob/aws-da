package main

import (
	controller "aws-da/internal/controllers"
	"aws-da/internal/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Load environment variables from .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	s3BucketName := os.Getenv("S3_BUCKET_NAME")

	if s3Endpoint == "" {
		log.Fatalln("S3_ENDPOINT environment variable is not set")
	}

	s3Client, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKey, s3SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	e := echo.New()

	fileController := &controller.FileController{
		S3Client:   s3Client,
		BucketName: s3BucketName,
	}

	routes.RegisterRoutes(e, fileController)

	// Start the Echo server
	e.Logger.Fatal(e.Start(":8080"))
}
