package scrl

import (
	"fmt"
)

type Error struct {
	pos     Pos
	message string
}

func NewError(pos Pos, spec string, args ...interface{}) Error {
	return Error{pos, fmt.Sprintf(spec, args...)}
}

func (self Error) Message() string {
	return self.message
}

func (self Error) Error() string {
	return fmt.Sprintf("Error in %v: %v", self.pos, self.message)
}
