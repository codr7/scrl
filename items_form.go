package scrl

import (
	"bufio"
)

type ItemsForm struct {
	BasicForm
	items []Form
}

func (self *ItemsForm) Init(pos Pos, items []Form) *ItemsForm {
	self.BasicForm.Init(pos)
	self.items = items
	return self
}

func (self ItemsForm) Emit(args *Forms, vm *VM, env Env) error {
	var fargs Forms
	fargs.Init(self.items)
	return fargs.Emit(vm, env)
}

func (self ItemsForm) Dump(out *bufio.Writer) error {
	for i, f := range self.items {
		if i > 0 {
			if _, err := out.WriteRune(' '); err != nil {
				return err
			}
		}

		if err := f.Dump(out); err != nil {
			return err
		}
	}

	return nil
}
