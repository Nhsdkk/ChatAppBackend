package validator_tests

import (
	"chat_app_backend/internal/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStructNotEmpty struct {
	stringVal string
	idVal     uuid.UUID `validator:"not_empty"`
	intVal    int       `validator:"not_empty"`
	intPtrVal *int      `validator:"not_empty"`
}

func TestValidator_NotEmpty_ShouldWorkWithFilledValues(t *testing.T) {
	intV := 123
	v := testStructNotEmpty{
		idVal:     uuid.New(),
		intVal:    1,
		intPtrVal: &intV,
	}

	validatorObject := validator.Validator[testStructNotEmpty]{}
	require.NoError(t, validatorObject.Validate(&v))
}

func TestValidator_NotEmpty_ShouldFailWithUnfilledValues(t *testing.T) {
	intV := 1
	v := testStructNotEmpty{
		intPtrVal: &intV,
	}

	validatorObject := validator.Validator[testStructNotEmpty]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v),
		`validation errors occurred:
field idVal is empty, but is required to be filled
field intVal is empty, but is required to be filled`,
	)
}
