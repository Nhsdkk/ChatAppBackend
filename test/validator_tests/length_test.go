package validator_tests

import (
	"chat_app_backend/internal/validator"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStructLength struct {
	StringVal string   `validator:"length eq 3"`
	SliceVal  []string `validator:"length lt 5"`
}

func TestValidator_Length_ShouldWorkWhenPassedCorrectValues(t *testing.T) {
	v := testStructLength{
		StringVal: "qwe",
		SliceVal:  []string{"123", "12412412", "12313123"},
	}

	validatorObject := validator.Validator[testStructLength]{}
	require.NoError(t, validatorObject.Validate(&v, context.Background()))
}

func TestValidator_Length_ShouldWorkWhenPassedWrongValues(t *testing.T) {
	v := testStructLength{
		StringVal: "qweq",
		SliceVal:  []string{"123", "12412412", "12313123", "123", "12412412", "12313123"},
	}

	validatorObject := validator.Validator[testStructLength]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v, context.Background()),
		`validation errors occurred:
length of value under field StringVal is not eq than 3
length of value under field SliceVal is not lt than 5`,
	)
}
