package main

import (
	"net/http"
	"os"
	"io"
	"io/ioutil"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	e.GET("/:filename", func(c echo.Context) error {
		file, err := os.Open(c.Param("filename"))
		if err != nil {
			log.Error(err)
			return err
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error(data)
			return err
		}
		mimeType := http.DetectContentType(data)
		return c.Blob(http.StatusOK, mimeType, data)
	})

	e.POST("/", func(c echo.Context) error {
		photo, err := c.FormFile("photo")
		if err != nil {
			log.Error(err)
			return err
		}
		src, err := photo.Open()
		if err != nil {
			log.Error(err)
			return err
		}
		defer src.Close()

		dst, err := os.Create(photo.Filename)
		if err != nil {
			log.Error(err)
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			log.Error(err)
			return err
		}

		return c.NoContent(http.StatusCreated)
	})

	e.PUT("/:filename", func(c echo.Context) error {
		photo, err := c.FormFile("photo")
		if err != nil {
			log.Error(err)
			return err
		}
		src, err := photo.Open()
		if err != nil {
			log.Error(err)
			return err
		}
		defer src.Close()

		dst, err := os.Create(c.Param("filename"))
		if err != nil {
			log.Error(err)
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			log.Error(err)
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/:filename", func(c echo.Context) error {
		if err := os.Remove(c.Param("filename")); err != nil {
			log.Error(err)
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Debug(e.Start(":1323"))
}
