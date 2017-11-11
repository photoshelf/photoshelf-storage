package router

import (
	"github.com/labstack/echo"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
)

func Load() (*echo.Echo, error) {
	e := echo.New()

	photoController := controller.New()
	container.Get(&photoController)

	g := e.Group("photos")
	g.GET("/:id", photoController.Get)
	g.POST("/", photoController.Post)
	g.PUT("/:id", photoController.Put)
	g.DELETE("/:id", photoController.Delete)

	return e, nil
}
