package validator_tests

import (
	"chat_app_backend/internal/validator"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStructComparison struct {
	uintVal     uint     `validator:"gt 5"`
	intVal      int      `validator:"gte 0"`
	floatVal    float32  `validator:"lt -5"`
	float64Val  float64  `validator:"lte -10.5"`
	uint8Val    uint8    `validator:"eq 5"`
	floatPtrVal *float64 `validator:"lt -5"`
}

func TestValidator_Comparison_ShouldWorkWithPassingCondition(t *testing.T) {
	floatV := -15.
	v := testStructComparison{
		uintVal:     6,
		intVal:      2,
		floatVal:    -6,
		float64Val:  -10.5,
		uint8Val:    5,
		floatPtrVal: &floatV,
	}

	validatorObject := validator.Validator[testStructComparison]{}
	require.NoError(t, validatorObject.Validate(&v))
}

func TestValidator_Comparison_ShouldFailWithNotPassingCondition(t *testing.T) {
	v := testStructComparison{
		uintVal:    6,
		intVal:     -1,
		floatVal:   -6,
		float64Val: -11.5,
		uint8Val:   4,
	}

	validatorObject := validator.Validator[testStructComparison]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v),
		`validation errors occurred:
the value in field intVal should be gte than 0 but it is not
the value in field uint8Val should be eq than 5 but it is not`,
	)
}
