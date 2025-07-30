package validator_tests

import (
	"chat_app_backend/internal/validator"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStructMultipleValidations struct {
	stringEnumValue string `validator:"not_empty;one_of [string1,string2]"`
	intValue        int    `validator:"not_empty;gt 5;lt 10"`
}

func TestValidator_MultipleValidations_ShouldWorkWhenPassedRightValues(t *testing.T) {
	v := testStructMultipleValidations{
		stringEnumValue: "string1",
		intValue:        7,
	}
	validatorObject := validator.Validator[testStructMultipleValidations]{}

	require.NoError(t, validatorObject.Validate(&v))
}

func TestValidator_MultipleValidations_ShouldFailWhenPassedWrongValues(t *testing.T) {
	v := testStructMultipleValidations{
		stringEnumValue: "string2",
		intValue:        0,
	}

	validatorObject := validator.Validator[testStructMultipleValidations]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v),
		`validation errors occurred:
field intValue is empty, but is required to be filled
the value in field intValue should be gt than 5 but it is not`,
	)
}
