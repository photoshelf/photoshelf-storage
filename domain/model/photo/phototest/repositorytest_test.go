package phototest_test

import (
	"github.com/photoshelf/photoshelf-storage/domain/model/photo/phototest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomTestData(t *testing.T) {
	t.Run("it return same values ", func(t *testing.T) {
		first := phototest.RandomTestData(t)
		second := phototest.RandomTestData(t)

		assert.EqualValues(t, first, second)
	})

	t.Run("it have max size 20M", func(t *testing.T) {
		for _, data := range phototest.RandomTestData(t) {
			assert.True(t, len(data) < 20000000)
		}
	})
}
