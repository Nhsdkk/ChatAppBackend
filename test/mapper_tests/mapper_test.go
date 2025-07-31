package mapper_tests

import (
	"chat_app_backend/internal/mapper"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMapper_SingleArgument_ShouldMapValuesWhenPassedStructsWithEnoughArguments(t *testing.T) {
	v1 := struct {
		V1 int
		V2 []int
		V3 []struct {
			V1 string
		}
	}{
		V1: 0,
		V2: []int{1, 3, 2},
		V3: []struct {
			V1 string
		}{
			{V1: "qwe"},
			{V1: "qwe1"},
		},
	}

	var v2 struct {
		V2 []int
		V3 []struct {
			V1 string
		}
	}

	err := mapper.Mapper{}.Map(&v2, v1)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(v2.V2, v1.V2))
	require.True(t, reflect.DeepEqual(v2.V3, v1.V3))
}

func TestMapper_MultipleArgument_ShouldMapValuesWhenPassedStructsWithEnoughArguments(t *testing.T) {
	v1 := struct {
		V1 int
		V2 []int
		V3 []struct {
			V1 string
		}
	}{
		V1: 0,
		V2: []int{1, 3, 2},
		V3: []struct {
			V1 string
		}{
			{V1: "qwe"},
			{V1: "qwe1"},
		},
	}

	v3 := struct {
		V5 []struct {
			intV []int
		}
	}{
		[]struct{ intV []int }{
			{intV: []int{1, 3, 2}},
			{intV: []int{1, 4, 2}},
		},
	}

	var v2 struct {
		V2 []int
		V3 []struct {
			V1 string
		}
		V5 []struct {
			intV []int
		}
	}

	err := mapper.Mapper{}.Map(&v2, v1, v3)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(v2.V2, v1.V2))
	require.True(t, reflect.DeepEqual(v2.V3, v1.V3))
	require.True(t, reflect.DeepEqual(v2.V5, v3.V5))
}

func TestMapper_MultipleArgument_ShouldFailWhenNameCollisionIsDetected(t *testing.T) {
	v1 := struct {
		V1 int
		V2 []int
		V3 []struct {
			V1 string
		}
	}{
		V1: 0,
		V2: []int{1, 3, 2},
		V3: []struct {
			V1 string
		}{
			{V1: "qwe"},
			{V1: "qwe1"},
		},
	}

	v3 := struct {
		V3 []struct {
			intV []int
		}
	}{
		[]struct{ intV []int }{
			{intV: []int{1, 3, 2}},
			{intV: []int{1, 4, 2}},
		},
	}

	var v2 struct {
		V2 []int
		V3 []struct {
			V1 string
		}
	}

	err := mapper.Mapper{}.Map(&v2, v1, v3)
	require.EqualError(
		t,
		err,
		`name collision found as field V3 exists in more than one struct`,
	)
}

func TestMapper_SingleArgument_ShouldFailWhenSourceDoesNotHaveEnoughArguments(t *testing.T) {
	v1 := struct {
		V1 int
		V2 []int
		V3 []struct {
			V1 string
		}
	}{
		V1: 0,
		V2: []int{1, 3, 2},
		V3: []struct {
			V1 string
		}{
			{V1: "qwe"},
			{V1: "qwe1"},
		},
	}

	var v2 struct {
		V2 []int
		V5 []struct {
			V1 string
		}
	}

	err := mapper.Mapper{}.Map(&v2, v1)
	require.EqualError(
		t,
		err,
		`src does not have field V5, which dest has`,
	)
}
