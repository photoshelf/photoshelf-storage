package controller

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/domain/model"
	"github.com/photoshelf/photoshelf-storage/infrastructure/utility"
	"github.com/photoshelf/photoshelf-storage/presentation/view"
	"net/http"
)

type PhotoController struct {
	Service service.PhotoService `inject:""`
}

func (controller *PhotoController) Get(c echo.Context) error {
	id := model.IdentifierOf(c.Param("id"))
	photo, err := controller.Service.Find(*id)
	if err != nil {
		log.Error(err)
		return err
	}

	mimeType := http.DetectContentType(photo.Image())
	return c.Blob(http.StatusOK, mimeType, photo.Image())
}

func (controller *PhotoController) Post(c echo.Context) error {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		log.Error(err)
		return err
	}

	data, err := utility.Read(*fileHeader)
	if err != nil {
		log.Error(err)
		return err
	}

	photo := model.NewPhoto(data)
	id, err := controller.Service.Save(*photo)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusCreated, view.Created{Id: id.Value()})
}

func (controller *PhotoController) Put(c echo.Context) error {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		log.Error(err)
		return err
	}

	data, err := utility.Read(*fileHeader)
	if err != nil {
		log.Error(err)
		return err
	}

	id := model.IdentifierOf(c.Param("id"))
	if _, err := controller.Service.Save(*model.PhotoOf(*id, data)); err != nil {
		log.Error(err)
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (controller *PhotoController) Delete(c echo.Context) error {
	id := model.IdentifierOf(c.Param("id"))

	if err := controller.Service.Delete(*id); err != nil {
		log.Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
