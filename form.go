package scrl

import (
	"bufio"
)

type Form interface {
	Pos() Pos
	Emit(args *Forms, vm *VM, env Env) error
	Dump(out *bufio.Writer) error
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
}

func (self BasicForm) Pos() Pos {
	return self.pos
}
