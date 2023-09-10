package scrl

import (
	"io"
)

type IdForm struct {
	BasicForm
	name string
}

func NewIdForm(pos Pos, name string) *IdForm {
	return new(IdForm).Init(pos, name)
}

func (self *IdForm) Init(pos Pos, name string) *IdForm {
	self.BasicForm.Init(pos)
	self.name = name
	return self
}

func (self *IdForm) Emit(args *Forms, vm *Vm, env Env) error {
	found := env.Find(self.name)

	if found == nil {
		return NewError(self.pos, "Unknown identifier: %v", self.name)
	}

	if err := found.Emit(args, vm, env, self.pos); err != nil {
		return err
	}

	return nil
}

func (self IdForm) Quote(vm *Vm) Val {
	return NewVal(&AbcLib.SymType, vm.Sym(self.name))
}

func (self IdForm) Eq(other Form) bool {
	f, ok := other.(*IdForm)
	return ok && f.name == self.name
}

func (self IdForm) Dump(out io.Writer) error {
	_, err := io.WriteString(out, self.name)
	return err
}

func (self IdForm) String() string {
	return self.name
}
