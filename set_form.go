package scrl

import (
	"io"
)

type SetForm struct {
	ItemsForm
}

func NewSetForm(pos Pos, items ...Form) *SetForm {
	return new(SetForm).Init(pos, items...)
}

func (self *SetForm) Init(pos Pos, items ...Form) *SetForm {
	self.ItemsForm.Init(pos, items)
	return self
}

func (self *SetForm) Emit(args *Forms, vm *Vm, env Env) error {
	if err := self.ItemsForm.Emit(args, vm, env); err != nil {
		return err
	}

	vm.Ops[vm.Emit(true)] = NewSetOp(self.pos, len(self.items))
	return nil
}

func (self SetForm) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "{"); err != nil {
		return err
	}

	if err := self.ItemsForm.Dump(out); err != nil {
		return err
	}

	if _, err := io.WriteString(out, "}"); err != nil {
		return err
	}

	return nil
}
