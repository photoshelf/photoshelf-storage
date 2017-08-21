package main

import (
	"flag"
	"fmt"
	"github.com/photoshelf/photoshelf-storage/infrastructure"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore"
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
	storageType := flag.String("t", configuration.Storage.Type, "storage type [file|leveldb|boltdb]")
	storagePath := flag.String("d", configuration.Storage.Path, "storage path")
	flag.Parse()

	var repository model.Repository
	switch *storageType {
	case "file":
		repository = datastore.NewFileStorage(*storagePath)
	case "leveldb":
		repository, err = datastore.NewLeveldbStorage(*storagePath)
		if err != nil {
			log.Fatal(err)
			return
		}
	case "boltdb":
		repository, err = datastore.NewBoltdbStorage(*storagePath)
		if err != nil {
			log.Fatal(err)
			return
		}
	default:
		log.Fatal(fmt.Sprintf("unknown storage type : %s", *storageType))
		return
	}
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
