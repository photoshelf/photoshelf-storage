package boltdb_storage

import (
	"github.com/boltdb/bolt"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
)

type BoltdbStorage struct {
	db *bolt.DB
}

func New(path string) (*BoltdbStorage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("photos"))
		return err
	}); err != nil {
		return nil, err
	}

	return &BoltdbStorage{db}, nil
}

func (storage *BoltdbStorage) Save(photograph photo.Photo) (*photo.Identifier, error) {
	data := photograph.Image()
	id := photograph.ID()
	if photograph.IsNew() {
		id = *photo.NewIdentifier(data)
	}

	if err := storage.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("photos")).Put([]byte(id.Value()), data)
	}); err != nil {
		return nil, err
	}

	return &id, nil
}

func (storage *BoltdbStorage) Read(id photo.Identifier) (*photo.Photo, error) {
	var photograph *photo.Photo
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("photos")).Get([]byte(id.Value()))
		if data == nil {
			return photo.ErrNotFound
		}
		photograph = photo.Of(id, data)
		return nil
	}); err != nil {
		return nil, &photo.ResourceError{ID: id, Err: err}
	}

	return photograph, nil
}

func (storage *BoltdbStorage) Delete(id photo.Identifier) error {
	return storage.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("photos")).Delete([]byte(id.Value()))
	})
}
