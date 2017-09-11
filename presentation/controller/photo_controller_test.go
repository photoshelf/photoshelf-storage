package controller

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/photoshelf/photoshelf-storage/application/service/mock_service"
	"github.com/photoshelf/photoshelf-storage/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	dataPath := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")
	os.Setenv("TEST_DATA_PATH", dataPath)

	code := m.Run()

	os.Unsetenv("TEST_DATA_PATH")
	os.Exit(code)
}

func TestGet(t *testing.T) {
	// Setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	body, _ := os.Open(path.Join(os.Getenv("TEST_DATA_PATH"), identifier.Value()))
	data, _ := ioutil.ReadAll(body)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Find(*identifier).
		Return(model.PhotoOf(*identifier, data), nil)

	photoController := &PhotoController{mockPhotoService}

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(identifier.Value())

	// Assertions
	if assert.NoError(t, photoController.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, data, rec.Body.Bytes())
	}
}

func TestGetWithError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Find(*model.IdentifierOf("not_found")).
		Return(nil, errors.New("file not found"))

	photoController := &PhotoController{mockPhotoService}

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("not_found")

	// Assertions
	assert.Error(t, photoController.Get(c))
}

func TestPost(t *testing.T) {
	// setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	file, _ := os.Open(path.Join(os.Getenv("TEST_DATA_PATH"), identifier.Value()))
	data, _ := ioutil.ReadAll(file)
	file.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Save(gomock.Any()).
		Return(identifier, nil)

	photoController := &PhotoController{mockPhotoService}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("photo", file.Name())
	part.Write(data)
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, photoController.Post(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestPostWithError(t *testing.T) {
	// setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Save(*identifier).
		Times(0)

	photoController := &PhotoController{mockPhotoService}

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	assert.Error(t, photoController.Post(c))
}

func TestPut(t *testing.T) {
	// setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	file, _ := os.Open(path.Join(os.Getenv("TEST_DATA_PATH"), identifier.Value()))
	data, _ := ioutil.ReadAll(file)
	file.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Save(*model.PhotoOf(*identifier, data)).
		Return(identifier, nil)

	photoController := &PhotoController{mockPhotoService}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("photo", file.Name())
	part.Write(data)
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(identifier.Value())

	// Assertions
	if assert.NoError(t, photoController.Put(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestPutWithError(t *testing.T) {
	// setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Save(*identifier).
		Times(0)

	photoController := &PhotoController{mockPhotoService}

	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(identifier.Value())

	// Assertions
	assert.Error(t, photoController.Put(c))
}

func TestDelete(t *testing.T) {
	// setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Delete(*identifier).
		Return(nil)

	photoController := &PhotoController{mockPhotoService}

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(identifier.Value())

	// Assertions
	if assert.NoError(t, photoController.Delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteWithError(t *testing.T) {
	// setup
	identifier := model.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhotoService := mock_service.NewMockPhotoService(ctrl)
	mockPhotoService.EXPECT().
		Delete(*identifier).
		Return(errors.New("error"))

	photoController := &PhotoController{mockPhotoService}

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(identifier.Value())

	assert.Error(t, photoController.Delete(c))
}
