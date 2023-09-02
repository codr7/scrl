package scrl

import (
	"io"
)

var AbcLib AbcLibT

func init() {
	AbcLib.Init("abc")
}

type BoolType struct {
	BasicType
}

func (_ BoolType) IsTrue(v Val) bool {
	return v.d.(bool)
}

func (_ BoolType) Dump(v Val, out io.Writer) error {
	var err error

	if v.d.(bool) {
		_, err = io.WriteString(out, "T")
	} else {
		_, err = io.WriteString(out, "F")
	}

	return err
}

type IntType struct {
	BasicType
}

func (_ IntType) IsTrue(v Val) bool {
	return v.d != 0
}

type PrimType struct {
	BasicType
}

func (_ PrimType) Emit(v Val, args *Forms, vm *VM, env Env, pos Pos) error {
	p := v.d.(*Prim)

	for i := 0; i < p.arity; i++ {
		if err := args.Pop().Emit(args, vm, env); err != nil {
			return err
		}
	}

	vm.Ops[vm.Emit(true)] = NewPrimCallOp(pos, p)
	return nil
}

type StrType struct {
	BasicType
}

func (_ StrType) IsTrue(v Val) bool {
	return len(v.d.(string)) > 0
}

type AbcLibT struct {
	BasicLib
	BoolType BoolType
	IntType  IntType
	PrimType PrimType
	StrType  StrType
}

func (self *AbcLibT) Init(name string) *AbcLibT {
	self.BasicLib.Init(name)
	self.BoolType.Init("Bool")
	self.IntType.Init("Int")
	self.PrimType.Init("Prim")
	self.StrType.Init("Str")

	self.Bind("T", NewVal(&self.BoolType, true))
	self.Bind("F", NewVal(&self.BoolType, false))

	self.Bind("trace", NewVal(&self.PrimType, NewPrim("trace", 0,
		func(_ *Prim, vm *VM, pos Pos, pc PC) (PC, error) {
			vm.Trace = !vm.Trace
			vm.task.Stack.Push(NewVal(&self.BoolType, vm.Trace))
			return pc, nil
		})))

	return self
}
