package main

import (
	"aws-da/internal/controllers"
	customMiddleware "aws-da/internal/middlewares"

	"aws-da/internal/routes"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	s3Client, err := minio.New(os.Getenv("S3_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	e := echo.New()

	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	authController := &controllers.AuthController{DB: db}
	fileController := &controllers.FileController{
		S3Client:   s3Client,
		BucketName: os.Getenv("S3_BUCKET_NAME"),
		DB:         db, // Pass DB connection here
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtMiddleware := customMiddleware.JWTMiddleware(jwtSecret)

	routes.RegisterAuthRoutes(e, authController)
	routes.RegisterS3Routes(e, fileController, jwtMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
