package main

import (
	"fmt"
	"github.com/photoshelf/photoshelf-storage/application"
	"github.com/photoshelf/photoshelf-storage/presentation/router"
	"log"
	"net"
	"os"
)

func main() {
	conf, err := application.Configure(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	address := fmt.Sprintf(":%d", conf.Server.Port)

	switch conf.Server.Mode {
	case "rest":
		e, err := router.LoadEchoServer()
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
		e.Logger.Info(e.Start(address))

	case "grpc":
		s := router.LoadGrpcServer()

		listener, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
		s.Serve(listener)

	default:
		log.Fatalf("No such as server mode: %s", conf.Server.Mode)
		os.Exit(-1)
	}
}
