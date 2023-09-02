package scrl

import (
	"io"
)

type ListForm struct {
	ItemsForm
}

func NewListForm(pos Pos, items ...Form) *ListForm {
	return new(ListForm).Init(pos, items...)
}

func (self *ListForm) Init(pos Pos, items ...Form) *ListForm {
	self.ItemsForm.Init(pos, items)
	return self
}

func (self ListForm) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "("); err != nil {
		return err
	}

	if err := self.ItemsForm.Dump(out); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ")"); err != nil {
		return err
	}

	return nil
}
