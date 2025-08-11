package mapper

import (
	"reflect"
)

func IsPointerType(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer || t.Kind() == reflect.UnsafePointer || t.Kind() == reflect.Ptr
}
