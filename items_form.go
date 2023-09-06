package scrl

import (
	"io"
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

func (self ItemsForm) Emit(args *Forms, vm *Vm, env Env) error {
	var fargs Forms
	fargs.Init(self.items)
	return fargs.Emit(vm, env)
}

func (self ItemsForm) Dump(out io.Writer) error {
	for i, f := range self.items {
		if i > 0 {
			if _, err := io.WriteString(out, " "); err != nil {
				return err
			}
		}

		if err := f.Dump(out); err != nil {
			return err
		}
	}

	return nil
}
