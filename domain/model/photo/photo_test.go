package photo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	instance := New([]byte("image"))
	assert.Equal(t, []byte("image"), instance.image)
	assert.Equal(t, 0, len(instance.id.value))
}

func TestOf(t *testing.T) {
	instance := Of(*IdentifierOf("id"), []byte("image"))
	assert.Equal(t, []byte("image"), instance.image)
	assert.Equal(t, "id", instance.id.value)
}

func TestPhoto_ID(t *testing.T) {
	instance := Of(*IdentifierOf("id"), []byte("image"))
	actual := instance.ID()
	assert.Equal(t, "id", actual.value)
}

func TestPhoto_Image(t *testing.T) {
	instance := Of(*IdentifierOf("id"), []byte("image"))
	actual := instance.Image()
	assert.Equal(t, []byte("image"), actual)
}

func TestPhoto_IsNew(t *testing.T) {
	t.Run("without id ,returns false", func(t *testing.T) {
		instance := New([]byte("image"))
		assert.True(t, instance.IsNew())
	})

	t.Run("with id, returns false", func(t *testing.T) {
		instance := Of(*IdentifierOf("id"), []byte("image"))
		assert.False(t, instance.IsNew())
	})
}
