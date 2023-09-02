package scrl

import (
	"io"
	"strings"
)

type Val struct {
	t Type
	d any
}

func NewVal(t Type, d any) Val {
	var v Val
	v.Init(t, d)
	return v
}

func (self *Val) Init(t Type, d any) {
	self.t = t
	self.d = d
}

func (self Val) Compare(other Val) int {
	return self.t.Compare(self, other)
}

func (self Val) Eq(other Val) bool {
	if self.t != other.t {
		return false
	}

	return self.t.Eq(self, other)
}

func (self Val) IsTrue() bool {
	return self.t.IsTrue(self)
}

func (self Val) Emit(args *Forms, vm *VM, env Env, pos Pos) error {
	return self.t.Emit(self, args, vm, env, pos)
}

func (self Val) Write(out io.Writer) error {
	return self.t.Write(self, out)
}

func (self Val) Dump(out io.Writer) error {
	return self.t.Dump(self, out)
}

func (self Val) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

func ValCompare(l, r Val) int {
	return l.Compare(r)
}

type ValSet = Set[Val, Val]

func NewValSet(compare Compare[Val]) *ValSet {
	return NewSet[Val, Val](ValCompare)
}
