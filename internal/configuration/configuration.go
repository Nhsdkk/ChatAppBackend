package configuration

import (
	"errors"
	"fmt"
	"reflect"
)

type IConfiguration interface {
	Get(vType interface{}) (interface{}, error)
	AddConfiguration(value interface{}) IConfiguration
}

type Configuration struct {
	configurations map[string]interface{}
}

func (c *Configuration) Get(vType interface{}) (interface{}, error) {
	t := reflect.TypeOf(vType)

	if t == nil {
		panic("can't accept interface type")
	}

	if t.Kind() != reflect.Pointer {
		panic("can't get non pointer type")
	}

	key := ""
	if t.Elem() == nil {
		key = t.Name()
	} else {
		key = t.Elem().Name()
	}

	value, exists := c.configurations[key]
	if !exists {
		return nil, errors.New(fmt.Sprintf("can't find configuration of type %s", t))
	}

	return value, nil
}

func (c *Configuration) AddConfiguration(value interface{}) IConfiguration {
	t := reflect.TypeOf(value)

	if t == nil {
		panic("can't accept interface type")
	}

	if t.Kind() != reflect.Pointer {
		panic("configuration is not of pointer type")
	}

	key := ""
	if t.Elem() == nil {
		key = t.Name()
	} else {
		key = t.Elem().Name()
	}

	if _, exists := c.configurations[key]; exists {
		panic(fmt.Sprintf("configuration with type %s already exists", key))
	}

	c.configurations[key] = value
	return c
}

func validateConstructorFunction[T interface{}](constructorType reflect.Type) error {
	if constructorType.IsVariadic() {
		return errors.New("variadic functions are unsupported")
	}

	if constructorType.NumOut() != 2 ||
		!constructorType.Out(1).Implements(reflect.TypeFor[error]()) ||
		constructorType.Out(0).Kind() != reflect.Pointer ||
		!constructorType.Out(0).Elem().AssignableTo(reflect.TypeFor[T]()) {
		return errors.New(fmt.Sprintf("can't construct object of type %s using this constructor as it does not follow this signatore (func[T interface{}](...) (*T, error))", reflect.TypeFor[T]()))
	}

	for i := range constructorType.NumIn() {
		argT := constructorType.In(i)
		if argT.Kind() != reflect.Pointer {
			return errors.New(fmt.Sprintf("can't use this constructor as there is an argument of type %s, which is not pointer", argT.Name()))
		}
	}

	return nil
}

func BuildFromConfiguration[T interface{}](configuration IConfiguration, constructor interface{}, additionalArguments ...interface{}) (*T, error) {
	constructorType := reflect.TypeOf(constructor)

	if constructorType.Kind() != reflect.Func {
		return nil, errors.New("can't construct as constructor is not a function")
	}

	if constructorValidationError := validateConstructorFunction[T](constructorType); constructorValidationError != nil {
		return nil, constructorValidationError
	}

	additionalArgumentValues := make(map[string]reflect.Value)
	for _, item := range additionalArguments {
		t := reflect.TypeOf(item)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		if _, exists := additionalArgumentValues[t.Name()]; exists {
			panic("two arguments of the same type passed")
		}

		additionalArgumentValues[t.Name()] = reflect.ValueOf(item)
	}

	inp := make([]reflect.Value, 0)
	for idx := range constructorType.NumIn() {
		argType := constructorType.In(idx)
		val := reflect.New(argType.Elem())

		arg, err := configuration.Get(val.Interface())
		if err != nil {
			additionalArgVal, exists := additionalArgumentValues[argType.Elem().Name()]
			if !exists {
				return nil, errors.New(fmt.Sprintf("can't find value for argument of type %s", argType.Elem().Name()))
			}
			inp = append(inp, additionalArgVal)
			continue
		}

		inp = append(inp, reflect.ValueOf(arg))
	}

	out := reflect.ValueOf(constructor).Call(inp)

	v := out[0].Interface().(*T)

	err, ok := out[1].Interface().(error)
	if !ok {
		return v, nil
	}

	return v, err
}

func CreateConfiguration() IConfiguration {
	return &Configuration{
		configurations: make(map[string]interface{}),
	}
}
