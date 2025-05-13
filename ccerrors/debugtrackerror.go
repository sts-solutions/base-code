package ccerrors

import "runtime/debug"

type DebugTrackError struct {
	message    string
	innerError error
	details    map[string]interface{}
}

func NewDebugTrackError(message string, err error, details map[string]interface{}) *DebugTrackError {
	returnedErr := &DebugTrackError{
		message:    message,
		innerError: err,
		details:    details,
	}
	returnedErr.details["stack_trace"] = string(debug.Stack())
	return returnedErr
}

func (e *DebugTrackError) Error() string {
	return e.message
}

func (e *DebugTrackError) Is(err error) bool {
	_, ok := err.(*DebugTrackError)
	return ok
}
