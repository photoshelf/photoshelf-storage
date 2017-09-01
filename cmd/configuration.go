package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/labstack/gommon/log"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Type string
		Path string
	}
}

func load() (*Configuration, error) {
	configurationFile, err := ioutil.ReadFile("./application.yml")
	if err != nil {
		log.Warn(err)
	}

	instance := &Configuration{}
	if err := yaml.Unmarshal(configurationFile, instance); err != nil {
		return nil, err
	}
	instance.parse()

	return instance, nil
}

func (configuration *Configuration) parse() {
	configuration.Server.Port = *flag.Int("p", configuration.Server.Port, "port number")
	configuration.Storage.Type = *flag.String("t", configuration.Storage.Type, "storage type [file|leveldb|boltdb]")
	configuration.Storage.Path = *flag.String("d", configuration.Storage.Path, "storage path")
	flag.Parse()
}

func configure() (*Configuration, error) {
	configuration, err := load()
	if err != nil {
		return nil, err
	}

	var repository model.Repository
	switch configuration.Storage.Type {
	case "file":
		repository = datastore.NewFileStorage(configuration.Storage.Path)
	case "leveldb":
		repository, err = datastore.NewLeveldbStorage(configuration.Storage.Path)
		if err != nil {
			return nil, err
		}
	case "boltdb":
		repository, err = datastore.NewBoltdbStorage(configuration.Storage.Path)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("unknown storage type : %s", configuration.Storage.Type))
	}

	photoController := new(controller.PhotoController)
	if err := inject.Populate(photoController, new(service.PhotoService), repository); err != nil {
		return nil, err
	}

	return configuration, nil
}
