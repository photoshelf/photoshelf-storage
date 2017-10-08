package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"flag"
	"path"
	"os"
)

func TestCanLoadDefault(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{os.Args[0]}

	configuration, err := load()
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1213, configuration.Server.Port)
		assert.EqualValues(t, "leveldb", configuration.Storage.Type)
		assert.EqualValues(t, "./photos", configuration.Storage.Path)
	}
}

func TestCanLoadFromFile(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{os.Args[0]}

	configurationPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata", "test.yml")
	os.Args = append(os.Args, "-c", configurationPath)

	configuration, err := load()
	if assert.NoError(t, err) {
		assert.EqualValues(t, configuration.Server.Port, 12345)
		assert.EqualValues(t, configuration.Storage.Type, "hoge")
		assert.EqualValues(t, configuration.Storage.Path, "fuga")
	}
}

func TestCanParseFlags(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{os.Args[0]}
	os.Args = append(os.Args, "-p", "54321", "-t", "foo", "-s", "bar")

	configuration, err := load()
	if assert.NoError(t, err) {
		assert.EqualValues(t, configuration.Server.Port, 54321)
		assert.EqualValues(t, configuration.Storage.Type, "foo")
		assert.EqualValues(t, configuration.Storage.Path, "bar")
	}
}
