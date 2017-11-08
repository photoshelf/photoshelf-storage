package application

import (
	"errors"
	"flag"
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/domain/model"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/boltdb_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/file_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/leveldb_storage"
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

func load() *Configuration {
	configuration := &Configuration{}

	flag.Var(configuration, "c", "configuration file path")
	flag.IntVar(
		&configuration.Server.Port,
		"p",
		1213,
		"port number",
	)
	flag.StringVar(
		&configuration.Storage.Type,
		"t",
		"leveldb",
		"storage type [file|leveldb|boltdb]",
	)
	flag.StringVar(
		&configuration.Storage.Path,
		"s",
		"./photos",
		"storage path",
	)
	flag.Parse()

	return configuration
}

func Configure() (*Configuration, error) {
	configuration := load()

	var repository model.Repository
	var err error
	switch configuration.Storage.Type {
	case "file":
		repository = file_storage.NewFileStorage(configuration.Storage.Path)
	case "leveldb":
		repository, err = leveldb_storage.NewLeveldbStorage(configuration.Storage.Path)
		if err != nil {
			return nil, err
		}
	case "boltdb":
		repository, err = boltdb_storage.NewBoltdbStorage(configuration.Storage.Path)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("unknown storage type : %s", configuration.Storage.Type))
	}

	photoController := new(controller.PhotoController)
	if err := inject.Populate(photoController, service.New(), repository); err != nil {
		return nil, err
	}
	container.Set(photoController)

	return configuration, nil
}
