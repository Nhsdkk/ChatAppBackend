package validator_tests

import (
	"chat_app_backend/internal/validator"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testStructNestedArray struct {
	stringVal string
	idVal     uuid.UUID `validator:"not_empty"`
	intVal    int       `validator:"not_empty"`
	arrVal    []struct {
		stringVal string
		idVal     uuid.UUID `validator:"not_empty"`
		intVal    int       `validator:"not_empty"`
	}
}

type testStructNested struct {
	stringVal string
	idVal     uuid.UUID `validator:"not_empty"`
	intVal    int       `validator:"not_empty"`
	arrVal    struct {
		stringVal string
		idVal     uuid.UUID `validator:"not_empty"`
		intVal    int       `validator:"not_empty"`
	}
}

func TestValidator_NestedStructs_ShouldWorkWithArrayOfStructsWithFilledValues(t *testing.T) {
	v := testStructNestedArray{
		idVal:  uuid.New(),
		intVal: 1,
		arrVal: []struct {
			stringVal string
			idVal     uuid.UUID `validator:"not_empty"`
			intVal    int       `validator:"not_empty"`
		}{
			{
				idVal:  uuid.New(),
				intVal: 1,
			},
		},
	}

	validatorObject := validator.Validator[testStructNestedArray]{}
	require.NoError(t, validatorObject.Validate(&v, context.Background()))
}

func TestValidator_NestedStructs_ShouldFailWithArrayOfStructsWithUnfilledValues(t *testing.T) {
	v := testStructNestedArray{
		intVal: 1,
		arrVal: []struct {
			stringVal string
			idVal     uuid.UUID `validator:"not_empty"`
			intVal    int       `validator:"not_empty"`
		}{
			{},
		},
	}

	validatorObject := validator.Validator[testStructNestedArray]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v, context.Background()),
		`validation errors occurred:
field idVal is empty, but is required to be filled
field idVal is empty, but is required to be filled
field intVal is empty, but is required to be filled`,
	)
}

func TestValidator_NestedStructs_ShouldWorkWithStructsWithFilledValues(t *testing.T) {
	v := testStructNested{
		idVal:  uuid.New(),
		intVal: 1,
		arrVal: struct {
			stringVal string
			idVal     uuid.UUID `validator:"not_empty"`
			intVal    int       `validator:"not_empty"`
		}{
			idVal:  uuid.New(),
			intVal: 1,
		},
	}

	validatorObject := validator.Validator[testStructNested]{}
	require.NoError(t, validatorObject.Validate(&v, context.Background()))
}

func TestValidator_NestedStructs_ShouldFailWithStructsWithUnfilledValues(t *testing.T) {
	v := testStructNested{
		intVal: 1,
		arrVal: struct {
			stringVal string
			idVal     uuid.UUID `validator:"not_empty"`
			intVal    int       `validator:"not_empty"`
		}{},
	}

	validatorObject := validator.Validator[testStructNested]{}
	require.EqualError(
		t,
		validatorObject.Validate(&v, context.Background()),
		`validation errors occurred:
field idVal is empty, but is required to be filled
field idVal is empty, but is required to be filled
field intVal is empty, but is required to be filled`,
	)
}
