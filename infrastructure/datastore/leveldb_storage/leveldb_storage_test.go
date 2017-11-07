package leveldb_storage

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
	t.Run("with wrong directory (file)", func(t *testing.T) {
		dbPath := path.Join(os.TempDir(), "readonly")
		file, err := os.Create(dbPath)
		assert.NoError(t, err)
		file.Close()

		instance, err := NewLeveldbStorage(dbPath)
		if assert.Error(t, err) {
			assert.Nil(t, instance)
		}
	})

	t.Run("with correct directory", func(t *testing.T) {
		instance, err := NewLeveldbStorage(path.Join(os.TempDir(), "leveldb"))
		if assert.NoError(t, err) {
			assert.NotNil(t, instance)
		}
	})
}

func TestLeveldbStorage_Save(t *testing.T) {
	t.Run("without identifier", func(t *testing.T) {
		var instance *LeveldbStorage
		var photo *model.Photo

		for _, testcase := range []struct {
			name     string
			function func(t *testing.T)
		}{
			{
				"when save photo, generate new identifier",
				func(t *testing.T) {

					identifier, err := instance.Save(*photo)
					if assert.NoError(t, err) {
						assert.NotNil(t, identifier)
					}
				},
			},
		} {
			instance = createInstance()
			photo = model.NewPhoto(readTestData())

			t.Run(testcase.name, testcase.function)
		}
	})

	t.Run("with identifier", func(t *testing.T) {
		var instance *LeveldbStorage
		var photo *model.Photo

		for _, testcase := range []struct {
			name     string
			function func(t *testing.T)
		}{
			{
				"when save photo",
				func(t *testing.T) {
					identifier, err := instance.Save(*photo)
					assert.NoError(t, err)

					t.Run("returns identifier has same value", func(t *testing.T) {
						actual := photo.Id()
						assert.EqualValues(t, actual.Value(), identifier.Value())
					})

					t.Run("stored same binary", func(t *testing.T) {
						actual, err := instance.db.Get([]byte("testdata"), nil)
						if err != nil {
							assert.Fail(t, "fail load data.")
						}
						assert.EqualValues(t, readTestData(), actual)
					})
				},
			},
			{
				"when db closed, returns error",
				func(t *testing.T) {
					instance.db.Close()

					_, err := instance.Save(*photo)
					assert.Error(t, err)
				},
			},
		} {
			instance = createInstance()
			photo = model.PhotoOf(*model.IdentifierOf("testdata"), readTestData())

			t.Run(testcase.name, testcase.function)
		}
	})
}

func TestLeveldbStorage_Read(t *testing.T) {
	var instance *LeveldbStorage

	for _, testcase := range []struct {
		name     string
		function func(t *testing.T)
	}{
		{
			"when try to read no key, returns err",
			func(t *testing.T) {
				_, err := instance.Read(*model.IdentifierOf("noKey"))
				assert.Error(t, err)
			},
		}, {
			"same value between src and stored value",
			func(t *testing.T) {
				photo, err := instance.Read(*model.IdentifierOf("testdata"))
				if assert.NoError(t, err) {
					assert.EqualValues(t, readTestData(), photo.Image())
				}
			},
		},
	} {
		instance = createInstance()
		err := instance.db.Put([]byte("testdata"), readTestData(), nil)
		assert.NoError(t, err)

		t.Run(testcase.name, testcase.function)
	}
}

func TestLeveldbStorage_Delete(t *testing.T) {
	var instance *LeveldbStorage

	for _, testcase := range []struct {
		name     string
		function func(t *testing.T)
	}{
		{
			"when delete existing key, returns no error",
			func(t *testing.T) {
				err := instance.Delete(*model.IdentifierOf("testdata"))
				if assert.NoError(t, err) {
					actual, _ := instance.db.Get([]byte("testdata"), nil)
					assert.EqualValues(t, []byte{}, actual)
				}
			},
		},
	} {
		instance = createInstance()
		err := instance.db.Put([]byte("testdata"), readTestData(), nil)
		assert.NoError(t, err)

		t.Run(testcase.name, testcase.function)
	}
}

func BenchmarkWithEmptyData(b *testing.B) {
	var instance *LeveldbStorage
	data := readTestData()

	for _, testcase := range []struct {
		name     string
		function func(b *testing.B)
	}{
		{
			"write override",
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					photo := model.PhotoOf(*model.IdentifierOf("testdata"), data)
					instance.Save(*photo)
				}
			},
		}, {
			"write new",
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					photo := model.PhotoOf(*model.IdentifierOf(fmt.Sprintf("testdata-%d", i)), data)
					instance.Save(*photo)
				}
			},
		},
	} {
		instance = createInstance()

		b.ResetTimer()
		b.Run(testcase.name, testcase.function)
	}
}

func BenchmarkWithData(b *testing.B) {
	var instance *LeveldbStorage
	data := readTestData()

	for _, testcase := range []struct {
		name     string
		function func(b *testing.B)
	}{
		{
			"read same data",
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					instance.Read(*model.IdentifierOf("testdata-0"))
				}
			},
		}, {
			"read different data",
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					instance.Read(*model.IdentifierOf(fmt.Sprintf("testdata-%d", i)))
				}
			},
		},
	} {
		instance = createInstance()
		for i := 0; i < 100; i++ {
			key := []byte(fmt.Sprintf("testdata-%d", i))
			err := instance.db.Put(key, data, nil)
			if err != nil {
				assert.NoError(b, err, "failure testdata setting.")
			}
		}

		b.ResetTimer()
		b.Run(testcase.name, testcase.function)
	}
}

func readTestData() []byte {
	testdataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	body, err := os.Open(path.Join(testdataPath, "e3158990bdee63f8594c260cd51a011d"))
	if err != nil {
		panic(err)
	}
	bytea, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return bytea
}

func createInstance() *LeveldbStorage {
	dataPath := path.Join(os.TempDir(), "leveldb")
	if err := os.RemoveAll(dataPath); err != nil {
		panic(err)
	}

	instance, err := NewLeveldbStorage(dataPath)
	if err != nil {
		panic(err)
	}
	return instance
}
