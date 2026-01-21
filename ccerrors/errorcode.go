package ccerrors

import (
	"fmt"
	"strconv"
)

const (
	unkownPrefix       = 90
	validationPrefix   = 91
	persisetncePrefix  = 92
	externalCallPrefix = 93
	invalidCodePrefix  = 99
)

type ErrorCode struct {
	prefix int
	value  int
	name   string
}

func new(name string, prefix, value int) ErrorCode {
	return ErrorCode{
		prefix: prefix,
		value:  value,
		name:   name,
	}
}

func FromCode(code int) ErrorCode {
	codeStr := fmt.Sprintf("%02d", code)
	codePrefix := codeStr[:2]
	actualCode := codeStr[2:]

	prefix, _ := strconv.Atoi(codePrefix)
	val, _ := strconv.Atoi(actualCode)

	return ErrorCode{
		prefix: prefix,
		value:  val,
		name:   "",
	}
}

func (ec ErrorCode) Prefix() int {
	return ec.prefix
}

func (ec ErrorCode) Value() int {
	return ec.value
}

func (ec ErrorCode) Name() string {
	return ec.name
}
func (ec ErrorCode) Code() int {
	codeStr := fmt.Sprintf("%d%d", ec.prefix, ec.value)
	code, _ := strconv.Atoi(codeStr)
	return code
}

func (ec ErrorCode) IsUnknown() bool {
	return ec.prefix == unkownPrefix
}

func (ec ErrorCode) IsValidation() bool {
	return ec.prefix == validationPrefix
}

func (ec ErrorCode) IsPersistence() bool {
	return ec.prefix == persisetncePrefix
}

func (ec ErrorCode) IsExternalCall() bool {
	return ec.prefix == externalCallPrefix
}

func (ec ErrorCode) IsInvalidCode() bool {
	return ec.prefix == invalidCodePrefix
}

func NewUnknownErrorCode(value int, name string) ErrorCode {
	return new(name, unkownPrefix, value)
}

func NewValidationErrorCode(value int, name string) ErrorCode {
	return new(name, validationPrefix, value)
}

func NewPersistenceErrorCode(value int, name string) ErrorCode {
	return new(name, persisetncePrefix, value)
}
func NewExternalCallErrorCode(value int, name string) ErrorCode {
	return new(name, externalCallPrefix, value)
}
func NewInvalidCodeErrorCode(value int, name string) ErrorCode {
	return new(name, invalidCodePrefix, value)
}
