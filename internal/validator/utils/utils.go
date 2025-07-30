package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagName = "validator"
const validationsSeparator = ";"
const wordSeparator = " "
const arraySeparator = ","

var arrayRegexp = regexp.MustCompile("\\[(?P<inside>(?:[^\\[\\]\\,]+\\,)*(?:[^\\[\\]\\,]+){1})\\]")

type ValidationTypes string

const (
	NotEmpty       ValidationTypes = "not_empty"
	Greater                        = "gt"
	Less                           = "lt"
	GreaterOrEqual                 = "gte"
	LessOrEqual                    = "lte"
	OneOf                          = "one_of"
	Equal                          = "eq"
	Length                         = "length"
)

var validationTypeConverterMap = map[string]ValidationTypes{
	"not_empty": NotEmpty,
	"gt":        Greater,
	"lt":        Less,
	"gte":       GreaterOrEqual,
	"lte":       LessOrEqual,
	"eq":        Equal,
	"one_of":    OneOf,
	"length":    Length,
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

func parseValidationType(valType string) (*ValidationTypes, error) {
	value, ok := validationTypeConverterMap[valType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("can't parse validation type from %s", valType))
	}
	return &value, nil
}

type Validation struct {
	ValidationType ValidationTypes
	Arguments      []reflect.Value
}

type Tag struct {
	Validations []Validation
}

func ParseTag(t reflect.StructField) (*Tag, error) {
	fieldType := t.Type
	fieldTag, exists := t.Tag.Lookup(tagName)
	if !exists {
		return nil, nil
	}

	validations, validationsParseError := parseValidations(fieldTag, fieldType)

	if validationsParseError != nil {
		return nil, validationsParseError
	}

	return &Tag{
		Validations: validations,
	}, nil
}

func parseValidations(tagString string, t reflect.Type) ([]Validation, error) {
	validations := make([]Validation, 0)

	stringValidations := strings.Split(tagString, validationsSeparator)

	for _, stringValidation := range stringValidations {
		if len(stringValidation) == 0 {
			continue
		}

		words := strings.Split(stringValidation, wordSeparator)

		validationType, validationTypeParseError := parseValidationType(words[0])
		if validationTypeParseError != nil {
			return nil, validationTypeParseError
		}

		arguments := make([]reflect.Value, 0)

		switch *validationType {
		case NotEmpty:
			break
		case Greater, GreaterOrEqual, Less, LessOrEqual, Equal:
			arg, err := parseComparison(words, t)
			if err != nil {
				return nil, err
			}
			arguments = arg
			break
		case OneOf:
			arg, err := parseOneOf(words, t)
			if err != nil {
				return nil, err
			}
			arguments = arg
		case Length:
			arg, err := parseLength(words)
			if err != nil {
				return nil, err
			}
			arguments = arg
		}

		validations = append(validations, Validation{
			ValidationType: *validationType,
			Arguments:      arguments,
		})
	}

	return validations, nil
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
		return convertStringToReflectType(v, t.Elem())
	default:
		return nil, errors.New(fmt.Sprintf("can't convert %s to struct", t.Kind()))
	}

}

func parseComparison(stringArguments []string, t reflect.Type) ([]reflect.Value, error) {
	if len(stringArguments[1:]) != 1 {
		return nil, errors.New(fmt.Sprintf("comparison validation requires exactly one argument, but %d found", len(stringArguments[1:])))
	}

	argument, err := convertStringToReflectType(stringArguments[1], t)
	if err != nil {
		return nil, err
	}
	return []reflect.Value{reflect.ValueOf(argument)}, nil
}

func parseLength(stringArguments []string) ([]reflect.Value, error) {
	if len(stringArguments[1:]) != 2 {
		return nil, errors.New(fmt.Sprintf("length validation requires exactly 2 arguments, but %d found", len(stringArguments)))
	}

	comparisonType, err := parseValidationType(stringArguments[1])
	if err != nil {
		return nil, err
	}

	switch *comparisonType {
	case OneOf, NotEmpty, Length:
		return nil, errors.New(fmt.Sprintf("length accepts comparison validator as first argument, but found %s validator, which is not", *comparisonType))
	default:
		break
	}

	uintValue, err := convertStringToReflectType(stringArguments[2], reflect.TypeFor[int]())

	return []reflect.Value{reflect.ValueOf(stringArguments[1]), reflect.ValueOf(uintValue)}, nil
}

func parseOneOf(stringArguments []string, t reflect.Type) ([]reflect.Value, error) {
	if len(stringArguments[1:]) != 1 {
		return nil, errors.New(fmt.Sprintf("comparison validation requires exactly one argument, but %d found", len(stringArguments[1:])))
	}

	argument := stringArguments[1]

	index := arrayRegexp.SubexpIndex("inside")
	innerArray := arrayRegexp.FindStringSubmatch(argument)
	if index >= len(innerArray) {
		return nil, errors.New("array argument does not match the style or is empty")
	}

	arguments := make([]reflect.Value, 0)

	for _, item := range strings.Split(innerArray[index], arraySeparator) {
		val, err := convertStringToReflectType(item, t)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, reflect.ValueOf(val))
	}

	return arguments, nil
}
