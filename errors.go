package jsont

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedCharacter = errors.New("unexpected character")
	ErrUnexpectedNewline   = errors.New("unexpected newline")
	ErrUnexpectedFinished  = errors.New("unexpected finished")
	ErrTooShort            = errors.New("input too short")
	ErrUnclosed            = errors.New("unclosed")
	ErrIllegalUnicode      = errors.New("illegal unicode")
	ErrIllegalEscape       = errors.New("illegal escape character")
)

type ErrorWithPosition struct {
	raw      []byte
	inner    error
	position int
}

func errWithPosition(pos int, err error) error {
	return &ErrorWithPosition{
		inner:    err,
		position: pos,
	}
}

func (err ErrorWithPosition) Error() string {
	return fmt.Sprintf("%s at `%s`...", err.inner.Error(), err.raw[:err.position])
}
