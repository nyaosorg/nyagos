package dbg

import (
	"reflect"
)

func TypeName(obj interface{}) string {
	return reflect.TypeOf(obj).String()
}
