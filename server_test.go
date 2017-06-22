package main

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestGet(t *testing.T) {
	// Setup
	conf := Configuration{
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
	body, _ := os.Open(path.Join(conf.Storage.Directory, "e3158990bdee63f8594c260cd51a011d"))
	data, _ := ioutil.ReadAll(body)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("e3158990bdee63f8594c260cd51a011d")
	cc := &CustomContext{c, conf}

	// Assertions
	if assert.NoError(t, get(cc)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, data, rec.Body.Bytes())
	}
}
