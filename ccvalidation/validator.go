package ccvalidation

import "context"

type Validator[T any] interface {
	Validate(src T) Result
	ValidateCtx(ctx context.Context, src T) Result
}

type validator[T any] struct {
	breakOnFailure bool
	steps          []validationStep[T]
}

func (v validator[T]) Validate(src T) Result {
	return v.ValidateCtx(context.Background(), src)
}

func (v validator[T]) ValidateCtx(ctx context.Context, src T) Result {
	result := Result{}

	if !v.hasValidationItems() {
		result.AddFailureMessage("No validation steps defined")
		return result
	}

	for _, step := range v.steps {
		err := step.action(ctx, src)
		if err != nil {
			if res, ok := err.(Result); ok {
				if res.IsSuccess() {
					continue
				}
				for _, v := range res.GetErrors() {
					result.AddFailure(v)
				}
			} else {
				result.AddFailure(err)
			}
			if step.breakOnFailure || v.breakOnFailure {
				return result
			}
		}
	}

	return result
}

func NewValidator[T any]() *validator[T] {
	return &validator[T]{
		breakOnFailure: false,
		steps:          make([]validationStep[T], 0),
	}
}

func New[T any]() *validator[T] {
	return NewValidator[T]()
}

func (v *validator[T]) BreakOnFailure() *validator[T] {
	v.breakOnFailure = true
	return v
}

func (v *validator[T]) AddStep(steps ...func(req T) error) {
	if steps == nil {
		steps = []func(req T) error{func(T) error { return nil }}
	}

	for _, step := range steps {
		ctxFunc := func(ctx context.Context, req T) error {
			return step(req)
		}
		v.AddStepCtx(ctxFunc)
	}
}

func (v *validator[T]) AddStepCtx(steps ...func(ctx context.Context, req T) error) {
	if steps == nil {
		steps = []func(ctx context.Context, req T) error{
			func(ctx context.Context, req T) error {
				return nil
			},
		}
	}

	for _, step := range steps {
		validationStep := validationStep[T]{
			breakOnFailure: false,
			action:         step,
		}
		v.steps = append(v.steps, validationStep)
	}
}

func (v *validator[T]) AddValidator(vldtr Validator[T]) {
	if v.breakOnFailure {
		val, _ := vldtr.(*validator[T])
		val.breakOnFailure = true
		val.steps = append(val.steps, v.steps...)
		vldtr = val
	}

	ctxFunc := func(ctx context.Context, req T) error {
		return vldtr.ValidateCtx(ctx, req)
	}
	v.AddStepCtx(ctxFunc)

}

func (v *validator[T]) hasValidationItems() bool {
	return len(v.steps) > 0
}
