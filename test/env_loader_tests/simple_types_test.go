package env_loader_tests

import (
	"chat_app_backend/internal/env_loader"
	"github.com/stretchr/testify/require"
	"os"
	path "path/filepath"
	"testing"
)

type testStructEnvSimpleTypes struct {
	IntVal    int     `env:"int_val"`
	FloatVal  float32 `env:"float_val"`
	StringVal string  `env:"string_val"`
}

func TestEnvLoader_SimpleTypes_ShouldWorkWhenUsingEnvFileWithAbsolutePath(t *testing.T) {
	var v testStructEnvSimpleTypes
	currentFolderAbsPath, _ := path.Abs(".")

	loader, err := env_loader.CreateLoaderFromFile(path.Join(currentFolderAbsPath, "./test_data/.env"))
	require.NoError(t, err)

	require.NoError(t, loader.LoadDataIntoStruct(&v))

	require.Equal(t, v.StringVal, "qweqweqweqw")
	require.Equal(t, v.IntVal, 1)
	require.Equal(t, v.FloatVal, float32(15.3))
}

func TestEnvLoader_SimpleTypes_ShouldWorkWhenUsingEnvFileWithRelativePath(t *testing.T) {
	var v testStructEnvSimpleTypes

	loader, err := env_loader.CreateLoaderFromFile("./test_data/.env")
	require.NoError(t, err)

	require.NoError(t, loader.LoadDataIntoStruct(&v))

	require.Equal(t, v.StringVal, "qweqweqweqw")
	require.Equal(t, v.IntVal, 1)
	require.Equal(t, v.FloatVal, float32(15.3))
}

func TestEnvLoader_SimpleTypes_ShouldWorkWhenUsingEnvironment(t *testing.T) {
	env := map[string]string{
		"testStructEnvSimpleTypes_int_val":    "1",
		"testStructEnvSimpleTypes_float_val":  "15.3",
		"testStructEnvSimpleTypes_string_val": "qweqweqweqw",
	}

	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

	var v testStructEnvSimpleTypes

	loader := env_loader.CreateLoaderFromEnv()

	require.NoError(t, loader.LoadDataIntoStruct(&v))

	require.Equal(t, v.StringVal, "qweqweqweqw")
	require.Equal(t, v.IntVal, 1)
	require.Equal(t, v.FloatVal, float32(15.3))
}
