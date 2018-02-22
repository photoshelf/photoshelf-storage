package photo

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("id does not exists")
	ErrCannotRead = errors.New("photo can't read")
)

type ResourceError struct {
	Id  Identifier
	Err error
}

func (err *ResourceError) Error() string {
	return fmt.Sprintf("%s: %s", err.Id.value, err.Err.Error())
}
