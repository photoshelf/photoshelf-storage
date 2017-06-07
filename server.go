package main

import (
	"net/http"
	"os"
	"io"
	"io/ioutil"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"crypto/md5"
	"fmt"
	"time"
	"flag"
)

type Created struct {
	Id string
}

func main() {
	port := flag.Int("port", 1323, "port number")
	flag.Parse()

	e := echo.New()

	e.GET("/:id", func(c echo.Context) error {
		file, err := os.Open(c.Param("id"))
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

		data, err := ioutil.ReadAll(src)
		if err != nil {
			log.Error(err)
			return err
		}

		dataHash := fmt.Sprintf("%x", md5.Sum(data))
		filename := fmt.Sprintf("%x", md5.Sum([]byte(dataHash + time.Now().String())))
		dst, err := os.Create(filename)
		if err != nil {
			log.Error(err)
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			log.Error(err)
			return err
		}

		return c.JSON(http.StatusCreated, Created{filename})
	})

	e.PUT("/:id", func(c echo.Context) error {
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

		dst, err := os.Create(c.Param("id"))
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

	e.DELETE("/:id", func(c echo.Context) error {
		if err := os.Remove(c.Param("id")); err != nil {
			log.Error(err)
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	address := fmt.Sprintf(":%d", *port)
	e.Logger.Debug(e.Start(address))
}
