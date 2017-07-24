package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
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

var conf = Configuration{
	Server: struct {
		Port int
	}{
		1234,
	},
	Storage: struct {
		Directory string
	}{
		"testdata",
	},
}

func TestGet(t *testing.T) {
	// Setup
	body, _ := os.Open(path.Join(conf.Storage.Directory, "e3158990bdee63f8594c260cd51a011d"))
	data, _ := ioutil.ReadAll(body)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("e3158990bdee63f8594c260cd51a011d")
	cc := &CustomContext{c, conf, nil}

	// Assertions
	if assert.NoError(t, get(cc)) {
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
	cc := &CustomContext{c, conf, nil}

	// Assertions
	assert.Error(t, get(cc))
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
	cc := &CustomContext{c, conf, nil}

	// Assertions
	assert.Error(t, get(cc))
}

func TestPost(t *testing.T) {
	file, _ := os.Open(path.Join(conf.Storage.Directory, "e3158990bdee63f8594c260cd51a011d"))
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
	cc := &CustomContext{c, conf, nil}

	err := post(cc)

	var res map[string]string
	json.Unmarshal(rec.Body.Bytes(), &res)

	actualFile, _ := os.Open(path.Join(conf.Storage.Directory, res["Id"]))
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
	cc := &CustomContext{c, conf, nil}

	// Assertions
	assert.Error(t, post(cc))
}

func TestPut(t *testing.T) {
	file, _ := os.Open(path.Join(conf.Storage.Directory, "e3158990bdee63f8594c260cd51a011d"))
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
	cc := &CustomContext{c, conf, nil}

	err := put(cc)

	actualFile, _ := os.Open(path.Join(conf.Storage.Directory, "test"))
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
	cc := &CustomContext{c, conf, nil}

	// Assertions
	assert.Error(t, post(cc))
}

func TestDelete(t *testing.T) {
	src, _ := os.Open(path.Join(conf.Storage.Directory, "e3158990bdee63f8594c260cd51a011d"))
	src.Close()

	dst, _ := os.Create(path.Join(conf.Storage.Directory, "test"))
	dst.Close()

	io.Copy(dst, src)

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("test")
	cc := &CustomContext{c, conf, nil}

	err := delete(cc)
	_, exist := os.Stat(path.Join(conf.Storage.Directory, "test"))

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
	cc := &CustomContext{c, conf, nil}

	assert.Error(t, delete(cc))
}
