package mapper

import (
	"chat_app_backend/internal/mapper/utils"
	"errors"
	"fmt"
	"reflect"
)

type IDest struct{}

type IMapper interface {
	Map(src interface{}, dest interface{}) error
}

type Mapper struct{}

func (m Mapper) Map(src interface{}, dest interface{}) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if !mapper.IsPointerType(destVal.Type()) {
		return errors.New("dest should be a pointer")
	}

	if mapper.IsPointerType(srcVal.Type()) {
		return errors.New("src should not be a pointer")
	}

	return mapStruct(srcVal, destVal)
}

func mapStruct(srcVal reflect.Value, destVal reflect.Value) error {
	if !mapper.IsPointerType(destVal.Type()) {
		return errors.New("dest is not a pointer")
	}

	destVal = destVal.Elem()

	for fieldIdx := range destVal.NumField() {
		destField := destVal.Field(fieldIdx)
		srcField := srcVal.FieldByName(destVal.Type().Field(fieldIdx).Name)

		if !srcField.IsValid() {
			return errors.New(fmt.Sprintf("src does not have field %s, which dest has", destVal.Type().Field(fieldIdx).Name))
		}

		if !mapper.SameTypes(srcField.Type(), destField.Type()) {
			return errors.New(fmt.Sprintf("kinds of values does not match (%s and %s", srcVal.Kind(), destVal.Kind()))
		}

		if !destField.CanSet() {
			continue
		}

		if mapper.IsArrayType(destField.Type()) {
			mapArray(srcField, destField.Addr())
			continue
		}

		if mapper.IsSliceType(destField.Type()) {
			mapSlice(srcField, destField.Addr())
			continue
		}

		destField.Set(srcField)
	}

	return nil
}

func mapArray(srcVal reflect.Value, destVal reflect.Value) {
	array := reflect.New(reflect.ArrayOf(srcVal.Len(), srcVal.Type().Elem()))
	for idx := range srcVal.Len() {
		array.Elem().Index(idx).Set(srcVal.Index(idx))
	}
	destVal.Elem().Set(array.Elem())
}

func mapSlice(srcVal reflect.Value, destVal reflect.Value) {
	slice := reflect.MakeSlice(srcVal.Type(), srcVal.Len(), srcVal.Len())
	reflect.Copy(slice, srcVal)
	destVal.Elem().Set(slice)
}
