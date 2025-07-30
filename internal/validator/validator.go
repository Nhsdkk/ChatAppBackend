package validator

import (
	validator "chat_app_backend/internal/validator/utils"
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationFunction[T interface{}] = func(data *T) error

type IValidator[T interface{}] interface {
	Validate(value *T) error
	AttachValidator(function ValidationFunction[T]) IValidator[T]
}

type Validator[T interface{}] struct {
	additionalValidations []ValidationFunction[T]
}

func (v Validator[T]) AttachValidator(function ValidationFunction[T]) IValidator[T] {
	v.additionalValidations = append(v.additionalValidations, function)
	return v
}

func (v Validator[T]) Validate(value *T) error {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return errors.New("can't parse tags on non pointer struct type")
	}

	errs := validate(reflect.ValueOf(value).Elem())

	for _, validation := range v.additionalValidations {
		if err := validation(value); err != nil {
			errs = append(errs, err)
		}
	}

	msgs := make([]string, len(errs))

	for idx, err := range errs {
		msgs[idx] = err.Error()
	}

	if len(errs) != 0 {
		return errors.New(fmt.Sprintf("validation errors occurred:\n%s", strings.Join(msgs, "\n")))
	}
	return nil
}

func validate(value reflect.Value) []error {
	t := value.Type()

	errs := make([]error, 0)

	for fieldIdx := range t.NumField() {
		fieldValue := value.Field(fieldIdx)
		fieldType := t.Field(fieldIdx)

		tag, tagParsingError := validator.ParseTag(fieldType)
		if tagParsingError != nil {
			panic(tagParsingError)
		}

		if tag != nil {
			newErrs := handleValidation(fieldType, fieldValue, tag)
			errs = append(errs, newErrs...)
		}

		switch fieldType.Type.Kind() {
		case reflect.Array, reflect.Slice:
			if fieldType.Type.Elem().Kind() != reflect.Struct || fieldValue.Len() == 0 {
				continue
			}

			for arrayIdx := range fieldValue.Len() {
				newErrs := validate(fieldValue.Index(arrayIdx))
				errs = append(errs, newErrs...)
			}

			break
		case reflect.Struct:
			newErrs := validate(fieldValue)
			errs = append(errs, newErrs...)
			break
		default:
			break
		}
	}

	return errs
}

func handleValidation(fieldType reflect.StructField, value reflect.Value, tag *validator.Tag) []error {
	errs := make([]error, 0)
	for _, validation := range tag.Validations {
		switch validation.ValidationType {
		case validator.NotEmpty:
			if err := handleNotEmpty(value, fieldType.Name); err != nil {
				errs = append(errs, err)
			}
			break
		case validator.OneOf:
			if fieldType.Type.Kind() == reflect.Pointer || fieldType.Type.Kind() == reflect.UnsafePointer || fieldType.Type.Kind() == reflect.Uintptr {
				if value.IsNil() {
					break
				}
				value = value.Elem()
			}
			if err := handleOneOf(value, fieldType.Name, validation.Arguments); err != nil {
				errs = append(errs, err)
			}
			break
		case validator.Greater, validator.GreaterOrEqual, validator.Equal, validator.LessOrEqual, validator.Less:
			if fieldType.Type.Kind() == reflect.Pointer || fieldType.Type.Kind() == reflect.UnsafePointer || fieldType.Type.Kind() == reflect.Uintptr {
				if value.IsNil() {
					break
				}
				value = value.Elem()
			}
			if err := handleComparison(value, fieldType, validation.ValidationType, validation.Arguments[0]); err != nil {
				errs = append(errs, err)
			}
			break
		case validator.Length:
			if err := handleLength(value, fieldType.Name, validation.Arguments); err != nil {
				errs = append(errs, err)
			}
			break
		}
	}

	return errs
}

func handleComparison(value reflect.Value, fieldType reflect.StructField, validationType validator.ValidationTypes, argument reflect.Value) error {
	if !value.Comparable() {
		panic(fmt.Sprintf("can't apply comparison validation for field %s as its not comparable", fieldType.Name))
	}

	var cmpFunc func() bool

	switch value.Kind() {
	case reflect.String:
		cmpFunc = func() bool {
			v1, v2 := value.String(), argument.String()
			return compare(v1, v2, validationType)
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cmpFunc = func() bool {
			v1, v2 := value.Int(), argument.Int()
			return compare(v1, v2, validationType)
		}
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		cmpFunc = func() bool {
			v1, v2 := value.Uint(), argument.Uint()
			return compare(v1, v2, validationType)
		}
		break
	case reflect.Float32, reflect.Float64:
		cmpFunc = func() bool {
			v1, v2 := value.Float(), argument.Float()
			return compare(v1, v2, validationType)
		}
		break
	default:
		panic(fmt.Sprintf("can't apply comparison validation for field %s as its not comparable", fieldType.Name))
	}

	if !cmpFunc() {
		return errors.New(fmt.Sprintf("the value in field %v should be %v than %v but it is not", fieldType.Name, validationType, argument))
	}

	return nil
}

func handleNotEmpty(value reflect.Value, fieldName string) error {
	if value.IsZero() {
		return errors.New(fmt.Sprintf("field %s is empty, but is required to be filled", fieldName))
	}

	return nil
}

func handleOneOf(value reflect.Value, fieldName string, arguments []reflect.Value) error {
	for _, argument := range arguments {
		if value.Equal(argument) {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("field %s is required to have one of these values %v", fieldName, arguments))
}

func handleLength(value reflect.Value, fieldName string, arguments []reflect.Value) error {
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		break
	default:
		panic(fmt.Sprintf("can't apply length validator to field %s of type %s", fieldName, value.Kind()))
	}

	if !compare(
		value.Len(),
		int(arguments[1].Int()),
		validator.ValidationTypes(arguments[0].String()),
	) {
		return errors.New(fmt.Sprintf("length of value under field %s is not %s than %d", fieldName, arguments[0].Convert(reflect.TypeFor[string]()).String(), arguments[1].Int()))
	}

	return nil
}

func compare[T cmp.Ordered](v1 T, v2 T, cmpType validator.ValidationTypes) bool {
	switch cmpType {
	case validator.Greater:
		return v1 > v2
	case validator.GreaterOrEqual:
		return v1 >= v2
	case validator.Less:
		return v1 < v2
	case validator.LessOrEqual:
		return v1 <= v2
	case validator.Equal:
		return v1 == v2
	default:
		panic("attempt to use compare function with non comparable type of operation")
	}
}
