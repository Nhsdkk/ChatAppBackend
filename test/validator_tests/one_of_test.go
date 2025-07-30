package validator_tests

import (
	"chat_app_backend/internal/validator"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStructOneOf struct {
	stringValueEnum string `validator:"one_of [string1,string2]"`
	intValueEnum    int    `validator:"one_of [1,15]"`
}

func TestValidator_OneOf_ShouldWorkWhenValueIsInRange(t *testing.T) {
	v := testStructOneOf{
		stringValueEnum: "string1",
		intValueEnum:    15,
	}

	validatorObject := validator.Validator[testStructOneOf]{}
	require.NoError(t, validatorObject.Validate(&v))
}

func TestValidator_OneOf_ShouldFailWhenValueIsNotInRange(t *testing.T) {
	v := testStructOneOf{
		stringValueEnum: "string3",
		intValueEnum:    15,
	}

	validatorObject := validator.Validator[testStructOneOf]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v),
		`validation errors occurred:
field stringValueEnum is required to have one of these values [string1 string2]`,
	)
}
