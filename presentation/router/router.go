package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/photoshelf/photoshelf-storage/presentation/protobuf"
	"google.golang.org/grpc"
)

func LoadEchoServer() (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true

	photoController := controller.NewRestPhotoController()
	container.Get(&photoController)

	g := e.Group("photos")
	g.GET("/:id", photoController.Get)
	g.POST("/", photoController.Post)
	g.PUT("/:id", photoController.Put)
	g.DELETE("/:id", photoController.Delete)

	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit("20M"))

	return e, nil
}

func LoadGrpcServer() *grpc.Server {
	s := grpc.NewServer()

	photoServiceServer := controller.NewGrpcPhotoController()
	container.Get(&photoServiceServer)

	protobuf.RegisterPhotoServiceServer(s, photoServiceServer)

	return s
}
