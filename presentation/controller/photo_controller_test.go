package controller

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/photoshelf/photoshelf-storage/application/service"
	"github.com/photoshelf/photoshelf-storage/infrastructure/datastore"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var storagePath = path.Join(os.Getenv("GOPATH"), "src/github.com/photoshelf/photoshelf-storage", "testdata")

var repository = datastore.NewFileStorage(storagePath)
var photoService = service.NewPhotoService(repository)
var photoController = NewPhotoController(*photoService)

func TestGet(t *testing.T) {
	// Setup
	body, _ := os.Open(path.Join(storagePath, "e3158990bdee63f8594c260cd51a011d"))
	data, _ := ioutil.ReadAll(body)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("e3158990bdee63f8594c260cd51a011d")

	// Assertions
	if assert.NoError(t, photoController.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, data, rec.Body.Bytes())
	}
}

func TestGetNotFound(t *testing.T) {
	// Setup
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

func TestGetDirectory(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("dir")

	// Assertions
	assert.Error(t, photoController.Get(c))
}

func TestPost(t *testing.T) {
	file, _ := os.Open(path.Join(storagePath, "e3158990bdee63f8594c260cd51a011d"))
	data, _ := ioutil.ReadAll(file)
	file.Close()

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

	err := photoController.Post(c)

	var res map[string]string
	json.Unmarshal(rec.Body.Bytes(), &res)

	actualFile, _ := os.Open(path.Join(storagePath, res["Id"]))
	actual, _ := ioutil.ReadAll(actualFile)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, actual, data)
	}
}

func TestPostWithoutData(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	assert.Error(t, photoController.Post(c))
}

func TestPut(t *testing.T) {
	file, _ := os.Open(path.Join(storagePath, "e3158990bdee63f8594c260cd51a011d"))
	data, _ := ioutil.ReadAll(file)
	file.Close()

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
	c.SetParamValues("test")

	err := photoController.Put(c)

	actualFile, _ := os.Open(path.Join(storagePath, "test"))
	actual, _ := ioutil.ReadAll(actualFile)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, actual, data)
	}
}

func TestPutWithoutData(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("e3158990bdee63f8594c260cd51a011d")

	// Assertions
	assert.Error(t, photoController.Put(c))
}

func TestDelete(t *testing.T) {
	src, _ := os.Open(path.Join(storagePath, "e3158990bdee63f8594c260cd51a011d"))
	src.Close()

	dst, _ := os.Create(path.Join(storagePath, "test"))
	dst.Close()

	io.Copy(dst, src)

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("test")

	err := photoController.Delete(c)
	_, exist := os.Stat(path.Join(storagePath, "test"))

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Error(t, exist)
	}
}

func TestDeleteWithoutFile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("test")

	assert.Error(t, photoController.Delete(c))
}
