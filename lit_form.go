package scrl

import (
	"bufio"
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

func (self LitForm) Emit(args *Forms, vm *VM, env Env) error {
	vm.Ops[vm.Emit(true)] = NewPushOp(self.pos, self.val)
	return nil
}

func (self LitForm) Dump(out *bufio.Writer) error {
	return self.val.Dump(out)
}
