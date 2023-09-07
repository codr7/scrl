package scrl

import (
	"io"
)

type DequeForm struct {
	ItemsForm
}

func NewDequeForm(pos Pos, items ...Form) *DequeForm {
	return new(DequeForm).Init(pos, items...)
}

func (self *DequeForm) Init(pos Pos, items ...Form) *DequeForm {
	self.ItemsForm.Init(pos, items)
	return self
}

func (self *DequeForm) Emit(args *Forms, vm *Vm, env Env, ret bool) error {
	if err := self.ItemsForm.Emit(args, vm, env, false); err != nil {
		return err
	}

	vm.Emit(NewDequeOp(self.pos, len(self.items)), true)
	return nil
}

func (self DequeForm) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "["); err != nil {
		return err
	}

	if err := self.ItemsForm.Dump(out); err != nil {
		return err
	}

	if _, err := io.WriteString(out, "]"); err != nil {
		return err
	}

	return nil
}
