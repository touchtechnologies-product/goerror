package goerror

import (
	"fmt"
)

type Error interface {
	Error() string
	ErrorWithCause() string

	PrintInput() string
	StackTrace() string
	Cause() string

	IsCodeEqual(err error) bool

	GetReasons() []*Reason
	AddReason(fieldName, reason string, value interface{})

	Input() interface{}

	WithCause(cause error) Error
	WithInput(input interface{}) Error
	WithKeyValueInput(inputs ...interface{}) Error
	WithExtendMsg(msg string) Error
}

type GoError struct {
	Status    int
	Code      string
	Msg       string
	ExtendMsg string
	cause     string

	reasons []*Reason

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

func (e *GoError) Input() interface{} {
	return e.input
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

func (e *GoError) WithKeyValueInput(keyValues ...interface{}) Error {
	if len(keyValues) == 1 && keyValues[0] == nil {
		return e
	}

	if len(keyValues)%2 != 0 {
		e.input = keyValues
		return e
	}

	fields := map[string]interface{}{}
	for i := 0; i*2 < len(keyValues); i++ {
		key, ok := keyValues[i*2].(string)
		if !ok {
			fields[fmt.Sprintf("errf_%d", i)] = keyValues[i*2+1]
			continue
		}

		fields[key] = keyValues[i*2+1]
	}

	e.input = fields
	return e
}
