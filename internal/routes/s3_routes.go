package routes

import (
	controller "aws-da/internal/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, fileController *controller.FileController) {
	e.GET("/images", fileController.ListImages)
	e.POST("/upload", fileController.UploadFile)
}
