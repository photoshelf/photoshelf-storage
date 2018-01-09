package application

import (
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/boltdb_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/file_storage"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore/leveldb_storage"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("with no args, can load default", func(t *testing.T) {
		configuration := load()
		assert.EqualValues(t, 1323, configuration.Server.Port)
		assert.EqualValues(t, "boltdb", configuration.Storage.Type)
		assert.EqualValues(t, "./photos", configuration.Storage.Path)
	})

	t.Run("with specify c flag, can load from file", func(t *testing.T) {
		configurationPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata", "test.yml")
		configuration := load("-c", configurationPath)
		assert.EqualValues(t, configuration.Server.Port, 12345)
		assert.EqualValues(t, configuration.Storage.Type, "hoge")
		assert.EqualValues(t, configuration.Storage.Path, "fuga")
	})

	t.Run("with flags, can parse from flags", func(t *testing.T) {
		configuration := load("-p", "54321", "-t", "foo", "-s", "bar")
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
		conf := load()
		assert.NotEmpty(t, conf.String())
	})
}

func TestConfigure(t *testing.T) {
	t.Run("with leveldb type, returns instance specify", func(t *testing.T) {
		dbPath := path.Join(os.TempDir(), "leveldb")
		os.RemoveAll(dbPath)

		_, err := Configure("-t", "leveldb", "-s", dbPath)
		if assert.NoError(t, err) {
			assert.IsType(t, new(leveldb_storage.LeveldbStorage), actualRepository())
		}
	})

	t.Run("when fail to load leveldb, returns error", func(t *testing.T) {
		dbPath := path.Join(os.TempDir(), "readonly")
		os.RemoveAll(dbPath)
		if err := ioutil.WriteFile(dbPath, nil, 0200); err != nil {
			t.Fatal(err)
		}

		_, err := Configure("-t", "leveldb", "-s", path.Join(os.TempDir(), "readonly"))
		assert.Error(t, err)
	})

	t.Run("with file type, returns instance specify", func(t *testing.T) {
		_, err := Configure("-t", "file")
		if assert.NoError(t, err) {
			assert.IsType(t, new(file_storage.FileStorage), actualRepository())
		}
	})

	t.Run("with boltdb type, returns instance specify", func(t *testing.T) {
		dbPath := path.Join(os.TempDir(), "boltdb")
		os.RemoveAll(dbPath)

		_, err := Configure("-t", "boltdb", "-s", dbPath)
		if assert.NoError(t, err) {
			assert.IsType(t, new(boltdb_storage.BoltdbStorage), actualRepository())
		}
	})

	t.Run("when fail to load boltdb, returns error", func(t *testing.T) {
		dbPath := path.Join(os.TempDir(), "err_boltdb")
		os.RemoveAll(dbPath)
		os.MkdirAll(dbPath, 0600)

		_, err := Configure("-t", "boltdb", "-s", dbPath)
		assert.Error(t, err)
	})

	t.Run("with unknown type, returns error", func(t *testing.T) {
		_, err := Configure("-t", "unknown")
		assert.Error(t, err)
	})
}

func actualRepository() interface{} {
	photoController := controller.New()
	container.Get(photoController)

	pcv := reflect.Indirect(reflect.ValueOf(photoController))
	ps := pcv.Field(0).Interface()
	psv := reflect.ValueOf(ps)
	rv := reflect.Indirect(psv).Field(0)
	return rv.Interface()
}
