package service

import "github.com/photoshelf/photoshelf-storage/domain/model"

type PhotoServiceImpl struct {
	Repository model.Repository `inject:""`
}

func (service *PhotoServiceImpl) Save(photo model.Photo) (*model.Identifier, error) {
	return service.Repository.Save(photo)
}

func (service *PhotoServiceImpl) Find(id model.Identifier) (*model.Photo, error) {
	return service.Repository.Read(id)
}

func (service *PhotoServiceImpl) Delete(id model.Identifier) error {
	return service.Repository.Delete(id)
}
