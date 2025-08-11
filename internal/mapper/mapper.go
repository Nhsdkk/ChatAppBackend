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
					panic(fmt.Sprintf("malformed mapper tag on field with name %s (expected %s, but got %s)", field.Name, ExcludeTagValue, tag))
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

		var mappingError error

		switch {
		case destField.Type() == srcField.Type():
			destField.Set(*srcField)
		case srcField.CanConvert(destField.Type()):
			destField.Set(srcField.Convert(destField.Type()))
		case destField.Kind() == reflect.Array:
			mappingError = mapArray(*srcField, destField.Addr())
		case destField.Kind() == reflect.Slice:
			mappingError = mapSlice(*srcField, destField.Addr())
		case destField.Kind() == reflect.Struct:
			mappingError = mapStruct(destField.Addr(), *srcField)
		}

		if mappingError != nil {
			return mappingError
		}
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
				panic(fmt.Sprintf("malformed mapper tag on field with name %s (expected %s, but got %s)", field.Name, ExcludeTagValue, tag))
			}
			continue
		}

		if !srcFieldV.CanInterface() {
			return nil, errors.New(fmt.Sprintf("src field %s is unexported", field.Name))
		}

		return &srcFieldV, nil
	}

	return nil, errors.New(fmt.Sprintf("src does not have field %s, which dest has", field.Name))
}

func mapArray(srcVal reflect.Value, destVal reflect.Value) error {
	array := reflect.New(reflect.ArrayOf(srcVal.Len(), destVal.Type().Elem()))
	for idx := range srcVal.Len() {
		element := array.Elem().Index(idx)
		src := srcVal.Index(idx)

		var mappingError error

		switch {
		case src.Type() == element.Type():
			element.Set(src)
		case src.CanConvert(element.Type()):
			element.Set(src.Convert(element.Type()))
		case element.Kind() == reflect.Array:
			mappingError = mapArray(src, element.Addr())
		case element.Kind() == reflect.Slice:
			mappingError = mapSlice(src, element.Addr())
		case element.Kind() == reflect.Struct:
			mappingError = mapStruct(element.Addr(), srcVal)
		}

		if mappingError != nil {
			return mappingError
		}
	}

	destVal.Elem().Set(array.Elem())
	return nil
}

func mapSlice(srcVal reflect.Value, destVal reflect.Value) error {
	slice := reflect.MakeSlice(destVal.Type().Elem(), srcVal.Len(), srcVal.Len())
	for idx := range srcVal.Len() {
		element := slice.Index(idx)
		src := srcVal.Index(idx)

		var mappingError error

		switch {
		case src.Type() == element.Type():
			element.Set(src)
		case src.CanConvert(element.Type()):
			element.Set(src.Convert(element.Type()))
		case element.Kind() == reflect.Array:
			mappingError = mapArray(src, element.Addr())
		case element.Kind() == reflect.Slice:
			mappingError = mapSlice(src, element.Addr())
		case element.Kind() == reflect.Struct:
			mappingError = mapStruct(element.Addr(), srcVal)
		}

		if mappingError != nil {
			return mappingError
		}
	}

	destVal.Elem().Set(slice)
	return nil
}
