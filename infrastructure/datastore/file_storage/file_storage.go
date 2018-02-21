package file_storage

import (
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type FileStorage struct {
	baseDir string
}

func New(baseDir string) *FileStorage {
	return &FileStorage{baseDir}
}

func (storage *FileStorage) Save(photograph photo.Photo) (*photo.Identifier, error) {
	data := photograph.Image()
	id := photograph.Id()
	if photograph.IsNew() {
		id = photo.NewIdentifier(data)
	}

	filename := path.Join(storage.baseDir, id.Value())
	ioutil.WriteFile(filename, data, 0600)

	return id, nil
}

func (storage *FileStorage) Read(id photo.Identifier) (*photo.Photo, error) {
	filename := path.Join(storage.baseDir, id.Value())
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		pathErr := err.(*os.PathError)
		errno := pathErr.Err.(syscall.Errno)
		if errno == syscall.ENOENT {
			return nil, &photo.ResourceError{ID: id, Err: photo.ErrNotFound}
		}
		return nil, &photo.ResourceError{ID: id, Err: err}
	}

	return photo.Of(id, data), nil
}

func (storage *FileStorage) Delete(id photo.Identifier) error {
	return os.Remove(path.Join(storage.baseDir, id.Value()))
}
