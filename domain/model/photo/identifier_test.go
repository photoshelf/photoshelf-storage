package photo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_Value(t *testing.T) {
	t.Run("it is unique", func(t *testing.T) {
		var ids []*Identifier
		for i := 0; i < 1000; i++ {
			ids = append(ids, NewIdentifier([]byte("hello world.")))
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
	})
}

func ExampleIdentifier_Value() {
	id := IdentifierOf("example_id")
	fmt.Println(id.Value())
	// Output:
	// example_id
}
