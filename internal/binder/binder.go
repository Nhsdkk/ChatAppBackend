package binder

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"reflect"
	"regexp"
	"strconv"
)

type BindSource string

const BindSourceTag = "binder"

const (
	Path  BindSource = "path"
	Body             = "body"
	Query            = "query"
	Form             = "form"
	None             = "None"
)

var bindSourceMapping = map[string]BindSource{
	"path":  Path,
	"body":  Body,
	"query": Query,
	"form":  Form,
}

var bitSizes = map[reflect.Kind]int{
	reflect.Int8:    8,
	reflect.Int16:   16,
	reflect.Int32:   32,
	reflect.Int64:   64,
	reflect.Uint8:   8,
	reflect.Uint16:  16,
	reflect.Uint32:  32,
	reflect.Uint64:  64,
	reflect.Float64: 64,
	reflect.Float32: 32,
}

var tagRegexp = regexp.MustCompile("(?P<source>path|body|query|form)\\s*,\\s*(?P<fieldName>\\w+)")

func parseBinderTag(t reflect.StructField) (BindSource, string, error) {
	tagString, exists := t.Tag.Lookup(BindSourceTag)
	if !exists {
		return None, "", nil
	}

	sourceIdx := tagRegexp.SubexpIndex("source")
	fieldNameIdx := tagRegexp.SubexpIndex("fieldName")

	matches := tagRegexp.FindStringSubmatch(tagString)

	if len(matches) <= max(sourceIdx, fieldNameIdx) {
		return None, "", errors.New("invalid format ")
	}

	sourceString := matches[sourceIdx]
	source := bindSourceMapping[sourceString]

	fieldNameString := matches[fieldNameIdx]

	return source, fieldNameString, nil
}

func Bind[T interface{}](ctx *gin.Context) (interface{}, error) {
	t := reflect.TypeFor[T]()

	if t.Kind() != reflect.Struct {
		return nil, errors.New("can't bind request to non struct type")
	}

	jsonBody := make(map[string]json.RawMessage)
	body, readBodyError := io.ReadAll(ctx.Request.Body)
	if readBodyError == nil {
		_ = json.Unmarshal(body, &jsonBody)
	}

	v := reflect.New(t)

	for idx := range t.NumField() {
		fieldV := v.Elem().Field(idx)
		fieldT := t.Field(idx)

		if !fieldV.CanAddr() {
			return nil, errors.New(fmt.Sprintf("can't get addr of the field %s", fieldT.Name))
		}

		if !fieldV.CanSet() {
			return nil, errors.New(fmt.Sprintf("can't set value of the field %s", fieldT.Name))
		}

		source, sourceFieldName, tagParsingError := parseBinderTag(fieldT)

		if tagParsingError != nil {
			return nil, tagParsingError
		}

		switch source {
		case Body:
			jsonValue, ok := jsonBody[sourceFieldName]
			if !ok {
				break
			}

			if err := json.Unmarshal(jsonValue, fieldV.Addr().Interface()); err != nil {
				return nil, err
			}
		case Path:
			pathStringValue := ctx.Param(sourceFieldName)
			if pathStringValue == "" {
				return nil, errors.New(fmt.Sprintf("missing path value for field %s", fieldT.Name))
			}

			pathValue, convError := convertStringToReflectType(pathStringValue, fieldT.Type)

			if convError != nil {
				return nil, convError
			}

			fieldV.Set(reflect.ValueOf(pathValue))
		case Query:
			queryStringValue := ctx.Query(sourceFieldName)
			if queryStringValue == "" {
				return nil, errors.New(fmt.Sprintf("missing path value for field %s", fieldT.Name))
			}

			queryValue, convError := convertStringToReflectType(queryStringValue, fieldT.Type)

			if convError != nil {
				return nil, convError
			}

			fieldV.Set(reflect.ValueOf(queryValue))
		case Form:
			formStringValue := ctx.PostForm(sourceFieldName)
			if formStringValue == "" {
				return nil, errors.New(fmt.Sprintf("missing form value for field %s", fieldT.Name))
			}

			formValue, convError := convertStringToReflectType(formStringValue, fieldT.Type)

			if convError != nil {
				return nil, convError
			}

			fieldV.Set(reflect.ValueOf(formValue))
		}
	}

	return v.Elem().Interface(), nil
}

func convertStringToReflectType(v string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Int:
		return strconv.Atoi(v)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.ParseInt(v, 10, bitSizes[t.Kind()])
	case reflect.Float64, reflect.Float32:
		return strconv.ParseFloat(v, bitSizes[t.Kind()])
	case reflect.Uint:
		ui64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		return uint(ui64), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.ParseUint(v, 10, bitSizes[t.Kind()])
	case reflect.Bool:
		return strconv.ParseBool(v)
	case reflect.String:
		return v, nil
	case reflect.UnsafePointer, reflect.Pointer, reflect.Uintptr:
		return nil, errors.New(fmt.Sprintf("can't convert to %s", t.Kind()))
	default:
		if t.Implements(reflect.TypeFor[encoding.TextUnmarshaler]()) {
			result := reflect.New(t).Interface().(encoding.TextUnmarshaler)
			if err := result.UnmarshalText([]byte(v)); err != nil {
				return nil, err
			}

			return result, nil
		}

		return nil, errors.New(fmt.Sprintf("can't convert to %s", t.Kind()))
	}

}
