package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type AuthController struct {
	DB *sql.DB
}

func (ac *AuthController) Signup(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid input: %v", err))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to hash password: %v", err))
	}

	_, err = ac.DB.Exec("INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", user.Email, string(hashedPassword), user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
	}

	return c.JSON(http.StatusCreated, "User created successfully")
}

func (ac *AuthController) Login(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	var hashedPassword string
	err := ac.DB.QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&hashedPassword)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{"access_token": tokenString})
}
