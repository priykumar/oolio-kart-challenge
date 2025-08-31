package myerror

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput        = errors.New("invalid input")         // 400
	ErrNotFound            = errors.New("product not found")     // 404
	ErrValidationException = errors.New("validation exception")  // 422
	ErrInternalServer      = errors.New("internal server error") // 500
)

var Code2Err map[int]error = map[int]error{
	400: errors.New("invalid input"),
	404: errors.New("product not found"),
	422: errors.New("validation exception"),
	500: errors.New("internal server error"),
}

type KartError struct {
	Code int
	Msg  string
}

func (e KartError) Error() string {
	err := Code2Err[e.Code]
	if err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, err)
	}
	return e.Msg
}

func (e KartError) Unwrap() error {
	return Code2Err[e.Code]
}

func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsIValidationException(err error) bool {
	return errors.Is(err, ErrValidationException)
}

func IsInternalServer(err error) bool {
	return errors.Is(err, ErrInternalServer)
}
