package goerror

import (
	"fmt"
)

type Error interface {
	Error() string
	ErrorWithCause() string
	PrintInput() string
	IsCodeEqual(err error) bool
	WithCause(cause error) Error
	WithInput(input interface{}) Error
	WithExtendMsg(msg string) Error
	StackTrace() string
	Cause() string
}

type GoError struct {
	Status    int
	Code      string
	Msg       string
	ExtendMsg string
	cause     string

	input  interface{}
	frames []*frame
}

func (e *GoError) Error() string {
	if e.cause != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Msg, e.cause)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *GoError) PrintInput() string {
	if e.input == nil {
		return ""
	}

	return fmt.Sprintf("%v", e.input)
}

func (e *GoError) Cause() string {
	return e.cause
}

func (e *GoError) ErrorWithCause() string {
	return fmt.Sprintf("%s - %s", e.Error(), e.Cause())
}

func (e *GoError) IsCodeEqual(err error) bool {
	if ge, ok := err.(*GoError); ok {
		return ge.Code == e.Code
	}

	return false
}

func (e *GoError) WithCause(cause error) Error {
	e.cause = cause.Error()
	e.frames = trace(DefaultStackTraceSkipLine)

	return e
}

func (e *GoError) WithInput(input interface{}) Error {
	e.input = input

	return e
}

func (e *GoError) WithExtendMsg(extendMsg string) Error {
	e.ExtendMsg = extendMsg

	return e
}
