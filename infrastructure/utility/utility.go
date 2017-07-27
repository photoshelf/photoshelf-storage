package utility

import (
	"io/ioutil"
	"mime/multipart"
)

func Read(header multipart.FileHeader) ([]byte, error) {
	src, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}
