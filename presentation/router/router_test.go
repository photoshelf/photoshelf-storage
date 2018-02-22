package router

import (
	"github.com/golang/mock/gomock"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/presentation/mock_controller"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadEchoServer(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	con := mock_controller.NewMockPhotoController(ctrl)
	con.EXPECT().Get(gomock.Any()).Times(1)
	con.EXPECT().Post(gomock.Any()).Times(1)
	con.EXPECT().Put(gomock.Any()).Times(1)
	con.EXPECT().Delete(gomock.Any()).Times(1)
	container.Set(con)

	e, err := LoadEchoServer()
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(e)
	client := server.Client()

	t.Run("route GET /photos/:id", func(t *testing.T) {
		_, err := client.Get(server.URL + "/photos/test")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("route POST /photos/", func(t *testing.T) {
		_, err := client.Post(server.URL+"/photos/", "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("routes PUT /photos/:id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, server.URL+"/photos/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("routes DELETE /photos/:id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/photos/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestLoadGrpcServer(t *testing.T) {
	s := LoadGrpcServer()

	assert.IsType(t, &grpc.Server{}, s)
}
