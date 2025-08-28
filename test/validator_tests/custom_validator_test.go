package validator_tests

import (
	"chat_app_backend/internal/validator"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStructCustomValidators struct {
	id    int
	name  string
	email string
}

type IdValidator struct {
	RequiredId int
}

func (i IdValidator) Validate(id *int) bool {
	if *id != i.RequiredId {
		return false
	}

	return true
}

type NameValidator struct{}

func (n NameValidator) Validate(name *string) bool {
	if len(*name) < 10 {
		return false
	}

	return true
}

type EmailValidator struct{}

func (e EmailValidator) Validate(email *string) bool {
	if !strings.Contains(*email, "@") {
		return false
	}
	return true
}

func TestValidator_CustomValidators_ShouldWorkWithRightValue(t *testing.T) {
	v := testStructCustomValidators{
		id:    1,
		name:  "wow_long_username",
		email: "email@gmail.com",
	}
	validatorObject := validator.Validator[testStructCustomValidators]{}
	require.NoError(
		t,
		validatorObject.
			AttachValidator(
				validator.ExternalValidator[testStructCustomValidators, string]{}.
					RuleFor(
						func(data *testStructCustomValidators) *string {
							return &data.name
						},
					).
					Must(NameValidator{}).
					WithMessage("name is too short").
					Validate,
			).
			AttachValidator(
				validator.ExternalValidator[testStructCustomValidators, string]{}.
					RuleFor(
						func(data *testStructCustomValidators) *string {
							return &data.email
						},
					).
					Must(EmailValidator{}).
					WithMessage("email has wrong format").
					Validate,
			).
			AttachValidator(
				validator.ExternalValidator[testStructCustomValidators, int]{}.
					RuleFor(
						func(data *testStructCustomValidators) *int {
							return &data.id
						},
					).
					Must(IdValidator{RequiredId: 1}).
					WithMessage("id does not match").
					Validate,
			).
			Validate(&v),
	)
}

func TestValidator_CustomValidators_ShouldFailWithWrongValue(t *testing.T) {
	v := testStructCustomValidators{
		id:    2,
		name:  "usr",
		email: "email.com",
	}
	validatorObject := validator.Validator[testStructCustomValidators]{}
	require.EqualError(
		t,
		validatorObject.
			AttachValidator(
				validator.ExternalValidator[testStructCustomValidators, string]{}.
					RuleFor(
						func(data *testStructCustomValidators) *string {
							return &data.name
						},
					).
					Must(NameValidator{}).
					WithMessage("name is too short").
					Validate,
			).
			AttachValidator(
				validator.ExternalValidator[testStructCustomValidators, string]{}.
					RuleFor(
						func(data *testStructCustomValidators) *string {
							return &data.email
						},
					).
					Must(EmailValidator{}).
					WithMessage("email has wrong format").
					Validate,
			).
			AttachValidator(
				validator.ExternalValidator[testStructCustomValidators, int]{}.
					RuleFor(
						func(data *testStructCustomValidators) *int {
							return &data.id
						},
					).
					Must(IdValidator{RequiredId: 1}).
					WithMessage("id does not match").
					Validate,
			).
			Validate(&v),
		`validation errors occurred:
name is too short
email has wrong format
id does not match`,
	)
}
