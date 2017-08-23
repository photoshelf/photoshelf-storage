package main

import (
	"github.com/photoshelf/photoshelf-storage/router"
	"fmt"
	"log"
	"os"
)

func main() {
	conf, err := Configure()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	e := router.Load()

	address := fmt.Sprintf(":%d", conf.Server.Port)
	e.Logger.Debug(e.Start(address))
}
