package scrl

import (
	"fmt"
	"io"
	"strings"
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

func (_ IntType) Compare(l, r Val) int {
	if l.d.(int) < r.d.(int) {
		return -1
	}

	if l.d.(int) > r.d.(int) {
		return 1
	}

	return 0
}

func (_ IntType) IsTrue(v Val) bool {
	return v.d != 0
}

type PairType struct {
	BasicType
}

func (_ PairType) IsTrue(v Val) bool {
	p := v.d.(Pair)
	return p.left.IsTrue() && p.right.IsTrue()
}

func (_ PairType) Compare(l, r Val) int {
	return l.d.(Pair).left.Compare(r.d.(Pair).left)
}

func (_ PairType) Dump(v Val, out io.Writer) error {
	p := v.d.(Pair)

	if err := p.left.Dump(out); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ":"); err != nil {
		return err
	}

	if err := p.right.Dump(out); err != nil {
		return err
	}

	return nil
}

type PrimType struct {
	BasicType
}

func (_ PrimType) Emit(v Val, args *Forms, vm *VM, env Env, pos Pos) error {
	p := v.d.(*Prim)

	for i := 0; i < p.arity; i++ {
		if err := args.PopFront().Emit(args, vm, env); err != nil {
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

func (_ StrType) Compare(l, r Val) int {
	return strings.Compare(l.d.(string), r.d.(string))
}

func (_ StrType) Dump(v Val, out io.Writer) error {
	_, err := fmt.Fprintf(out, "\"%v\"", v.d.(string))
	return err
}

type ValSet = Set[Val]

func NewValSet(compare Compare[Val]) *ValSet {
	return NewSet[Val](ValCompare)
}

type SetType struct {
	BasicType
}

func (_ SetType) IsTrue(v Val) bool {
	return v.d.(*ValSet).Len() > 0
}

type AbcLibT struct {
	BasicLib
	BoolType BoolType
	IntType  IntType
	PairType PairType
	PrimType PrimType
	SetType  SetType
	StrType  StrType
}

func (self *AbcLibT) Init(name string) *AbcLibT {
	self.BasicLib.Init(name)
	self.BoolType.Init("Bool")
	self.IntType.Init("Int")
	self.PairType.Init("Pair")
	self.PrimType.Init("Prim")
	self.SetType.Init("Set")
	self.StrType.Init("Str")

	self.Bind("T", NewVal(&self.BoolType, true))
	self.Bind("F", NewVal(&self.BoolType, false))

	self.Bind("trace", NewVal(&self.PrimType, NewPrim("trace", 0,
		func(_ *Prim, vm *VM, pos Pos, pc PC) (PC, error) {
			vm.Trace = !vm.Trace
			vm.task.Stack.PushBack(NewVal(&self.BoolType, vm.Trace))
			return pc, nil
		})))

	return self
}
