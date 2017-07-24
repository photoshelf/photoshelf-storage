package controller

import (
	"github.com/duck8823/photoshelf-storage/service"
	"github.com/labstack/echo"
	"github.com/duck8823/photoshelf-storage/model"
	"net/http"
	"github.com/labstack/gommon/log"
	"github.com/duck8823/photoshelf-storage/presentation/view"
	"github.com/duck8823/photoshelf-storage/infrastructure/utility"
)

type PhotoController struct {
	service service.PhotoService
}

func NewPhotoController(service service.PhotoService) *PhotoController {
	return &PhotoController{service}
}

func (controller *PhotoController) Get(c echo.Context) error {
	id := model.IdentifierOf(c.Param("id"))
	photo, err := controller.service.Find(*id)
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
	id, err := controller.service.Save(photo)
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
	if _, err := controller.service.Save(model.PhotoOf(*id, data)); err != nil {
		log.Error(err)
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (controller *PhotoController) Delete(c echo.Context) error {
	id := model.IdentifierOf(c.Param("id"))

	if err := controller.service.Delete(*id); err != nil {
		log.Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
