package application

import (
	"flag"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/boltdb_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/file_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/leveldb_storage"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("with no args, can load default", func(t *testing.T) {
		resetFlag()

		configuration := load()
		assert.EqualValues(t, 1213, configuration.Server.Port)
		assert.EqualValues(t, "leveldb", configuration.Storage.Type)
		assert.EqualValues(t, "./photos", configuration.Storage.Path)
	})

	t.Run("with specify c flag, can load from file", func(t *testing.T) {
		resetFlag()

		configurationPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata", "test.yml")
		os.Args = append(os.Args, "-c", configurationPath)

		configuration := load()
		assert.EqualValues(t, configuration.Server.Port, 12345)
		assert.EqualValues(t, configuration.Storage.Type, "hoge")
		assert.EqualValues(t, configuration.Storage.Path, "fuga")
	})

	t.Run("with flags, can parse from flags", func(t *testing.T) {
		resetFlag()

		os.Args = append(os.Args, "-p", "54321", "-t", "foo", "-s", "bar")

		configuration := load()
		assert.EqualValues(t, configuration.Server.Port, 54321)
		assert.EqualValues(t, configuration.Storage.Type, "foo")
		assert.EqualValues(t, configuration.Storage.Path, "bar")
	})
}

func TestConfiguration_Set(t *testing.T) {
	t.Run("with no file, returns error", func(t *testing.T) {
		wrongPath := path.Join(os.TempDir(), "wrong_path")
		if err := os.RemoveAll(wrongPath); err != nil {
			t.Fatal(err)
		}

		conf := &Configuration{}
		err := conf.Set(wrongPath)
		assert.Error(t, err)
	})

	t.Run("with wrong data, returns error", func(t *testing.T) {
		wrongDataPath := path.Join(os.TempDir(), "wrong_data")
		if err := ioutil.WriteFile(wrongDataPath, []byte("This is not yml format"), 0700); err != nil {
			t.Fatal(err)
		}

		conf := &Configuration{}
		err := conf.Set(wrongDataPath)
		assert.Error(t, err)
	})
}

func TestConfiguration_String(t *testing.T) {
	t.Run("when not empty, returns value", func(t *testing.T) {
		resetFlag()

		conf := load()
		assert.NotEmpty(t, conf.String())
	})
}

func TestConfigure(t *testing.T) {
	t.Run("with leveldb type, returns instance specify", func(t *testing.T) {
		resetFlag()

		os.Args = append(os.Args, "-t", "leveldb")

		_, err := Configure()
		if assert.NoError(t, err) {
			assert.IsType(t, new(leveldb_storage.LeveldbStorage), actualRepository())
		}
	})

	t.Run("when fail to load leveldb, returns error", func(t *testing.T) {
		resetFlag()
		file, err := os.Create(path.Join(os.TempDir(), "readonly"))
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		os.Args = append(os.Args, "-t", "leveldb", "-s", path.Join(os.TempDir(), "readonly"))

		_, err = Configure()
		assert.Error(t, err)
	})

	t.Run("with file type, returns instance specify", func(t *testing.T) {
		resetFlag()

		os.Args = append(os.Args, "-t", "file")

		_, err := Configure()
		if assert.NoError(t, err) {
			assert.IsType(t, new(file_storage.FileStorage), actualRepository())
		}
	})

	t.Run("with boltdb type, returns instance specify", func(t *testing.T) {
		resetFlag()
		dbPath := path.Join(os.TempDir(), "boltdb")
		os.RemoveAll(dbPath)

		os.Args = append(os.Args, "-t", "boltdb", "-s", dbPath)

		_, err := Configure()
		if assert.NoError(t, err) {
			assert.IsType(t, new(boltdb_storage.BoltdbStorage), actualRepository())
		}
	})

	t.Run("when fail to load boltdb, returns error", func(t *testing.T) {
		resetFlag()
		dbPath := path.Join(os.TempDir(), "err_boltdb")
		os.RemoveAll(dbPath)
		os.MkdirAll(dbPath, 0600)

		os.Args = append(os.Args, "-t", "boltdb", "-s", dbPath)

		_, err := Configure()
		assert.Error(t, err)
	})

	t.Run("with unknown type, returns error", func(t *testing.T) {
		resetFlag()

		os.Args = append(os.Args, "-t", "unknown")

		_, err := Configure()
		assert.Error(t, err)
	})
}

func resetFlag() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{os.Args[0]}
}

func actualRepository() interface{} {
	var photoController controller.PhotoController
	container.Get(&photoController)

	return photoController.Service.(*service.PhotoServiceImpl).Repository
}
