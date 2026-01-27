package ccvalidation

import "context"

type conditionalValidator[TCond any, TRequest any] struct {
	breakOnFailure   bool
	validators       map[any]Validator[TRequest]
	defaultValidator Validator[TRequest]
	condition        func(TRequest) TCond
}

func NewConditionalValidator[TCond any, TRequest any]() *conditionalValidator[TCond, TRequest] {
	return &conditionalValidator[TCond, TRequest]{
		validators: make(map[any]Validator[TRequest]),
	}
}

func (v *conditionalValidator[TCond, TRequest]) WithCondition(
	condition func(TRequest) TCond) *conditionalValidator[TCond, TRequest] {
	v.condition = condition
	return v
}

func (v *conditionalValidator[TCond, TRequest]) WithValidator(
	condition TCond,
	validator Validator[TRequest],
) *conditionalValidator[TCond, TRequest] {
	if v.breakOnFailure {
		val := NewValidator[TRequest]().BreakOnFailure()
		val.AddValidator(validator)
		validator = val
	}
	v.validators[condition] = validator
	return v
}

func (v *conditionalValidator[TCond, TRequest]) WithDefaultValidator(
	validator Validator[TRequest],
) *conditionalValidator[TCond, TRequest] {
	if v.breakOnFailure {
		val := NewValidator[TRequest]().BreakOnFailure()
		val.AddValidator(validator)
		validator = val
	}
	v.defaultValidator = validator
	return v
}

func (v *conditionalValidator[TCond, TRequest]) BreakOnFailure() *conditionalValidator[TCond, TRequest] {
	v.breakOnFailure = true
	return v
}

func (v *conditionalValidator[TCond, TRequest]) ValidateCtx(ctx context.Context, req TRequest) Result {
	condition := v.condition(req)
	validator, ok := v.validators[condition]
	if !ok {
		if v.defaultValidator != nil {
			return v.defaultValidator.ValidateCtx(ctx, req)
		}
		result := Result{}
		result.AddErrorMessage("no validator found for the given condition")
		return result
	}

	result := validator.ValidateCtx(ctx, req)
	return result
}
