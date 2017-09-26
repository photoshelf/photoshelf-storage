package file_storage

import (
	"fmt"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var storage *FileStorage
var testdata []byte

func TestMain(m *testing.M) {
	testdataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	body, _ := os.Open(path.Join(testdataPath, "e3158990bdee63f8594c260cd51a011d"))
	testdata, _ = ioutil.ReadAll(body)

	dataPath := path.Join(os.TempDir(), "file")
	os.RemoveAll(dataPath)
	os.MkdirAll(dataPath, 0700)

	storage = NewFileStorage(dataPath)

	os.Exit(m.Run())
}

func TestEmptyDirectory(t *testing.T) {
	os.RemoveAll(storage.baseDir)
	os.MkdirAll(storage.baseDir, 0700)

	t.Run("same data between src and dst", func(t *testing.T) {
		photo := model.PhotoOf(*model.IdentifierOf("testdata"), testdata)
		_, err := storage.Save(*photo)

		if assert.NoError(t, err) {
			file, _ := os.Open(path.Join(storage.baseDir, "testdata"))
			actual, _ := ioutil.ReadAll(file)

			assert.EqualValues(t, testdata, actual)

			files, _ := ioutil.ReadDir(storage.baseDir)
			assert.EqualValues(t, 1, len(files))
		}
	})
}

func TestExistData(t *testing.T) {
	os.RemoveAll(storage.baseDir)
	os.MkdirAll(storage.baseDir, 0700)
	err := ioutil.WriteFile(path.Join(storage.baseDir, "testdata"), testdata, 0700)
	assert.NoError(t, err, "failure testdata setting.")

	t.Run("same data between src and read", func(t *testing.T) {
		file, _ := os.Open(path.Join(storage.baseDir, "testdata"))
		actual, _ := ioutil.ReadAll(file)

		assert.EqualValues(t, testdata, actual)
	})

	t.Run("deleted data", func(t *testing.T) {
		err := storage.Delete(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			files, _ := ioutil.ReadDir(storage.baseDir)
			assert.EqualValues(t, 0, len(files))
		}
	})
}

func BenchmarkWithEmptyData(b *testing.B) {
	os.RemoveAll(storage.baseDir)
	os.MkdirAll(storage.baseDir, 0700)
	err := ioutil.WriteFile(path.Join(storage.baseDir, "testdata"), testdata, 0700)
	assert.NoError(b, err, "failure testdata setting.")

	b.Run("write override", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			photo := model.PhotoOf(*model.IdentifierOf("testdata"), testdata)
			storage.Save(*photo)
		}
	})

	b.Run("write new", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			photo := model.PhotoOf(*model.IdentifierOf(fmt.Sprintf("testdata-%d", i)), testdata)
			storage.Save(*photo)
		}
	})
}

func BenchmarkWithData(b *testing.B) {
	os.RemoveAll(storage.baseDir)
	os.MkdirAll(storage.baseDir, 0700)
	err := ioutil.WriteFile(path.Join(storage.baseDir, "testdata"), testdata, 0700)
	assert.NoError(b, err, "failure testdata setting.")

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("testdata-%d", i)
		if err := ioutil.WriteFile(path.Join(storage.baseDir, key), testdata, 0700); err != nil {
			assert.NoError(b, err, "failure testdata setting.")
		}
	}

	b.Run("read same data", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			storage.Read(*model.IdentifierOf("testdata"))
		}
	})

	b.Run("read different data", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			storage.Read(*model.IdentifierOf(fmt.Sprintf("testdata-%d", i)))
		}
	})
}
