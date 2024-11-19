package routes

import (
	"aws-da/internal/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Echo, authController *controllers.AuthController) {
	e.POST("/signup", authController.Signup)
	e.POST("/login", authController.Login)
}
