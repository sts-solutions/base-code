package ccerrors

type TransientError struct {
	message    string
	innerError error
	details    map[string]interface{}
}

func NewTransientError(message string, err error, details map[string]interface{}) *TransientError {
	returnedErr := &TransientError{
		message:    "transient error handling nats msg, nak-ing to retry later",
		innerError: err,
		details:    details,
	}
	return returnedErr
}

func (e *TransientError) Error() string {
	return e.message
}

func (e *TransientError) Is(err error) bool {
	_, ok := err.(*TransientError)
	return ok
}
