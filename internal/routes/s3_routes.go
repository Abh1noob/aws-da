package routes

import (
	controller "aws-da/internal/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, fileController *controller.FileController) {
	e.GET("/images", fileController.ListImages)
	e.POST("/upload", fileController.UploadFile)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	})
}
