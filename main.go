package main

import (
	"fmt"
	"github.com/photoshelf/photoshelf-storage/router"
	"log"
	"os"
)

func main() {
	conf, err := configure()
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
