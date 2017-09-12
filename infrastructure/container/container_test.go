package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

type testtype struct {
	Foo string
	Bar string
}

type testdata struct {
	emp interface{}
	val interface{}
}

var testdataSet = []testdata{
	{new(string), "value"},
	{new(int), 1234},
	{new(float64), 12.34},
	{new(testtype), testtype{"hoge", "fuga"}},
}

func TestMain(m *testing.M) {
	for _, testdata := range testdataSet {
		Set(testdata.val)
	}
	os.Exit(m.Run())
}

func TestContainer(t *testing.T) {
	for _, testdata := range testdataSet {
		t.Run(fmt.Sprintf("type=%s", reflect.TypeOf(testdata.emp).Name()), func(t *testing.T) {
			a := testdata.emp
			Get(a)

			actual := reflect.Indirect(reflect.ValueOf(a)).Interface()
			expect := reflect.Indirect(reflect.ValueOf(testdata.val)).Interface()

			assert.EqualValues(t, expect, actual)
		})
	}
}
