package env_loader_tests

import (
	"chat_app_backend/internal/env_loader"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testStructEnvArrayTypes struct {
	IntArray    []int     `env:"int_array"`
	FloatArray  []float64 `env:"float_array"`
	StringArray []string  `env:"string_array"`
}

func TestEnvLoader_ArrayTypes_ShouldWorkWhenPassedArrayFromEnvironment(t *testing.T) {
	env := map[string]string{
		"testStructEnvArrayTypes_int_array":    "[1, 2, 3]",
		"testStructEnvArrayTypes_float_array":  "[15.3, 14.5]",
		"testStructEnvArrayTypes_string_array": "[qwe,qqq]",
	}

	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

	var v testStructEnvArrayTypes

	envLoader := env_loader.CreateLoaderFromEnv()

	require.NoError(t, envLoader.LoadDataIntoStruct(&v))

	require.Equal(t, v.StringArray, []string{"qwe", "qqq"})
	require.Equal(t, v.FloatArray, []float64{15.3, 14.5})
	require.Equal(t, v.IntArray, []int{1, 2, 3})
}
