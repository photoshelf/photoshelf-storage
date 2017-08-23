package datastore

import (
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/syndtr/goleveldb/leveldb"
)

type LeveldbStorage struct {
	db *leveldb.DB
}

func NewLeveldbStorage(path string) (*LeveldbStorage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LeveldbStorage{db}, nil
}

func (storage *LeveldbStorage) Save(photo model.Photo) (*model.Identifier, error) {
	data := photo.Image()
	id := photo.Id()
	if photo.IsNew() {
		id = *model.NewIdentifier(data)
	}

	if err := storage.db.Put([]byte(id.Value()), data, nil); err != nil {
		return nil, err
	}

	return &id, nil
}

func (storage *LeveldbStorage) Read(id model.Identifier) (*model.Photo, error) {
	data, err := storage.db.Get([]byte(id.Value()), nil)
	if err != nil {
		return nil, err
	}

	return model.PhotoOf(id, data), nil
}

func (storage *LeveldbStorage) Delete(id model.Identifier) error {
	return storage.db.Delete([]byte(id.Value()), nil)
}
