package errors

import "fmt"

var (
	ServiceNotFound          = Error{"service not found", ""}
	InvalidInitSign          = Error{"invalid service initialize function", ""}
	InvalidOutArgumentsCount = Error{"invalid output arguments, should be 2", ""}
	InvalidOutSign           = Error{"invalid output arguments, should be (interface, error)", ""}
	NotResolveDep            = Error{"DI cannot resolve dependency", ""}
	CircleLoop               = Error{"DI cannot resolve dependency, circle loop", ""}
	serviceInitTpl           = Error{"service init error: %s", ""}
)

type Error struct {
	msg  string
	data string
}

func (e Error) Error() string {
	return e.msg + " in " + e.data
}

func (e Error) GetMsg() string {
	return e.msg
}

func (e *Error) SetData(d string) {
	e.data = d
}

func ServiceInitError(e error) Error {
	err := serviceInitTpl
	err.msg = fmt.Sprintf(err.msg, e.Error())
	return err
}
