package validator_tests

import (
	"chat_app_backend/internal/validator"
	"errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type testStructCustomValidators struct {
	name  string
	email string
}

func nameValidator(v *testStructCustomValidators) error {
	if len(v.name) < 10 {
		return errors.New("name is too short")
	}
	return nil
}

func emailValidator(v *testStructCustomValidators) error {
	if !strings.Contains(v.email, "@") {
		return errors.New("email has wrong format")
	}
	return nil
}

func TestValidator_CustomValidators_ShouldWorkWithRightValue(t *testing.T) {
	v := testStructCustomValidators{
		name:  "wow_long_username",
		email: "email@gmail.com",
	}
	validatorObject := validator.Validator[testStructCustomValidators]{}
	require.NoError(t, validatorObject.AttachValidator(nameValidator).AttachValidator(emailValidator).Validate(&v))
}

func TestValidator_CustomValidators_ShouldFailWithWrongValue(t *testing.T) {
	v := testStructCustomValidators{
		name:  "usr",
		email: "email.com",
	}
	validatorObject := validator.Validator[testStructCustomValidators]{}
	require.EqualError(
		t,
		validatorObject.AttachValidator(nameValidator).AttachValidator(emailValidator).Validate(&v),
		`validation errors occurred:
name is too short
email has wrong format`,
	)
}
