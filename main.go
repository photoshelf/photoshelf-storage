package main

import (
	"log"
	"os"
	"net"
	"google.golang.org/grpc"
	"context"
	"github.com/photoshelf/photoshelf-storage/presentation"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	//conf, err := application.Configure(os.Args[1:]...)
	//if err != nil {
	//	log.Fatal(err)
	//	os.Exit(-1)
	//}
	//e, err := router.Load()
	//if err != nil {
	//	log.Fatal(err)
	//	os.Exit(-1)
	//}
	//
	//address := fmt.Sprintf(":%d", conf.Server.Port)
	//e.Logger.Debug(e.Start(address))

	lis, err := net.Listen("tcp", ":1323")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	db, err := leveldb.OpenFile("leveldb", nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	s := grpc.NewServer()
	presentation.RegisterPhotoServiceServer(s, &server{db: db})
	go func() {
		s.Serve(lis)
	}()

	conn, err := grpc.Dial(":1323", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := presentation.NewPhotoServiceClient(conn)

	//id := "hoge"
	//_, err = c.Find(context.Background(), &presentation.Id{Value: &id})
	//if err != nil {
	//	log.Fatal(err)
	//}

	image := []byte("Hello World.")
	identifier, err := c.Save(context.Background(), &presentation.Photo{Image: image})
	log.Print(*identifier.Value)

	newid := "newid"
	photo, err := c.Find(context.Background(), &presentation.Id{Value: &newid})
	log.Print(photo)
}

type server struct {
	db *leveldb.DB
}

func (s *server) Save(ctx context.Context, in *presentation.Photo) (*presentation.Id, error) {
	var id *string
	if in.Id == nil {
		newId := "newid"
		id = &newId
	} else {
		id = in.Id.Value
	}

	s.db.Put([]byte(*id), in.Image, nil)
	return &presentation.Id{Value: id}, nil
}

func (s *server) Find(ctx context.Context, in *presentation.Id) (*presentation.Photo, error) {
	image, err := s.db.Get([]byte(*in.Value), nil)
	if err != nil {
		return nil, err
	}
	return &presentation.Photo{Id: in, Image: image}, nil
}
