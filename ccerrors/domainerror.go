package ccerrors

import (
	"fmt"
	"runtime"
)

type DomainError struct {
	code       ErrorCode
	detail     map[string]interface{}
	innerError error
	stackTrace []uintptr
}

type DomainErrorFunc func(*DomainError)

func NewDomainError(err error, code ErrorCode) *DomainError {
	return &DomainError{
		code:       code,
		innerError: err,
		stackTrace: callers(),
	}
}

func (e *DomainError) ErrorCode() ErrorCode {
	return e.code
}

func (e *DomainError) StackTrace() []string {
	result := make([]string, len(e.stackTrace))
	for i, pc := range e.stackTrace {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		frame := fmt.Sprintf("%s:%d %s", file, line, fn.Name())
		result[i] = frame
	}
	return result
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

func (e *DomainError) SetDetail(key string, value any) {
	if e.detail == nil {
		e.detail = make(map[string]interface{})
	}
	e.detail[key] = value
}

func (e *DomainError) Detail() map[string]interface{} {
	if e.innerError != nil {
		if e.detail == nil {
			e.detail = make(map[string]interface{})
		}

		if innerDomainErr, ok := e.innerError.(*DomainError); ok {
			for k, v := range innerDomainErr.Detail() {
				e.detail[k] = v
			}
		}
	}

	return e.detail
}

func (e *DomainError) Error() string {
	if e.innerError != nil {
		return fmt.Sprintf("%s: %v", e.code.Name(), e.innerError)

	}
	return e.code.Name()
}
