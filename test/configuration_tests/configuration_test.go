package configuration_tests

import (
	configuration2 "chat_app_backend/internal/configuration"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStructConfiguration1 struct {
	v1 string
	v2 []int
}

type testStructConfiguration2 struct {
	v3 float32
}

type testStructConstructorResult1 struct {
	v1 string
	v2 []int
}

type testStructConstructorResult2 struct {
	v1 string
	v2 []int
	v3 float32
}

func TestConfiguration_ShouldReturnValueWithOneArgument(t *testing.T) {
	constructor := func(configuration1 *testStructConfiguration1) (*testStructConstructorResult1, error) {
		return &testStructConstructorResult1{
			v1: configuration1.v1,
			v2: configuration1.v2,
		}, nil
	}

	cfg := &testStructConfiguration1{
		v1: "v1",
		v2: []int{1, 2, 3},
	}

	configuration := configuration2.CreateConfiguration().
		AddConfiguration(cfg)

	constructedObject, err := configuration2.BuildFromConfiguration[testStructConstructorResult1](configuration, constructor)
	require.NoError(t, err)
	require.Equal(t, constructedObject.v1, cfg.v1)
	require.Equal(t, constructedObject.v2, cfg.v2)
}

func TestConfiguration_ShouldReturnErrorWhenConstructorReturnsError(t *testing.T) {
	constructor := func(configuration1 *testStructConfiguration1) (*testStructConstructorResult1, error) {
		return nil, errors.New("test error")
	}

	cfg := &testStructConfiguration1{
		v1: "v1",
		v2: []int{1, 2, 3},
	}

	configuration := configuration2.CreateConfiguration().
		AddConfiguration(cfg)

	constructedObject, err := configuration2.BuildFromConfiguration[testStructConstructorResult1](configuration, constructor)
	require.EqualError(t, err, "test error")
	require.Equal(t, constructedObject, (*testStructConstructorResult1)(nil))
}

func TestConfiguration_ShouldPanicWhenConfigurationExists(t *testing.T) {
	cfg := &testStructConfiguration1{
		v1: "v1",
		v2: []int{1, 2, 3},
	}

	cfg1 := &testStructConfiguration1{
		v1: "v2",
		v2: []int{1, 2, 3},
	}

	require.PanicsWithValue(
		t,
		"configuration with type testStructConfiguration1 already exists",
		func() {
			_ = configuration2.CreateConfiguration().
				AddConfiguration(cfg).
				AddConfiguration(cfg1)
		},
	)
}

func TestConfiguration_ShouldReturnValueWithMultipleArguments(t *testing.T) {
	constructor := func(configuration1 *testStructConfiguration1, configuration2 *testStructConfiguration2) (*testStructConstructorResult2, error) {
		return &testStructConstructorResult2{
			v1: configuration1.v1,
			v2: configuration1.v2,
			v3: configuration2.v3,
		}, nil
	}

	cfg := &testStructConfiguration1{
		v1: "v1",
		v2: []int{1, 2, 3},
	}

	cfg1 := &testStructConfiguration2{
		v3: 32.3,
	}

	configuration := configuration2.CreateConfiguration().
		AddConfiguration(cfg).
		AddConfiguration(cfg1)

	constructedObject, err := configuration2.BuildFromConfiguration[testStructConstructorResult2](configuration, constructor)
	require.NoError(t, err)
	require.Equal(t, constructedObject.v1, cfg.v1)
	require.Equal(t, constructedObject.v2, cfg.v2)
	require.Equal(t, constructedObject.v3, cfg1.v3)
}

func TestConfiguration_ShouldReturnValueWithAdditionalArguments(t *testing.T) {
	constructor := func(configuration1 *testStructConfiguration1, floatV *float32) (*testStructConstructorResult2, error) {
		return &testStructConstructorResult2{
			v1: configuration1.v1,
			v2: configuration1.v2,
			v3: *floatV,
		}, nil
	}

	cfg := &testStructConfiguration1{
		v1: "v1",
		v2: []int{1, 2, 3},
	}

	floatV := float32(32.3)

	configuration := configuration2.CreateConfiguration().
		AddConfiguration(cfg)

	constructedObject, err := configuration2.BuildFromConfiguration[testStructConstructorResult2](configuration, constructor, &floatV)
	require.NoError(t, err)
	require.Equal(t, constructedObject.v1, cfg.v1)
	require.Equal(t, constructedObject.v2, cfg.v2)
	require.Equal(t, constructedObject.v3, floatV)
}
