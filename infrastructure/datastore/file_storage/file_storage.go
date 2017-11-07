package file_storage

import (
	"github.com/photoshelf/photoshelf-storage/domain/model"
	"io/ioutil"
	"os"
	"path"
)

type FileStorage struct {
	baseDir string
}

func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{baseDir}
}

func (storage *FileStorage) Save(photo model.Photo) (*model.Identifier, error) {
	data := photo.Image()
	id := photo.Id()
	if photo.IsNew() {
		id = *model.NewIdentifier(data)
	}

	filename := path.Join(storage.baseDir, id.Value())
	ioutil.WriteFile(filename, data, 0600)

	return &id, nil
}

func (storage *FileStorage) Read(id model.Identifier) (*model.Photo, error) {
	filename := path.Join(storage.baseDir, id.Value())
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return model.PhotoOf(id, data), nil
}

func (storage *FileStorage) Delete(id model.Identifier) error {
	if err := os.Remove(path.Join(storage.baseDir, id.Value())); err != nil {
		return err
	}
	return nil
}
