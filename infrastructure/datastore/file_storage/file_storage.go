package file_storage

import (
	"github.com/photoshelf/photoshelf-storage/application/errors"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"io/ioutil"
	"os"
	"path"
)

type FileStorage struct {
	baseDir string
}

func New(baseDir string) *FileStorage {
	return &FileStorage{baseDir}
}

func (storage *FileStorage) Save(photograph photo.Photo) (*photo.Identifier, error) {
	data := photograph.Image()
	id := photograph.ID()
	if photograph.IsNew() {
		id = *photo.NewIdentifier(data)
	}

	filename := path.Join(storage.baseDir, id.Value())
	ioutil.WriteFile(filename, data, 0600)

	return &id, nil
}

func (storage *FileStorage) Read(id photo.Identifier) (*photo.Photo, error) {
	filename := path.Join(storage.baseDir, id.Value())
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.NotFound(id.Value())
	}

	return photo.Of(id, data), nil
}

func (storage *FileStorage) Delete(id photo.Identifier) error {
	return os.Remove(path.Join(storage.baseDir, id.Value()))
}
