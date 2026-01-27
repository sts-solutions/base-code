package ccvalidation

import "context"

type validationStep[T any] struct {
	breakOnFailure bool
	action         func(ctx context.Context, req T) error
}

func (v *validationStep[T]) BreakOnFailure() *validationStep[T] {
	v.breakOnFailure = true
	return v
}
