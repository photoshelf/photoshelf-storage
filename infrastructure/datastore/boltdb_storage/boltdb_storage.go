package boltdb_storage

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/photoshelf/photoshelf-storage/domain/model"
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

func (storage *BoltdbStorage) Save(photo model.Photo) (*model.Identifier, error) {
	data := photo.Image()
	id := photo.Id()
	if photo.IsNew() {
		id = *model.NewIdentifier(data)
	}

	if err := storage.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("photos")).Put([]byte(id.Value()), data)
	}); err != nil {
		return nil, err
	}

	return &id, nil
}

func (storage *BoltdbStorage) Read(id model.Identifier) (*model.Photo, error) {
	var photo *model.Photo
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("photos")).Get([]byte(id.Value()))
		if data == nil {
			return errors.New(fmt.Sprintf("no such id : %s", id.Value()))
		}
		photo = model.PhotoOf(id, data)
		return nil
	}); err != nil {
		return nil, err
	}

	return photo, nil
}

func (storage *BoltdbStorage) Delete(id model.Identifier) error {
	return storage.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("photos")).Delete([]byte(id.Value()))
	})
}
