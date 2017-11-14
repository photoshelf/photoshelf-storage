package controller

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/photoshelf/photoshelf-storage/application/errors"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/photoshelf/photoshelf-storage/presentation/view"
	"io/ioutil"
	"net/http"
)

type PhotoController interface {
	Get(c echo.Context) error
	Post(c echo.Context) error
	Put(c echo.Context) error
	Delete(c echo.Context) error
}

type photoControllerImpl struct {
	Service service.PhotoService `inject:""`
}

func New() PhotoController {
	return &photoControllerImpl{}
}

func (controller *photoControllerImpl) Get(c echo.Context) error {
	id := photo.IdentifierOf(c.Param("id"))
	photograph, err := controller.Service.Find(*id)
	if err != nil {
		log.Error(err)
		switch e := err.(type) {
		case errors.NotFoundError:
			return c.JSON(http.StatusNotFound, e.Error())
		default:
			return e
		}
	}

	mimeType := http.DetectContentType(photograph.Image())
	return c.Blob(http.StatusOK, mimeType, photograph.Image())
}

func (controller *photoControllerImpl) Post(c echo.Context) error {
	data, err := readPhotoBytes(c)
	if err != nil {
		log.Error(err)
		return err
	}

	photograph := photo.New(data)
	id, err := controller.Service.Save(*photograph)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusCreated, view.Created{ID: id.Value()})
}

func (controller *photoControllerImpl) Put(c echo.Context) error {
	data, err := readPhotoBytes(c)
	if err != nil {
		log.Error(err)
		return err
	}

	id := photo.IdentifierOf(c.Param("id"))
	if _, err := controller.Service.Save(*photo.Of(*id, data)); err != nil {
		log.Error(err)
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (controller *photoControllerImpl) Delete(c echo.Context) error {
	id := photo.IdentifierOf(c.Param("id"))

	if err := controller.Service.Delete(*id); err != nil {
		log.Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func readPhotoBytes(c echo.Context) ([]byte, error) {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return nil, err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}
