package scrl

import (
	"fmt"
	"io"
)

type Type interface {
	Init(name string)
	Name() string
	String() string

	Compare(l, r Val) int
	Emit(v Val, args *Forms, vm *VM, env Env, pos Pos) error
	Eq(l, r Val) bool
	IsTrue(v Val) bool
	Dump(v Val, out io.Writer) error
	Write(v Val, out io.Writer) error
}

type BasicType struct {
	name string
}

func (self *BasicType) Init(name string) {
	self.name = name
}

func (self *BasicType) Name() string {
	return self.name
}

func (self *BasicType) String() string {
	return self.name
}

func (_ BasicType) Emit(v Val, args *Forms, vm *VM, env Env, pos Pos) error {
	vm.Ops[vm.Emit(true)] = NewPushOp(pos, v)
	return nil
}

func (_ BasicType) Compare(l, r Val) int {
	return 0
}

func (_ BasicType) Eq(l, r Val) bool {
	return l.d == r.d
}

func (_ BasicType) IsTrue(_ Val) bool {
	return true
}

func (_ BasicType) Dump(v Val, out io.Writer) error {
	_, err := fmt.Fprint(out, v.d)
	return err
}

func (_ BasicType) Write(v Val, out io.Writer) error {
	_, err := fmt.Fprint(out, v.d)
	return err
}
