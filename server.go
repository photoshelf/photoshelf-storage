package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

type Created struct {
	Id string
}

type Configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Directory string
	}
}

func main() {
	configurationFile, err := ioutil.ReadFile("./application.yml")
	if err != nil {
		log.Warn(err)
	}

	configuration := Configuration{}
	if err := yaml.Unmarshal(configurationFile, &configuration); err != nil {
		log.Fatal(err)
		return
	}

	port := flag.Int("port", configuration.Server.Port, "port number")
	flag.Parse()

	e := echo.New()

	e.GET("/:id", func(c echo.Context) error {
		file, err := os.Open(path.Join(configuration.Storage.Directory, c.Param("id")))
		if err != nil {
			log.Error(err)
			return err
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error(err)
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
		filename := fmt.Sprintf("%x", md5.Sum([]byte(dataHash+time.Now().String())))
		dst, err := os.Create(path.Join(configuration.Storage.Directory, filename))
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

		dst, err := os.Create(path.Join(configuration.Storage.Directory, c.Param("id")))
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
		if err := os.Remove(path.Join(configuration.Storage.Directory, c.Param("id"))); err != nil {
			log.Error(err)
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	address := fmt.Sprintf(":%d", *port)
	e.Logger.Debug(e.Start(address))
}
