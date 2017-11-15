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
	ID  Identifier
	Err error
}

func (err *ResourceError) Error() string {
	return fmt.Sprintf("%s: %s", err.ID, err.Error())
}
