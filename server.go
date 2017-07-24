package main

import (
	"flag"
	"fmt"
	"github.com/duck8823/photoshelf-storage/service"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/duck8823/photoshelf-storage/infrastructure"
	"github.com/duck8823/photoshelf-storage/presentation/controller"
)

type Configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Directory string
	}
}

func main() {
	configurationFile, err := ioutil.ReadFile("./application.yml")
	if err != nil {
		log.Warn(err)
	}

	configuration := &Configuration{}
	if err := yaml.Unmarshal(configurationFile, configuration); err != nil {
		log.Fatal(err)
		return
	}

	repository := infrastructure.NewFileStorage(configuration.Storage.Directory)
	photoService := service.NewPhotoService(repository)
	photoController := controller.NewPhotoController(*photoService)

	e := echo.New()
	e.GET("/:id", photoController.Get)
	e.POST("/", photoController.Post)
	e.PUT("/:id", photoController.Put)
	e.DELETE("/:id", photoController.Delete)

	port := flag.Int("port", configuration.Server.Port, "port number")
	flag.Parse()

	address := fmt.Sprintf(":%d", *port)
	e.Logger.Debug(e.Start(address))
}
