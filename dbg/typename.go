package dbg

import (
	"reflect"
)

func TypeName(obj interface{}) string {
	return reflect.ValueOf(obj).Type().String()
}
