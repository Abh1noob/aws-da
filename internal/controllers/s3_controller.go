package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

type FileController struct {
	S3Client   *minio.Client
	BucketName string
	DB         *sql.DB
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

	userEmail, err := getEmailFromJWT(c)

	if err != nil {
		fmt.Println("Error getting email from JWT:", err)
		return c.JSON(http.StatusUnauthorized, "Invalid token in S3")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error getting uploaded file: %v", err))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error opening file: %v", err))
	}
	defer src.Close()

	fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", fc.BucketName, fileName)

	_, err = fc.S3Client.PutObject(c.Request().Context(), fc.BucketName, fileName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error uploading file to S3: %v", err))
	}

	isVisible := c.FormValue("is_visible") == "true"

	_, err = fc.DB.Exec("INSERT INTO posts (email, image_url, is_visible) VALUES ($1, $2, $3)", userEmail, imageURL, isVisible)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error saving file info to database: %v", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", fileName))
}

func getEmailFromJWT(c echo.Context) (string, error) {

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("missing token")
	}

	tokenString = strings.Split(tokenString, "Bearer ")[1]
	fmt.Println("Received Token: ", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("email not found in token")
	}

	return email, nil
}

func (fc *FileController) ListPublicImages(c echo.Context) error {
	var images []map[string]string

	rows, err := fc.DB.Query("SELECT username, image_url FROM posts join users on posts.email=users.email WHERE is_visible = true")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching public images: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var username, imageURL string
		if err := rows.Scan(&username, &imageURL); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error scanning row: %v", err))
		}

		images = append(images, map[string]string{
			"username": username,
			"image":    imageURL,
		})
	}

	return c.JSON(http.StatusOK, images)
}

func (fc *FileController) ListPrivateImages(c echo.Context) error {

	userEmail, err := getEmailFromJWT(c)
	if err != nil {
		fmt.Println("Error fetching user email:", err)
		return c.String(http.StatusUnauthorized, "Invalid token for private")
	}

	fmt.Println("User email:", userEmail)

	rows, err := fc.DB.Query(
		"SELECT username, image_url FROM posts JOIN users ON posts.email = users.email WHERE is_visible = false AND posts.email = $1",
		userEmail,
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching private images: %v", err))
	}
	defer rows.Close()

	var images []map[string]string
	for rows.Next() {
		var username, imageURL string
		if err := rows.Scan(&username, &imageURL); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error scanning row: %v", err))
		}

		images = append(images, map[string]string{
			"username": username,
			"image":    imageURL,
		})
	}

	return c.JSON(http.StatusOK, images)
}
