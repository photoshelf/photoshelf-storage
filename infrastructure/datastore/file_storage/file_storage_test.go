package file_storage

import (
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var filestorage *FileStorage
var testdata []byte

func TestMain(m *testing.M) {
	testdataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	body, _ := os.Open(path.Join(testdataPath, "e3158990bdee63f8594c260cd51a011d"))
	testdata, _ = ioutil.ReadAll(body)

	dataPath := path.Join(os.TempDir(), "file")
	os.RemoveAll(dataPath)
	os.MkdirAll(dataPath, 0700)

	filestorage = NewFileStorage(dataPath)

	os.Exit(m.Run())
}

func TestEmptyDirectory(t *testing.T) {
	os.RemoveAll(filestorage.baseDir)
	os.MkdirAll(filestorage.baseDir, 0700)

	t.Run("same data between src and dst", func(t *testing.T) {
		photo := model.PhotoOf(*model.IdentifierOf("testdata"), testdata)
		_, err := filestorage.Save(*photo)

		if assert.NoError(t, err) {
			file, _ := os.Open(path.Join(filestorage.baseDir, "testdata"))
			actual, _ := ioutil.ReadAll(file)

			assert.EqualValues(t, testdata, actual)

			files, _ := ioutil.ReadDir(filestorage.baseDir)
			assert.EqualValues(t, 1, len(files))
		}
	})
}

func TestExistData(t *testing.T) {
	os.RemoveAll(filestorage.baseDir)
	os.MkdirAll(filestorage.baseDir, 0700)
	err := ioutil.WriteFile(path.Join(filestorage.baseDir, "testdata"), testdata, 0700)
	assert.NoError(t, err, "failure testdata setting.")

	t.Run("same data between src and read", func(t *testing.T) {
		file, _ := os.Open(path.Join(filestorage.baseDir, "testdata"))
		actual, _ := ioutil.ReadAll(file)

		assert.EqualValues(t, testdata, actual)
	})

	t.Run("deleted data", func(t *testing.T) {
		err := filestorage.Delete(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			files, _ := ioutil.ReadDir(filestorage.baseDir)
			assert.EqualValues(t, 0, len(files))
		}
	})
}
