package router

import (
	"github.com/labstack/echo"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
)

func Load() *echo.Echo {
	e := echo.New()

	photoController := container.Get("PhotoController").(*controller.PhotoController)

	g := e.Group("photos")
	g.GET("/:id", photoController.Get)
	g.POST("/", photoController.Post)
	g.PUT("/:id", photoController.Put)
	g.DELETE("/:id", photoController.Delete)

	return e
}
