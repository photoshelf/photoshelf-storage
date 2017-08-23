package container

import (
	"reflect"
	"regexp"
)

type singletonContainer struct {
	values map[string]interface{}
}

var sharedInstance *singletonContainer = &singletonContainer{make(map[string]interface{})}

func Get(name string) interface{} {
	return sharedInstance.values[name]
}

func Set(instance interface{}) {
	rep := regexp.MustCompile(`^\*.*\.`)
	name := rep.ReplaceAllString(reflect.TypeOf(instance).String(), "")
	sharedInstance.values[name] = instance
}
