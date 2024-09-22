package app

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	internal = "internal"
	notFound = "not_found"
)

type Error struct {
	Code    string
	Message string
	Status  int
	line    int
	file    string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) SetStatus(status int) *Error {
	e.Status = status
	return e
}

func (e *Error) GetStatus() int {
	return e.Status
}

func (e *Error) WrapError(err error) *Error {
	e.Err = err
	return e
}

func ErrorWithCaller(err error) *Error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*Error); ok {
		return err
	}

	return errorf("%w", err)
}

func (e *Error) SetCode(code string) *Error {
	e.Code = strings.ToUpper(strings.ReplaceAll(code, " ", "_"))
	return e
}

func Errorf(format string, args ...any) *Error {
	return errorf(format, args...)
}

func errorf(format string, args ...any) *Error {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("failed to get caller")
	}

	err := fmt.Errorf(format, args...)

	return &Error{
		Message: err.Error(),
		Err:     errors.Unwrap(err),
		line:    line,
		file:    file,
	}
}
