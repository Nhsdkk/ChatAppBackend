package validator

import (
	"context"
	"errors"
)

type IValidation[T1, T2 any] interface {
	Must(validation IInternalValidation[T2]) IValidation[T1, T2]
	WithMessage(message string) IValidation[T1, T2]
	Optional() IValidation[T1, T2]
	Validate(data *T1, ctx context.Context) error
}

type IExternalValidator[T1, T2 any] interface {
	RuleFor(func(data T1) T2) IValidation[T1, T2]
}

type IInternalValidation[T any] interface {
	Validate(data *T, ctx context.Context) bool
}

type ExternalValidator[T1 any, T2 any] struct {
}

func (e ExternalValidator[T1, T2]) RuleFor(f func(data *T1) *T2) IValidation[T1, T2] {
	return Validation[T1, T2]{
		transformation: f,
	}
}

type Validation[T1, T2 any] struct {
	transformation func(data *T1) *T2
	validation     IInternalValidation[T2]
	optional       bool
	message        string
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

func (v Validation[T1, T2]) Validate(data *T1, ctx context.Context) error {
	switch {
	case v.optional && v.transformation(data) == nil:
		return nil
	case v.transformation(data) == nil, !v.validation.Validate(v.transformation(data), ctx):
		return errors.New(v.message)
	default:
		return nil
	}
}
