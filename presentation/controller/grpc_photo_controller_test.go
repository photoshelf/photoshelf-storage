package controller

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/photoshelf/photoshelf-storage/application/mock_service"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/photoshelf/photoshelf-storage/presentation/protobuf"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestGrpcPhotoControllerImpl_Find(t *testing.T) {
	t.Run("when service no error, returns bytes", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Find(*identifier).
			Return(photo.Of(*identifier, readTestData(t)), nil)

		photoController := &grpcPhotoControllerImpl{mockPhotoService}

		actual, err := photoController.Find(context.Background(), &protobuf.Id{Value: identifier.Value()})
		if assert.NoError(t, err) {
			assert.Equal(t, identifier.Value(), actual.Id.Value)
			assert.Equal(t, readTestData(t), actual.Image)
		}
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Find(*photo.IdentifierOf("not_found")).
			Return(nil, errors.New("error not found"))

		photoController := &grpcPhotoControllerImpl{mockPhotoService}

		_, err := photoController.Find(context.Background(), &protobuf.Id{Value: "not_found"})
		assert.Error(t, err)
	})
}

func TestGrpcPhotoController_Save(t *testing.T) {
	t.Run("when service no error, returns identifier", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(gomock.Any()).
			Return(identifier, nil)

		photoController := &grpcPhotoControllerImpl{mockPhotoService}

		actual, err := photoController.Save(context.Background(), &protobuf.Photo{})
		if assert.NoError(t, err) {
			assert.NotNil(t, actual)
		}
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(gomock.Any()).
			Return(nil, errors.New("mock error"))

		photoController := &grpcPhotoControllerImpl{mockPhotoService}

		_, err := photoController.Save(context.Background(), &protobuf.Photo{})
		assert.Error(t, err)
	})
}

func TestGrpcPhotoController_Delete(t *testing.T) {
	t.Run("when service no error, returns status ok", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Delete(*identifier).
			Return(nil)

		photoController := &grpcPhotoControllerImpl{mockPhotoService}

		_, err := photoController.Delete(context.Background(), &protobuf.Id{Value: identifier.Value()})
		assert.NoError(t, err)
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Delete(*identifier).
			Return(errors.New("error"))

		photoController := &grpcPhotoControllerImpl{mockPhotoService}

		_, err := photoController.Delete(context.Background(), &protobuf.Id{Value: identifier.Value()})
		assert.Error(t, err)
	})
}
