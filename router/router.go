package router

import (
	"github.com/facebookgo/inject"
	"github.com/labstack/echo"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
)

func Load() (*echo.Echo, error) {
	e := echo.New()

	var graph inject.Graph
	var photoController controller.PhotoController
	if err := graph.Provide(
		&inject.Object{Value: &photoController},
	); err != nil {
		return nil, err
	}

	g := e.Group("photos")
	g.GET("/:id", photoController.Get)
	g.POST("/", photoController.Post)
	g.PUT("/:id", photoController.Put)
	g.DELETE("/:id", photoController.Delete)

	return e, nil
}
