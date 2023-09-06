package scrl

import (
	"io"
)

type PairForm struct {
	BasicForm
	left, right Form
}

func NewPairForm(pos Pos, left, right Form) *PairForm {
	return new(PairForm).Init(pos, left, right)
}

func (self *PairForm) Init(pos Pos, left, right Form) *PairForm {
	self.BasicForm.Init(pos)
	self.left = left
	self.right = right
	return self
}

func (self PairForm) Emit(args *Forms, vm *Vm, env Env) error {
	if err := self.left.Emit(args, vm, env); err != nil {
		return err
	}

	if err := self.right.Emit(args, vm, env); err != nil {
		return err
	}

	vm.Ops[vm.Emit(true)] = NewPairOp(self.pos)
	return nil
}

func (self PairForm) Dump(out io.Writer) error {
	if err := self.left.Dump(out); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ":"); err != nil {
		return err
	}

	if err := self.right.Dump(out); err != nil {
		return err
	}

	return nil
}
