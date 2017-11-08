package service

import (
	"github.com/photoshelf/photoshelf-storage/domain/model"
)

type PhotoService interface {
	Save(photo model.Photo) (*model.Identifier, error)
	Find(id model.Identifier) (*model.Photo, error)
	Delete(id model.Identifier) error
}

type photoServiceImpl struct {
	Repository model.Repository `inject:""`
}

func New() PhotoService {
	return &photoServiceImpl{}
}

func (service *photoServiceImpl) Save(photo model.Photo) (*model.Identifier, error) {
	return service.Repository.Save(photo)
}

func (service *photoServiceImpl) Find(id model.Identifier) (*model.Photo, error) {
	return service.Repository.Read(id)
}

func (service *photoServiceImpl) Delete(id model.Identifier) error {
	return service.Repository.Delete(id)
}
