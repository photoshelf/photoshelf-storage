package leveldb_storage

import (
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/syndtr/goleveldb/leveldb"
)

type LeveldbStorage struct {
	db *leveldb.DB
}

func New(path string) (*LeveldbStorage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LeveldbStorage{db}, nil
}

func (storage *LeveldbStorage) Save(photograph photo.Photo) (*photo.Identifier, error) {
	data := photograph.Image()
	id := photograph.Id()
	if photograph.IsNew() {
		id = *photo.NewIdentifier(data)
	}

	if err := storage.db.Put([]byte(id.Value()), data, nil); err != nil {
		return nil, err
	}

	return &id, nil
}

func (storage *LeveldbStorage) Read(id photo.Identifier) (*photo.Photo, error) {
	data, err := storage.db.Get([]byte(id.Value()), nil)
	if err != nil {
		return nil, err
	}

	return photo.Of(id, data), nil
}

func (storage *LeveldbStorage) Delete(id photo.Identifier) error {
	return storage.db.Delete([]byte(id.Value()), nil)
}
