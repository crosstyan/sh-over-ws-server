package hub

import "fmt"

const (
	Existed = iota
)

type AppError interface {
	error
	Code() int
	Msg() string
}

type hubError struct {
	code int
	msg  string
}

func (e hubError) Error() string {
	return fmt.Sprintf("hub error: %s, code: %d", e.msg, e.code)
}

func (e hubError) Code() int {
	return e.code
}

func (e hubError) Msg() string {
	return e.msg
}
