package errors

import (
	"fmt"
)

type NotFoundError struct {
	resource string
}

func NotFound(resource string) error {
	return NotFoundError{resource: resource}
}

func (error NotFoundError) Error() string {
	return fmt.Sprintf("resource not found. such id as %s", error.resource)
}
