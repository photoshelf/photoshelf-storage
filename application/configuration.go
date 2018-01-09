package application

import (
	"flag"
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/boltdb_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/file_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/leveldb_storage"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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
	return yaml.Unmarshal(configurationFile, configuration)
}

func load(args ...string) *Configuration {
	configuration := &Configuration{}

	flg := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flg.Var(configuration, "c", "configuration file path")
	flg.IntVar(
		&configuration.Server.Port,
		"p",
		1323,
		"port number",
	)
	flg.StringVar(
		&configuration.Storage.Type,
		"t",
		"boltdb",
		"storage type [file|leveldb|boltdb]",
	)
	flg.StringVar(
		&configuration.Storage.Path,
		"s",
		"./photos",
		"storage path",
	)
	flg.Parse(args)

	return configuration
}

func Configure(args ...string) (*Configuration, error) {
	configuration := load(args...)

	var repository photo.Repository
	var err error
	switch configuration.Storage.Type {
	case "file":
		repository = file_storage.New(configuration.Storage.Path)
	case "leveldb":
		repository, err = leveldb_storage.New(configuration.Storage.Path)
		if err != nil {
			return nil, err
		}
	case "boltdb":
		repository, err = boltdb_storage.New(configuration.Storage.Path)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown storage type : %s", configuration.Storage.Type)
	}

	photoController := controller.New()
	if err := inject.Populate(photoController, service.New(), repository); err != nil {
		return nil, err
	}
	container.Set(photoController)

	return configuration, nil
}
