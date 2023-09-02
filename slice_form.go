package scrl

import (
	"io"
)

type SliceForm struct {
	ItemsForm
}

func NewSliceForm(pos Pos, items ...Form) *SliceForm {
	return new(SliceForm).Init(pos, items...)
}

func (self *SliceForm) Init(pos Pos, items ...Form) *SliceForm {
	self.ItemsForm.Init(pos, items)
	return self
}

func (self *SliceForm) Emit(args *Forms, vm *VM, env Env) error {
	if err := self.ItemsForm.Emit(args, vm, env); err != nil {
		return err
	}

	vm.Ops[vm.Emit(true)] = NewSliceOp(self.pos, len(self.items))
	return nil
}

func (self SliceForm) Dump(out io.Writer) error {
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
