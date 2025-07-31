package mapper

import (
	"chat_app_backend/internal/mapper/utils"
	"errors"
	"fmt"
	"reflect"
)

const MapperTag = "mapper"
const ExcludeTagValue = "exclude"

type IDest struct{}

type IMapper interface {
	Map(dest interface{}, srcs ...interface{}) error
}

type Mapper struct{}

func (m Mapper) Map(dest interface{}, srcs ...interface{}) error {
	srcVals := make([]reflect.Value, 0)
	valNamesExistence := make(map[string]bool)

	for _, src := range srcs {
		srcVal := reflect.ValueOf(src)
		if mapper.IsPointerType(srcVal.Type()) {
			return errors.New("src should not be a pointer")
		}

		for fieldIdx := range srcVal.NumField() {
			field := srcVal.Type().Field(fieldIdx)

			tag, ok := field.Tag.Lookup(MapperTag)
			if ok {
				if tag != ExcludeTagValue {
					panic(fmt.Sprintf("malformed tag on field with name %s", field.Name))
				}
				continue
			}

			if _, ok := valNamesExistence[field.Name]; ok {
				return errors.New(fmt.Sprintf("name collision found as field %s exists in more than one struct", field.Name))
			}
			valNamesExistence[field.Name] = true
		}

		srcVals = append(srcVals, srcVal)
	}

	destVal := reflect.ValueOf(dest)

	if !mapper.IsPointerType(destVal.Type()) {
		return errors.New("dest should be a pointer")
	}

	return mapStruct(destVal, srcVals...)
}

func mapStruct(destVal reflect.Value, srcVals ...reflect.Value) error {
	if !mapper.IsPointerType(destVal.Type()) {
		return errors.New("dest is not a pointer")
	}

	destVal = destVal.Elem()

	for fieldIdx := range destVal.NumField() {
		destField := destVal.Field(fieldIdx)
		srcField, err := findValue(destVal.Type().Field(fieldIdx), srcVals...)

		if err != nil {
			return err
		}

		if !destField.CanSet() {
			continue
		}

		if mapper.IsArrayType(destField.Type()) {
			mapArray(*srcField, destField.Addr())
			continue
		}

		if mapper.IsSliceType(destField.Type()) {
			mapSlice(*srcField, destField.Addr())
			continue
		}

		destField.Set(*srcField)
	}

	return nil
}

func findValue(field reflect.StructField, srcVals ...reflect.Value) (*reflect.Value, error) {
	for _, srcVal := range srcVals {
		srcFieldV := srcVal.FieldByName(field.Name)

		if !srcFieldV.IsValid() {
			continue
		}

		srcFieldT, _ := srcVal.Type().FieldByName(field.Name)
		tag, ok := srcFieldT.Tag.Lookup(MapperTag)
		if ok {
			if tag != ExcludeTagValue {
				panic(fmt.Sprintf("malformed tag on field with name %s", srcFieldT.Name))
			}
			continue
		}

		if !mapper.SameTypes(srcFieldV.Type(), field.Type) {
			return nil, errors.New(fmt.Sprintf("kinds of values does not match (%s and %s)", srcFieldV.Kind(), field.Type.Kind()))
		}

		if !srcFieldV.CanInterface() {
			return nil, errors.New(fmt.Sprintf("src field %s is unexported", field.Name))
		}

		return &srcFieldV, nil
	}

	return nil, errors.New(fmt.Sprintf("src does not have field %s, which dest has", field.Name))
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
