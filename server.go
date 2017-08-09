package main

import (
	"flag"
	"fmt"
	"github.com/photoshelf/photoshelf-storage/infrastructure"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/photoshelf/photoshelf-storage/service"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	configurationFile, err := ioutil.ReadFile("./application.yml")
	if err != nil {
		log.Warn(err)
	}

	configuration := &infrastructure.Configuration{}
	if err := yaml.Unmarshal(configurationFile, configuration); err != nil {
		log.Fatal(err)
		return
	}

	port := flag.Int("p", configuration.Server.Port, "port number")
	imageDir := flag.String("d", configuration.Storage.Directory, "storage directory")
	flag.Parse()

	repository := infrastructure.NewFileStorage(*imageDir)
	photoService := service.NewPhotoService(repository)
	photoController := controller.NewPhotoController(*photoService)

	e := echo.New()

	g := e.Group("photos")
	g.GET("/:id", photoController.Get)
	g.POST("/", photoController.Post)
	g.PUT("/:id", photoController.Put)
	g.DELETE("/:id", photoController.Delete)

	address := fmt.Sprintf(":%d", *port)
	e.Logger.Debug(e.Start(address))
}
