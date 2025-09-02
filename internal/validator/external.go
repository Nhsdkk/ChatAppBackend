package validator

import (
	"chat_app_backend/internal/request_env"
	"context"
	"errors"
	"fmt"
	"strings"
)

type ValidationFunction[T any] = func(data *T, ctx context.Context, env request_env.RequestEnv) error

type ExternalValidatorGroup[T any] struct {
	validations      []ValidationFunction[T]
	exceptionFactory ExceptionFactory
}

func (vGroup ExternalValidatorGroup[T]) WithExceptionFactory(factory ExceptionFactory) ExternalValidatorGroup[T] {
	vGroup.exceptionFactory = factory
	return vGroup
}

func (vGroup ExternalValidatorGroup[T]) AttachValidation(function ValidationFunction[T]) ExternalValidatorGroup[T] {
	vGroup.validations = append(vGroup.validations, function)
	return vGroup
}

func (vGroup ExternalValidatorGroup[T]) Validate(data *T, ctx context.Context, env request_env.RequestEnv) error {
	msgs := make([]string, 0)

	for _, validation := range vGroup.validations {
		if err := validation(data, ctx, env); err != nil {
			msgs = append(msgs, err.Error())
		}
	}

	if len(msgs) == 0 {
		return nil
	}

	return vGroup.exceptionFactory(fmt.Sprintf("validation errors occurred:\n%s", strings.Join(msgs, "\n")))
}

func CreateValidatorGroup[T any]() ExternalValidatorGroup[T] {
	return ExternalValidatorGroup[T]{
		exceptionFactory: func(message string) error {
			return errors.New(message)
		},
	}
}

type ExceptionFactory = func(message string) error

type IValidation[T1, T2 any] interface {
	Must(validation IInternalValidation[T2]) IValidation[T1, T2]
	WithMessage(message string) IValidation[T1, T2]
	WithExceptionFactory(exceptionFactory ExceptionFactory) IValidation[T1, T2]
	Optional() IValidation[T1, T2]
	Validate(data *T1, ctx context.Context, env request_env.RequestEnv) error
}

type IExternalValidator[T1, T2 any] interface {
	RuleFor(func(data T1) T2) IValidation[T1, T2]
}

type IInternalValidation[T any] interface {
	Validate(data *T, ctx context.Context, env request_env.RequestEnv) bool
}

type ExternalValidator[T1 any, T2 any] struct {
}

func (e ExternalValidator[T1, T2]) RuleFor(f func(data *T1) *T2) IValidation[T1, T2] {
	return Validation[T1, T2]{
		transformation: f,
		exceptionFactory: func(message string) error {
			return errors.New(message)
		},
	}
}

type Validation[T1, T2 any] struct {
	transformation   func(data *T1) *T2
	validation       IInternalValidation[T2]
	exceptionFactory ExceptionFactory
	optional         bool
	message          string
}

func (v Validation[T1, T2]) WithExceptionFactory(exceptionFactory ExceptionFactory) IValidation[T1, T2] {
	v.exceptionFactory = exceptionFactory
	return v
}

func (v Validation[T1, T2]) Optional() IValidation[T1, T2] {
	v.optional = true
	return v
}

func (v Validation[T1, T2]) Must(validation IInternalValidation[T2]) IValidation[T1, T2] {
	v.validation = validation
	return v
}

func (v Validation[T1, T2]) WithMessage(message string) IValidation[T1, T2] {
	v.message = message
	return v
}

func (v Validation[T1, T2]) Validate(data *T1, ctx context.Context, env request_env.RequestEnv) error {
	switch {
	case v.optional && v.transformation(data) == nil:
		return nil
	case v.transformation(data) == nil, !v.validation.Validate(v.transformation(data), ctx, env):
		return v.exceptionFactory(v.message)
	default:
		return nil
	}
}
