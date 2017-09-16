package datastore

import (
	"testing"
	"os"
	"path"
	"io/ioutil"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/stretchr/testify/assert"
	"github.com/boltdb/bolt"
)

var boltdb *BoltdbStorage
var testdata []byte

func TestMain(m *testing.M) {
	testdataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	body, _ := os.Open(path.Join(testdataPath, "e3158990bdee63f8594c260cd51a011d"))
	testdata, _ = ioutil.ReadAll(body)

	dataPath := path.Join(os.TempDir(), "boltdb")
	boltdb, _ = NewBoltdbStorage(dataPath)

	code := m.Run()

	boltdb.db.Close()
	boltdb = nil
	os.Exit(code)
}

func TestEmptyBucket(t *testing.T) {
	boltdb.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("photos"))
		return nil
	})

	t.Run("same data between src and dst", func(t *testing.T) {
		photo := model.PhotoOf(*model.IdentifierOf("testdata"), testdata)
		_, err := boltdb.Save(*photo)

		if assert.NoError(t, err) {
			boltdb.db.View(func(tx *bolt.Tx) error {
				photos := tx.Bucket([]byte("photos"))
				actual := photos.Get([]byte("testdata"))

				assert.EqualValues(t, testdata, actual)
				assert.EqualValues(t, 1, photos.Stats().KeyN)

				return nil
			})
		}
	})
}

func TestExistData(t *testing.T) {
	err := boltdb.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("photos"))
		photos, err := tx.CreateBucketIfNotExists([]byte("photos"))
		if err != nil {
			return err
		}
		return photos.Put([]byte("testdata"), testdata)
	})
	assert.NoError(t, err, "failure testdata setting.")

	t.Run("same data between src and read", func(t *testing.T) {
		photo, err := boltdb.Read(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			assert.EqualValues(t, testdata, photo.Image())
		}
	})

	t.Run("deleted data", func(t *testing.T) {
		err := boltdb.Delete(*model.IdentifierOf("testdata"))
		if assert.NoError(t, err) {
			boltdb.db.View(func(tx *bolt.Tx) error {
				photos := tx.Bucket([]byte("photos"))
				assert.EqualValues(t, 0, photos.Stats().KeyN)
				return nil
			})
		}
	})
}

