package ccerrors

import (
	"fmt"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/sts-solutions/base-code/ccerrors/withinnererr"
)

type StackTrace struct {
	ExcludedStringPatterns []string
}

func (st StackTrace) GetString(err error) string {
	if err == nil {
		return ""
	}
	stackTrace := st.GetStrings(err)
	for len(stackTrace) == 0 {
		if withInnerErr, ok := err.(withinnererr.WithInnerErr); ok {
			innerErr := withInnerErr.GetInnerErr()
			if innerErr == nil {
				break
			}
			err = innerErr
			stackTrace = st.GetStrings(err)
			continue
		}

		return ""
	}

	resp := strings.Join(stackTrace, "\n")
	return resp
}

func (st StackTrace) GetStrings(err error) []string {
	if err == nil {
		return []string{}
	}

	stackTrace := st.stack(err)
	resp := []string{}

	for _, stackItem := range stackTrace {
		item := fmt.Sprintf("%+v", stackItem)
		if st.shouldExcludeStackTraceItem(item) {
			continue
		}
		resp = append(resp, item)
	}

	return resp
}

func (st StackTrace) stack(err error) (stackTrace errors.StackTrace) {
	if err == nil {
		return nil
	}

	for {
		stackError, ok := err.(interface {
			StackTrace() errors.StackTrace
		})
		if ok {
			stackTrace = stackError.StackTrace()
		}

		u, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			break
		}
		err = u.Unwrap()
	}

	return stackTrace
}

func (st StackTrace) shouldExcludeStackTraceItem(item string) bool {
	for _, exc := range st.ExcludedStringPatterns {
		if match, _ := regexp.Match(exc, []byte(item)); match {
			return true
		}
	}

	return false
}
