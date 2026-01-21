package ccerrors

import "strings"

type MultiError struct {
	errors []error
}

func (me *MultiError) Add(err error) {
	if err == nil {
		return
	}

	if me.errors == nil {
		me.errors = []error{}
	}

	me.errors = append(me.errors, err)
}

func (me *MultiError) HasErrors() bool {
	return len(me.errors) > 0
}

func (me *MultiError) Errors() []error {
	if me.errors == nil {
		return []error{}
	}
	return me.errors
}

func (me *MultiError) Error() string {
	if len(me.errors) == 0 {
		return ""
	}

	if len(me.errors) == 1 {
		return me.errors[0].Error()
	}

	var sb strings.Builder

	for i, err := range me.errors {
		if i != 0 {
			sb.WriteString(";")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}
