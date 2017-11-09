package phototest

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
			rand.Read(bytea)
			instance.randomData = append(instance.randomData, bytea)
		}
	}
	return instance.randomData
}
