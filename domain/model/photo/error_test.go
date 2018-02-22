package photo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceError_Error(t *testing.T) {
	t.Run("error not found", func(t *testing.T) {
		e := &ResourceError{*IdentifierOf("id"), ErrNotFound}
		assert.Equal(t, "id: id does not exists", e.Error())
	})

	t.Run("photo can't read", func(t *testing.T) {
		e := &ResourceError{*IdentifierOf("id"), ErrCannotRead}
		assert.Equal(t, "id: photo can't read", e.Error())
	})
}
