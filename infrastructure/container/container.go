package container

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"reflect"
)

type singletonContainer struct {
	values map[string]interface{}
}

var instance *singletonContainer

func init() {
	instance = &singletonContainer{
		make(map[string]interface{}),
	}
}

func Set(val interface{}) {
	key := reflect.ValueOf(val).Type().String()
	log.Info(fmt.Sprintf("added %s to components container.", key))
	instance.values[key] = val
}

func Get(ptr interface{}) {
	val := reflect.ValueOf(ptr)
	key := reflect.Indirect(val).Type().String()
	component := instance.values[key]
	if component == nil {
		log.Warn(fmt.Sprintf("component not found. such type of %s.", key))
		return
	}
	log.Info(fmt.Sprintf("found component of %s .", key))

	elm := reflect.ValueOf(ptr).Elem()
	elm.Set(reflect.ValueOf(component))
}
