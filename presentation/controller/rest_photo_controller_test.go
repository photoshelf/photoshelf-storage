package controller

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/photoshelf/photoshelf-storage/application/mock_service"
	"github.com/photoshelf/photoshelf-storage/domain/model/photo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestRestPhotoController_Get(t *testing.T) {
	t.Run("when service no error, returns bytes", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Find(*identifier).
			Return(photo.Of(*identifier, readTestData(t)), nil)

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(identifier.Value())

		if assert.NoError(t, photoController.Get(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, readTestData(t), rec.Body.Bytes())
		}
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Find(*photo.IdentifierOf("not_found")).
			Return(nil, errors.New("error not found"))

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("not_found")

		assert.Error(t, photoController.Get(c))
	})

	t.Run("when service NotFoundError, returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Find(*photo.IdentifierOf("not_found")).
			Return(nil, &photo.ResourceError{ID: *photo.IdentifierOf("not_found"), Err: photo.ErrNotFound})

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("not_found")

		assert.NoError(t, photoController.Get(c))
	})
}

func TestRestPhotoController_Post(t *testing.T) {
	t.Run("when service no error, returns status created", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(gomock.Any()).
			Return(identifier, nil)

		photoController := &restPhotoControllerImpl{mockPhotoService}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("photo", identifier.Value())
		if err != nil {
			t.Fatal(err)
		}
		part.Write(readTestData(t))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, photoController.Post(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(gomock.Any()).
			Return(nil, errors.New("mock error"))

		photoController := &restPhotoControllerImpl{mockPhotoService}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("photo", identifier.Value())
		if err != nil {
			t.Fatal(err)
		}
		part.Write(readTestData(t))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.Error(t, photoController.Post(c))
	})

	t.Run("with nil body, returns error", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(*identifier).
			Times(0)

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.Error(t, photoController.Post(c))
	})
}

func TestRestPhotoController_Put(t *testing.T) {
	t.Run("when service no error, returns identifier", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(*photo.Of(*identifier, readTestData(t))).
			Return(identifier, nil)

		photoController := &restPhotoControllerImpl{mockPhotoService}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("photo", identifier.Value())
		if err != nil {
			t.Fatal(err)
		}
		part.Write(readTestData(t))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(identifier.Value())

		if assert.NoError(t, photoController.Put(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(*photo.Of(*identifier, readTestData(t))).
			Return(nil, errors.New("mock error"))

		photoController := &restPhotoControllerImpl{mockPhotoService}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("photo", identifier.Value())
		if err != nil {
			t.Fatal(err)
		}
		part.Write(readTestData(t))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(identifier.Value())

		assert.Error(t, photoController.Put(c))
	})

	t.Run("with nil body, returns error", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Save(*identifier).
			Times(0)

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(identifier.Value())

		assert.Error(t, photoController.Put(c))
	})

	t.Run("when no photo form, returns error", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("wrong_field", "filename")
		if err != nil {
			t.Fatal(err)
		}
		part.Write(readTestData(t))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_, err = readPhotoBytes(c)
		assert.Error(t, err)
	})
}

func TestRestPhotoController_Delete(t *testing.T) {
	t.Run("when service no error, returns status ok", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Delete(*identifier).
			Return(nil)

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(identifier.Value())

		if assert.NoError(t, photoController.Delete(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("when service error, returns error", func(t *testing.T) {
		identifier := photo.IdentifierOf("e3158990bdee63f8594c260cd51a011d")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPhotoService := mock_service.NewMockPhotoService(ctrl)
		mockPhotoService.EXPECT().
			Delete(*identifier).
			Return(errors.New("error"))

		photoController := &restPhotoControllerImpl{mockPhotoService}

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(identifier.Value())

		assert.Error(t, photoController.Delete(c))
	})
}

func readTestData(tb testing.TB) []byte {
	tb.Helper()

	filename := path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata", "e3158990bdee63f8594c260cd51a011d")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		tb.Fatal(err)
	}
	return body
}
