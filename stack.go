package scrl

import (
	"bufio"
)

type Stack struct {
	items []Val
}

func (self *Stack) Init(items []Val) *Stack {
	self.items = items
	return self
}

func (self *Stack) Push(v Val) {
	self.items = append(self.items, v)
}

func (self Stack) Peek(i int) *Val {
	n := len(self.items) - 1

	if n < i {
		return nil
	}

	return &self.items[n-i]
}

func (self *Stack) Pop() *Val {
	i := len(self.items) - 1

	if i < 0 {
		return nil
	}

	v := self.items[i]
	self.items = self.items[:i]
	return &v
}

func (self *Stack) Cut(n int) []Val {
	l := len(self.items)

	if l < n {
		return nil
	}

	out := make([]Val, n)
	copy(out, self.items[l-n:])
	self.items = self.items[:l-n]
	return out
}

func (self *Stack) Clear() {
	self.items = nil
}

func (self Stack) Len() int {
	return len(self.items)
}

func (self Stack) Dump(out *bufio.Writer) error {
	if _, err := out.WriteRune('['); err != nil {
		return err
	}

	for i, v := range self.items {
		if i > 0 {
			if _, err := out.WriteRune(' '); err != nil {
				return err
			}
		}

		if err := v.Dump(out); err != nil {
			return err
		}
	}

	if _, err := out.WriteRune(']'); err != nil {
		return err
	}

	return nil
}
