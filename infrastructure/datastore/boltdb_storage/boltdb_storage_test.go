package boltdb_storage

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo/phototest"
)

func TestNew(t *testing.T) {
	t.Run("with wrong directory (file)", func(t *testing.T) {
		dbPath := path.Join(os.TempDir(), "readonly")
		os.RemoveAll(dbPath)
		if err := ioutil.WriteFile(dbPath, []byte("not empty"), 0200); err != nil {
			t.Fatal(err)
		}

		instance, err := New(dbPath)
		if assert.Error(t, err) {
			assert.Nil(t, instance)
		}
	})

	t.Run("with correct directory", func(t *testing.T) {
		instance, err := New(path.Join(os.TempDir(), "boltdb"))
		if assert.NoError(t, err) {
			assert.NotNil(t, instance)
		}
	})
}

func TestBoltdbStorage_Save(t *testing.T) {
	t.Run("save without identifier, generate new identifier", func(t *testing.T) {
		instance := createInstance(t)
		photograph := photo.New(readTestData(t))

		identifier, err := instance.Save(*photograph)
		if assert.NoError(t, err) {
			assert.NotNil(t, identifier)
		}

		instance.db.Close()
	})

	t.Run("save with identifier", func(t *testing.T) {
		instance := createInstance(t)
		photograph := *photo.Of(*photo.IdentifierOf("testdata"), readTestData(t))

		identifier, err := instance.Save(photograph)

		if assert.NoError(t, err) {
			t.Run("returns identifier has same value", func(t *testing.T) {
				actual := photograph.Id()
				assert.EqualValues(t, actual.Value(), identifier.Value())
			})

			t.Run("stored same binary", func(t *testing.T) {
				var actual []byte
				instance.db.View(func(tx *bolt.Tx) error {
					actual = tx.Bucket([]byte("photos")).Get([]byte("testdata"))
					return nil
				})
				assert.EqualValues(t, readTestData(t), actual)
			})
		}
		instance.db.Close()
	})
}

func TestBoltdbStorage_Read(t *testing.T) {
	instance := createInstance(t)
	err := instance.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("photos")).Put([]byte("testdata"), readTestData(t))
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("with no key, returns err", func(t *testing.T) {
		_, err := instance.Read(*photo.IdentifierOf("noKey"))
		assert.Error(t, err)
	})

	t.Run("returns same data with source", func(t *testing.T) {
		photograph, err := instance.Read(*photo.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			assert.EqualValues(t, readTestData(t), photograph.Image())
		}
	})

	instance.db.Close()
}

func TestBoltdbStorage_Delete(t *testing.T) {
	instance := createInstance(t)
	err := instance.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("photos")).Put([]byte("testdata"), readTestData(t))
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("when delete existing key, returns no error", func(t *testing.T) {
		err := instance.Delete(*photo.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			instance.db.View(func(tx *bolt.Tx) error {
				numOfKeys := tx.Bucket([]byte("photos")).Stats().KeyN
				assert.EqualValues(t, 0, numOfKeys)
				return nil
			})
		}
	})
}

func BenchmarkBoltdbStorage_Save(b *testing.B) {
	data := readTestData(b)

	b.Run("override", func(b *testing.B) {
		instance := createInstance(b)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", 0)
			photograph := *photo.Of(*photo.IdentifierOf(key), data)
			instance.Save(photograph)
		}
		b.StopTimer()

		instance.db.Close()
	})

	b.Run("with new key", func(b *testing.B) {
		instance := createInstance(b)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", i)
			photograph := *photo.Of(*photo.IdentifierOf(key), data)
			instance.Save(photograph)
		}
		b.StopTimer()

		instance.db.Close()
	})

	b.Run("random data", func(b *testing.B) {
		instance := createInstance(b)
		randomTestData := phototest.RandomTestData(b)

		b.ResetTimer()
		for i := 1; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", i%20)
			photograph := *photo.Of(*photo.IdentifierOf(key), randomTestData[i%20])
			instance.Save(photograph)
		}
		b.StopTimer()

		instance.db.Close()
	})
}

func BenchmarkBoltdbStorage_Read(b *testing.B) {
	dataSet := phototest.RandomTestData(b)
	instance := createInstance(b)

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("testdata-%d", i%100))
		err := instance.db.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("photos")).Put(key, dataSet[i%20])
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	b.Run("same key", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", 0)
			instance.Read(*photo.IdentifierOf(key))
		}
		b.StopTimer()

		instance.db.Close()
	})

	b.Run("sequential key", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", i)
			instance.Read(*photo.IdentifierOf(key))
		}
		b.StopTimer()

		instance.db.Close()
	})

	b.Run("random key", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("testdata-%d", rand.Intn(100))
			instance.Read(*photo.IdentifierOf(key))
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

func createInstance(tb testing.TB) *BoltdbStorage {
	tb.Helper()

	dataPath := path.Join(os.TempDir(), "boltdb")
	if err := os.RemoveAll(dataPath); err != nil {
		tb.Fatal(err)
	}

	instance, err := New(dataPath)
	if err != nil {
		tb.Fatal(err)
	}
	if err := instance.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("photos"))
		return err
	}); err != nil {
		tb.Fatal(err)
	}
	return instance
}
