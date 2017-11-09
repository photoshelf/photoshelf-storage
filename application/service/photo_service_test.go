package service

import (
	"errors"
	"github.com/facebookgo/inject"
	"github.com/golang/mock/gomock"
	"github.com/photoshelf/photoshelf-storage/domain/model/mock_photo"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPhotoServiceImpl_Find(t *testing.T) {
	t.Run("when repository returns object, it returns object", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		photograph := photo.Of(*photo.IdentifierOf("id"), []byte("test"))
		mock_repository := mock_photo.NewMockRepository(ctrl)
		mock_repository.EXPECT().
			Read(gomock.Any()).
			Return(photograph, nil)

		photo_service := New()
		if err := inject.Populate(photo_service, mock_repository); err != nil {
			t.Fatal(err)
		}

		actual, err := photo_service.Find(*photo.IdentifierOf("any"))
		if assert.NoError(t, err) {
			assert.EqualValues(t, photograph, actual)
		}
	})

	t.Run("when repository returns error, it returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock_repository := mock_photo.NewMockRepository(ctrl)
		mock_repository.EXPECT().
			Read(gomock.Any()).
			Return(nil, errors.New("expected error"))

		photo_service := New()
		if err := inject.Populate(photo_service, mock_repository); err != nil {
			t.Fatal(err)
		}

		actual, err := photo_service.Find(*photo.IdentifierOf("any"))
		if assert.Error(t, err) {
			assert.Nil(t, actual)
		}
	})
}

func TestPhotoServiceImpl_Save(t *testing.T) {
	t.Run("when repository returns object, it returns object", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := photo.IdentifierOf("id")
		mock_repository := mock_photo.NewMockRepository(ctrl)
		mock_repository.EXPECT().
			Save(gomock.Any()).
			Return(id, nil)

		photo_service := New()
		if err := inject.Populate(photo_service, mock_repository); err != nil {
			t.Fatal(err)
		}

		actual, err := photo_service.Save(*photo.Of(*id, nil))
		if assert.NoError(t, err) {
			assert.EqualValues(t, id, actual)
		}
	})

	t.Run("when repository returns error, it returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock_repository := mock_photo.NewMockRepository(ctrl)
		mock_repository.EXPECT().
			Save(gomock.Any()).
			Return(nil, errors.New("expected error"))

		photo_service := New()
		if err := inject.Populate(photo_service, mock_repository); err != nil {
			t.Fatal(err)
		}

		actual, err := photo_service.Save(*photo.Of(*photo.IdentifierOf("any"), nil))
		if assert.Error(t, err) {
			assert.Nil(t, actual)
		}
	})
}

func TestPhotoServiceImpl_Delete(t *testing.T) {
	t.Run("when repository returns error, it returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock_repository := mock_photo.NewMockRepository(ctrl)
		mock_repository.EXPECT().
			Delete(gomock.Any()).
			Return(errors.New("expected error"))

		photo_service := New()
		if err := inject.Populate(photo_service, mock_repository); err != nil {
			t.Fatal(err)
		}

		assert.Error(t, photo_service.Delete(*photo.IdentifierOf("any")))
	})
}
