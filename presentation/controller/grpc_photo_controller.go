package controller

import (
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/photoshelf/photoshelf-storage/presentation/protobuf"
	"golang.org/x/net/context"
)

type grpcPhotoControllerImpl struct {
	Service service.PhotoService `inject:""`
}

func NewGrpcPhotoController() protobuf.PhotoServiceServer {
	return &grpcPhotoControllerImpl{}
}

func (ctrl *grpcPhotoControllerImpl) Save(ctx context.Context, req *protobuf.Photo) (*protobuf.Id, error) {
	var model *photo.Photo
	if req.Id != nil {
		model = photo.Of(*photo.IdentifierOf(req.Id.Value), req.Image)
	} else {
		model = photo.New(req.Image)
	}

	id, err := ctrl.Service.Save(*model)
	if err != nil {
		return nil, err
	}

	return &protobuf.Id{Value: id.Value()}, nil
}

func (ctrl *grpcPhotoControllerImpl) Find(ctx context.Context, req *protobuf.Id) (*protobuf.Photo, error) {
	id := photo.IdentifierOf(req.Value)
	photograph, err := ctrl.Service.Find(*id)
	if err != nil {
		return nil, err
	}
	return &protobuf.Photo{Id: &protobuf.Id{Value: photograph.Id().Value()}, Image: photograph.Image()}, nil
}

func (ctrl *grpcPhotoControllerImpl) Delete(ctx context.Context, req *protobuf.Id) (*protobuf.Empty, error) {
	id := photo.IdentifierOf(req.Value)
	if err := ctrl.Service.Delete(*id); err != nil {
		return nil, err
	}
	return &protobuf.Empty{}, nil
}
