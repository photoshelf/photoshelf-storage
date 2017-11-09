package service

import (
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
)

type PhotoService interface {
	Save(photo photo.Photo) (*photo.Identifier, error)
	Find(id photo.Identifier) (*photo.Photo, error)
	Delete(id photo.Identifier) error
}

type photoServiceImpl struct {
	Repository photo.Repository `inject:""`
}

func New() PhotoService {
	return &photoServiceImpl{}
}

func (service *photoServiceImpl) Save(photo photo.Photo) (*photo.Identifier, error) {
	return service.Repository.Save(photo)
}

func (service *photoServiceImpl) Find(id photo.Identifier) (*photo.Photo, error) {
	return service.Repository.Read(id)
}

func (service *photoServiceImpl) Delete(id photo.Identifier) error {
	return service.Repository.Delete(id)
}
