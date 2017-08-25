package service

import (
	"github.com/photoshelf/photoshelf-storage/model"
)

type PhotoService struct {
	Repository model.Repository `inject:""`
}

func (service *PhotoService) Save(photo model.Photo) (*model.Identifier, error) {
	return service.Repository.Save(photo)
}

func (service *PhotoService) Find(id model.Identifier) (*model.Photo, error) {
	return service.Repository.Read(id)
}

func (service *PhotoService) Delete(id model.Identifier) error {
	return service.Repository.Delete(id)
}
