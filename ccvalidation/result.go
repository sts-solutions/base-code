package ccvalidation

import (
	"errors"
	"fmt"
	"strings"
)

// Result holds validation results and errors
type Result struct {
	failures []error
}

// Error implements the error interface for Result
// It returns a concatenated string of all error messages if there are any errors
// If there are no errors, it returns an empty string
func (e Result) Error() string {
	if e.IsSuccess() {
		return ""
	}
	return strings.Join(e.GetErrorMessages(), ";")
}

// Is checks if the given error is of type Result
func (e Result) Is(err error) bool {
	_, ok := err.(Result)
	return ok
}

// AddFailureMessage adds a validation failure message to the Result
// If the message is empty or whitespace, nothing is added
func (r *Result) AddFailureMessage(msg string) {
	if strings.TrimSpace(msg) == "" {
		return
	}

	err := errors.New(msg)
	r.AddFailure(err)
}

// AddErrorMessage adds a validation error message to the Result
// If the message is empty or whitespace, nothing is added
func (r *Result) AddErrorMessage(msg string) {
	r.AddFailureMessage(msg)
}

// AddFailure adds a validation failure to the Result
// If the failure is nil, nothing is added
func (r *Result) AddFailure(failure error) {
	if failure == nil {
		return
	}
	r.failures = append(r.failures, failure)
}

// AddError adds a validation error to the Result
// If the error is nil, nothing is added
func (r *Result) AddError(err error) {
	r.AddFailure(err)
}

// AddParameterIsNotValidError adds a validation error indicating that a parameter is not valid
// If the name is empty or whitespace, nothing is added
func (r *Result) AddParameterIsNotValidError(name string, value any) {
	if strings.TrimSpace(name) == "" {
		return
	}

	var valMsg string
	if v, ok := value.(string); ok {
		valMsg = fmt.Sprintf("'%s'", v)
	} else {
		valMsg = fmt.Sprintf("%v", value)
	}

	err := fmt.Errorf("%s is not valid: %s", name, valMsg)
	r.AddError(err)
}

// IsSuccess returns true when no error has been added to the result
func (r Result) IsSuccess() bool {
	return len(r.failures) == 0
}

// IsFailure returns true when any error has been added to the result
func (r Result) IsFailure() bool {
	return !r.IsSuccess()
}

// GetFailures returns a list of all failures in the result
// If no failures are found, returns an empty slice
func (r Result) GetFailures() []error {
	return r.failures
}

// GetErrors returns a list of all errors in the result
// If no errors are found, returns an empty slice
func (r Result) GetErrors() []error {
	return r.GetFailures()
}

// GetFailureMessages returns a list of all failure messages in the result
// If no failures are found, returns an empty slice
func (r Result) GetFailureMessages() []string {
	s := make([]string, 0, len(r.failures))
	for _, err := range r.failures {
		s = append(s, err.Error())
	}
	return s
}

// GetErrorMessages returns a list of all error messages in the result
// If no errors are found, returns an empty slice
func (r Result) GetErrorMessages() []string {
	return r.GetFailureMessages()
}
