package mapper

import (
	"reflect"
)

func IsPointerType(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer || t.Kind() == reflect.UnsafePointer || t.Kind() == reflect.Ptr
}

func IsSliceType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

func IsArrayType(t reflect.Type) bool {
	return t.Kind() == reflect.Array
}

func SameTypes(t1 reflect.Type, t2 reflect.Type) bool {
	return t1.Name() == t2.Name()
}
