package routes

import (
	"aws-da/internal/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterS3Routes(e *echo.Echo, fileController *controllers.FileController, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/images/public", fileController.ListPublicImages)
	e.GET("/images/private", fileController.ListPrivateImages, jwtMiddleware)
	e.POST("/upload", fileController.UploadFile, jwtMiddleware)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	}) // wrong place to keep this. i know.
}
