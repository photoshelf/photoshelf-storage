package main

import (
	"fmt"
	"github.com/photoshelf/photoshelf-storage/application"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/infrastructure/protobuf"
	"github.com/photoshelf/photoshelf-storage/presentation/controller"
	"github.com/photoshelf/photoshelf-storage/presentation/router"
	"google.golang.org/grpc"
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
	switch conf.Server.Mode {
	case "rest":
		e, err := router.Load()
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}

		address := fmt.Sprintf(":%d", conf.Server.Port)
		e.Logger.Debug(e.Start(address))

	case "grpc":
		listener, err := net.Listen("tcp", ":1323")
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}

		photoServiceServer := controller.NewGrpcPhotoController()
		container.Get(&photoServiceServer)

		s := grpc.NewServer()
		protobuf.RegisterPhotoServiceServer(s, photoServiceServer)
		s.Serve(listener)

	default:
		log.Fatalf("No such as server mode: %s", conf.Server.Mode)
		os.Exit(-1)
	}
}
