package service

import (
	"github.com/duck8823/photoshelf-storage/model"
)

type PhotoService struct {
	repository model.Repository
}

func NewPhotoService(repository model.Repository) *PhotoService {
	return &PhotoService{repository}
}

func (service *PhotoService) Save(photo model.Photo) (*model.Identifier, error) {
	return service.repository.Save(photo)
}

func (service *PhotoService) Find(id model.Identifier) (*model.Photo, error) {
	return service.repository.Read(id)
}

func (service *PhotoService) Delete(id model.Identifier) error {
	return service.repository.Delete(id)
}
