package service

import (
	"github.com/photoshelf/photoshelf-storage/domain/model"
)

type PhotoService interface {
	Save(photo model.Photo) (*model.Identifier, error)
	Find(id model.Identifier) (*model.Photo, error)
	Delete(id model.Identifier) error
}
