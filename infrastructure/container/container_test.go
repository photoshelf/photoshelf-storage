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

type testInterface interface {
	Hoge()
}

type testInterfaceTwo interface {
	Hoge()
}

type testImpl struct {
	Name string
}

func (*testImpl) Hoge() {}

var testdataSet []testdata

func TestMain(m *testing.M) {
	testdataSet = []testdata{
		{new(string), "value"},
		{new(int), 1234},
		{new(float64), 12.34},
		{new(testtype), testtype{"hoge", "fuga"}},
	}

	for _, testdata := range testdataSet {
		Set(testdata.val)
	}
	Set(&testImpl{"hoge"})
	os.Exit(m.Run())
}

func TestContainer(t *testing.T) {
	type hoge string

	testdataSet = append(testdataSet, testdata{new(hoge), nil})
	testdataSet = append(testdataSet, testdata{new(testImpl), testImpl{"hoge"}})
	testdataSet = append(testdataSet, testdata{new(testInterface), &testImpl{"hoge"}})
	testdataSet = append(testdataSet, testdata{new(testInterfaceTwo), nil})

	for _, testdata := range testdataSet {
		t.Run(fmt.Sprintf("type=%s", reflect.TypeOf(testdata.emp).String()), func(t *testing.T) {
			a := testdata.emp
			Get(a)

			if testdata.val != nil {
				actual := reflect.Indirect(reflect.ValueOf(a)).Interface()
				expect := reflect.ValueOf(testdata.val).Interface()

				assert.EqualValues(t, expect, actual)
			} else {
				assert.EqualValues(t, testdata.emp, a)
			}
		})
	}
}
