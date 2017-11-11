package router

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/log"
	"github.com/photoshelf/photoshelf-storage/infrastructure/container"
	"github.com/photoshelf/photoshelf-storage/presentation/mock_controller"
	"net"
	"net/http"
	"testing"
)

func TestLoad(t *testing.T) {

	port := randomPort(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	con := mock_controller.NewMockPhotoController(ctrl)
	con.EXPECT().Get(gomock.Any()).Times(1)
	con.EXPECT().Post(gomock.Any()).Times(1)
	con.EXPECT().Put(gomock.Any()).Times(1)
	con.EXPECT().Delete(gomock.Any()).Times(1)
	container.Set(con)

	e, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		log.Info(e.Start(fmt.Sprintf(":%s", port)))
	}()

	entrypoint := fmt.Sprintf("http://127.0.0.1:%s/photos/", port) + "%s"

	t.Run("route GET /photos/:id", func(t *testing.T) {
		_, err := http.Get(fmt.Sprintf(entrypoint, "test"))
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("route POST /photos/", func(t *testing.T) {
		_, err := http.Post(fmt.Sprintf(entrypoint, ""), "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("routes PUT /photos/:id", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(entrypoint, "test"), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("routes DELETE /photos/:id", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(entrypoint, "test"), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func randomPort(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := listener.Addr().String()
	_, port, err := net.SplitHostPort(addr)
	listener.Close()

	return port
}
