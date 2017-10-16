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

	dst, err := os.Create(path.Join(storage.baseDir, id.Value()))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := dst.Write(data); err != nil {
		return nil, err
	}
	return &id, nil
}

func (storage *FileStorage) Read(id model.Identifier) (*model.Photo, error) {
	file, err := os.Open(path.Join(storage.baseDir, id.Value()))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
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
