package scrl

import (
	"bufio"
	"fmt"
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

func (self *IdForm) Emit(args *Forms, vm *VM, env Env) error {
	found := env.Find(self.name)

	if found == nil {
		return fmt.Errorf("Unknown identifier: %v", self.name)
	}

	if err := found.Emit(args, vm, env, self.pos); err != nil {
		return err
	}

	return nil
}

func (self IdForm) Dump(out *bufio.Writer) error {
	_, err := out.WriteString(self.name)
	return err
}

func (self IdForm) String() string {
	return self.name
}
