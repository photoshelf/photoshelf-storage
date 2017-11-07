package file_storage

import (
	"fmt"
	"github.com/photoshelf/photoshelf-storage/domain/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("with correct directory", func(t *testing.T) {
		instance := NewFileStorage(path.Join(os.TempDir(), "file_storage"))
		assert.NotNil(t, instance)
	})
}

func TestFileStorage_Save(t *testing.T) {
	t.Run("save without identifier, generate new identifier", func(t *testing.T) {
		instance := createInstance(t)
		photo := model.NewPhoto(readTestData(t))

		identifier, err := instance.Save(*photo)
		if assert.NoError(t, err) {
			assert.NotNil(t, identifier)
		}
	})

	t.Run("save with identifier", func(t *testing.T) {
		instance := createInstance(t)
		photo := model.PhotoOf(*model.IdentifierOf("testdata"), readTestData(t))

		identifier, err := instance.Save(*photo)

		if assert.NoError(t, err) {
			t.Run("returns identifier has same value", func(t *testing.T) {
				actual := photo.Id()
				assert.EqualValues(t, actual.Value(), identifier.Value())
			})

			t.Run("stored same binary", func(t *testing.T) {
				file, err := os.Open(path.Join(instance.baseDir, "testdata"))
				if err != nil {
					t.Fatal(err)
				}
				actual, err := ioutil.ReadAll(file)
				if err != nil {
					t.Fatal(err)
				}

				assert.EqualValues(t, readTestData(t), actual)
			})
		}
	})
}

func TestFileStorage_Read(t *testing.T) {
	instance := createInstance(t)
	if err := ioutil.WriteFile(path.Join(instance.baseDir, "testdata"), readTestData(t), 0700); err != nil {
		t.Fatal(err)
	}

	t.Run("with no key, returns err", func(t *testing.T) {
		_, err := instance.Read(*model.IdentifierOf("noKey"))
		assert.Error(t, err)
	})

	t.Run("returns same data with source", func(t *testing.T) {
		photo, err := instance.Read(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			assert.EqualValues(t, readTestData(t), photo.Image())
		}
	})
}

func TestFileStorage_Delete(t *testing.T) {
	instance := createInstance(t)
	if err := ioutil.WriteFile(path.Join(instance.baseDir, "testdata"), readTestData(t), 0700); err != nil {
		t.Fatal(err)
	}

	t.Run("when delete existing key, returns no error", func(t *testing.T) {
		err := instance.Delete(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			files, _ := ioutil.ReadDir(instance.baseDir)
			assert.EqualValues(t, 0, len(files))
		}
	})
}

func BenchmarkFileStorage_Save(b *testing.B) {
	data := readTestData(b)

	b.Run("override", func(b *testing.B) {
		instance := createInstance(b)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", 0)
			photo := *model.PhotoOf(*model.IdentifierOf(key), data)
			instance.Save(photo)
		}
		b.StopTimer()
	})

	b.Run("with new key", func(b *testing.B) {
		instance := createInstance(b)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", i)
			photo := *model.PhotoOf(*model.IdentifierOf(key), data)
			instance.Save(photo)
		}
		b.StopTimer()
	})
}

func BenchmarkFileStorage_Read(b *testing.B) {
	data := readTestData(b)
	instance := createInstance(b)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("testdata-%d", i)
		if err := ioutil.WriteFile(path.Join(instance.baseDir, key), data, 0700); err != nil {
			b.Fatal(err)
		}
	}

	b.Run("same data", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", 0)
			instance.Read(*model.IdentifierOf(key))
		}
		b.StopTimer()
	})

	b.Run("sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", i)
			instance.Read(*model.IdentifierOf(key))
		}
		b.StopTimer()
	})
}

func readTestData(tb testing.TB) []byte {
	tb.Helper()

	testdataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	body, err := os.Open(path.Join(testdataPath, "e3158990bdee63f8594c260cd51a011d"))
	if err != nil {
		tb.Fatal(err)
	}
	bytea, err := ioutil.ReadAll(body)
	if err != nil {
		tb.Fatal(err)
	}
	return bytea
}

func createInstance(tb testing.TB) *FileStorage {
	tb.Helper()

	dataPath := path.Join(os.TempDir(), "file_storage")
	if err := os.RemoveAll(dataPath); err != nil {
		tb.Fatal(err)
	}
	if err := os.MkdirAll(dataPath, 0700); err != nil {
		tb.Fatal(err)
	}
	return NewFileStorage(dataPath)
}
