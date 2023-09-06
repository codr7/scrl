package scrl

import (
	"io"
)

type LitForm struct {
	BasicForm
	val Val
}

func NewLitForm(pos Pos, v Val) *LitForm {
	return new(LitForm).Init(pos, v)
}

func (self *LitForm) Init(pos Pos, v Val) *LitForm {
	self.BasicForm.Init(pos)
	self.val = v
	return self
}

func (self LitForm) Emit(args *Forms, vm *Vm, env Env) error {
	vm.Emit(NewPushOp(self.pos, self.val), true)
	return nil
}

func (self LitForm) Dump(out io.Writer) error {
	return self.val.Dump(out)
}
