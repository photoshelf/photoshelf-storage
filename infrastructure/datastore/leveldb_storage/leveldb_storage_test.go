package leveldb_storage

import (
	"errors"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var storage *LeveldbStorage
var testdata []byte

func TestMain(m *testing.M) {
	testdataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	body, _ := os.Open(path.Join(testdataPath, "e3158990bdee63f8594c260cd51a011d"))
	testdata, _ = ioutil.ReadAll(body)

	dataPath := path.Join(os.TempDir(), "leveldb")
	storage, _ = NewLeveldbStorage(dataPath)

	code := m.Run()

	storage.db.Close()
	storage = nil
	os.Exit(code)
}

func TestWithNoKeys(t *testing.T) {
	storage.db.Delete([]byte("testdata"), nil)

	t.Run("same data between src and dst", func(t *testing.T) {
		photo := model.PhotoOf(*model.IdentifierOf("testdata"), testdata)
		_, err := storage.Save(*photo)

		if assert.NoError(t, err) {
			actual, err := storage.db.Get([]byte("testdata"), nil)
			if err != nil {
				assert.Fail(t, "fail load data.")
			}
			assert.EqualValues(t, testdata, actual)
		}
	})
}

func TestExistData(t *testing.T) {
	err := storage.db.Put([]byte("testdata"), testdata, nil)
	assert.NoError(t, err, "failure testdata setting.")

	t.Run("same data between src and read", func(t *testing.T) {
		photo, err := storage.Read(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			assert.EqualValues(t, testdata, photo.Image())
		}
	})

	t.Run("deleted data", func(t *testing.T) {
		err := storage.Delete(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			actual, err := storage.db.Get([]byte("testdata"), nil)
			assert.EqualValues(t, []byte{}, actual)
			assert.EqualValues(t, errors.New("leveldb: not found"), err)
		}
	})
}
