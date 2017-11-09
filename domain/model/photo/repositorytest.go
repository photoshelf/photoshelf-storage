package photo

import (
	"math/rand"
	"testing"
)

type testDataSet struct {
	randomData [][]byte
}

var instance *testDataSet

func init() {
	instance = &testDataSet{}
}

func RandomTestData(tb testing.TB) [][]byte {
	if instance.randomData == nil {
		max := 20000000 // 20MB
		for i := 0; i < 20; i++ {
			size := rand.Intn(max)
			bytea := make([]byte, size)
			if _, err := rand.Read(bytea); err != nil {
				tb.Fatal(err)
			}
			instance.randomData = append(instance.randomData, bytea)
		}
	}
	return instance.randomData
}
