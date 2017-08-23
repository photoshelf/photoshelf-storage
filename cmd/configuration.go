package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"flag"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore"
	"fmt"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/labstack/gommon/log"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"errors"
)

type configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Type string
		Path string
	}
}

func (configuration *configuration) parse() {
	configuration.Server.Port = *flag.Int("p", configuration.Server.Port, "port number")
	configuration.Storage.Type = *flag.String("t", configuration.Storage.Type, "storage type [file|leveldb|boltdb]")
	configuration.Storage.Path = *flag.String("d", configuration.Storage.Path, "storage path")
	flag.Parse()
}

func configure() (*configuration, error) {
	configurationFile, err := ioutil.ReadFile("./application.yml")
	if err != nil {
		log.Warn(err)
	}

	configuration := &configuration{}
	if err := yaml.Unmarshal(configurationFile, configuration); err != nil {
		return nil, err
	}
	configuration.parse()

	var repository model.Repository
	switch configuration.Storage.Type {
	case "file":
		repository = datastore.NewFileStorage(configuration.Storage.Type)
	case "leveldb":
		repository, err = datastore.NewLeveldbStorage(configuration.Storage.Type)
		if err != nil {
			return nil, err
		}
	case "boltdb":
		repository, err = datastore.NewBoltdbStorage(configuration.Storage.Type)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("unknown storage type : %s", configuration.Storage.Type))
	}
	photoService := service.NewPhotoService(repository)

	container.Set(controller.NewPhotoController(*photoService))

	return configuration, nil
}
