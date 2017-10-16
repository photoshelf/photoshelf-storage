package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifierIsUnique(t *testing.T) {
	var ids = make([]*Identifier, 1000)
	for i := 0; i < 1000; i++ {
		id := NewIdentifier([]byte("hello world."))
		ids[i] = id
	}
	for i, id := range ids {
		var after = make([]string, 1000-(i+1))
		var before = make([]string, 1000-(len(after)+1))

		for i, id := range ids[i+1:] {
			after[i] = id.Value()
		}
		if i > 0 {
			for i, id := range ids[:i-1] {
				before[i] = id.Value()
			}
		}
		target := append(before, after...)

		assert.NotContains(t, id.value, target)
	}
}
