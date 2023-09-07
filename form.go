package scrl

import (
	"io"
)

type Form interface {
	Pos() Pos
	Emit(args *Forms, vm *Vm, env Env) error
	Dump(out io.Writer) error
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
