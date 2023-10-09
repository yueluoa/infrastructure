package gerror

import "fmt"

type Error struct {
	code  Code
	error error
}

func New(msg string) CodeError {
	return newError(msg)
}

func Newf(format string, args ...interface{}) CodeError {
	return newError(fmt.Sprintf(format, args...))
}

func newError(msg string) *Error {
	return &Error{
		code:  CommonError.code,
		error: new(msg),
	}
}

func WithCodeError(error CodeError) CodeError {
	return &Error{
		code:  error.Code(),
		error: error,
	}
}

func WithCode(code Code, msg string) CodeError {
	return &Error{
		code:  code,
		error: new(msg),
	}
}

func IsWithCode(code Code, err error) bool {
	if xe, ok := err.(CodeError); ok {
		return xe.Code() == code
	}

	return false
}

func GetWithCode(err error) Code {
	if xe, ok := err.(CodeError); ok {
		return xe.Code()
	}

	return 0
}

func Wrap(err error, msg string) CodeError {
	return withMessage(err, msg)
}

func Wrapf(err error, format string, args ...interface{}) CodeError {
	return withMessage(err, fmt.Sprintf(format, args...))
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func withMessage(err error, msg string) CodeError {
	if err == nil {
		return nil
	}
	ce, ok := err.(CodeError)
	if !ok {
		ce = newError(msg)
	}
	return &Error{
		code:  ce.Code(),
		error: newEntry(err, msg),
	}
}

func (e *Error) Code() Code { return e.code }

func (e *Error) Error() string {
	if e.error == nil {
		return "nil"
	}
	return e.error.Error()
}
func (e *Error) Unwrap() error { return e.error }

func (e *Error) Cause() error { return e.error }
