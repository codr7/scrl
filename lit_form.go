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

func (self *LitForm) Quote(vm *Vm) Val {
	return self.val
}

func (self LitForm) Eq(other Form) bool {
	f, ok := other.(*LitForm)
	return ok && f.val.Eq(self.val)
}

func (self LitForm) Dump(out io.Writer) error {
	return self.val.Dump(out)
}
