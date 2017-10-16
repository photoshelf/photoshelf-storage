package main

import (
	"fmt"
	"github.com/photoshelf/photoshelf-storage/application"
	"github.com/photoshelf/photoshelf-storage/presentation/router"
	"log"
	"os"
)

func main() {
	conf, err := application.Configure()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	e, err := router.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	address := fmt.Sprintf(":%d", conf.Server.Port)
	e.Logger.Debug(e.Start(address))
}
