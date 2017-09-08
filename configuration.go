package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
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

func (configuration *Configuration) String() string {
	if (Configuration{}) == *configuration {
		return ""
	}
	return fmt.Sprint(*configuration)
}

func (configuration *Configuration) Set(path string) error {
	configurationFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(configurationFile, configuration); err != nil {
		return err
	}
	return nil
}

var defaultConf = &Configuration{
	Server: struct {
		Port int
	}{
		1213,
	},
	Storage: struct {
		Type string
		Path string
	}{
		"leveldb",
		"./photos",
	},
}

func load() (*Configuration, error) {
	configuration := &Configuration{}

	flag.Var(configuration, "c", "configuration file path")
	flag.IntVar(
		&configuration.Server.Port,
		"p",
		defaultConf.Server.Port,
		"port number",
	)
	flag.StringVar(
		&configuration.Storage.Type,
		"t",
		defaultConf.Storage.Type,
		"storage type [file|leveldb|boltdb]",
	)
	flag.StringVar(
		&configuration.Storage.Path,
		"s",
		defaultConf.Storage.Path,
		"storage path",
	)
	flag.Parse()

	return configuration, nil
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
	container.Set(*photoController)

	return configuration, nil
}
