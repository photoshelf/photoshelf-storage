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

type CustomContext struct {
	echo.Context
	conf Configuration
}

type Created struct {
	Id string
}

type PhotoList struct {
	 Ids []string
}

type Configuration struct {
	Server struct {
		Port int
	}
	Storage struct {
		Directory string
	}
}

func get(c echo.Context) error {
	cc := c.(*CustomContext)

	file, err := os.Open(path.Join(cc.conf.Storage.Directory, c.Param("id")))
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
}

func list(c echo.Context) error {
	cc := c.(*CustomContext)

	files, err := ioutil.ReadDir(cc.conf.Storage.Directory)
	if err != nil {
		log.Error(err)
		return err
	}
	names := []string{}
	for _, file := range files {
		names = append(names, file.Name())
	}

	return c.JSON(http.StatusOK, PhotoList{names})
}

func post(c echo.Context) error {
	cc := c.(*CustomContext)

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
	dst, err := os.Create(path.Join(cc.conf.Storage.Directory, filename))
	if err != nil {
		log.Error(err)
		return err
	}
	defer dst.Close()

	if _, err := dst.Write(data); err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusCreated, Created{filename})
}

func put(c echo.Context) error {
	cc := c.(*CustomContext)

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

	dst, err := os.Create(path.Join(cc.conf.Storage.Directory, c.Param("id")))
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
}

func delete(c echo.Context) error {
	cc := c.(*CustomContext)

	if err := os.Remove(path.Join(cc.conf.Storage.Directory, c.Param("id"))); err != nil {
		log.Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
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

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, configuration}
			return h(cc)
		}
	})

	e.GET("/:id", get)
	e.GET("/", list)
	e.POST("/", post)
	e.PUT("/:id", put)
	e.DELETE("/:id", delete)

	address := fmt.Sprintf(":%d", *port)
	e.Logger.Debug(e.Start(address))
}
