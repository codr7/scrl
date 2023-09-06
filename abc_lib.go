package scrl

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var AbcLib AbcLibT

func init() {
	AbcLib.Init("abc")
}

type AbcLibT struct {
	BasicLib
	BoolType  BoolType
	DequeType DequeType
	IntType   IntType
	MacroType MacroType
	MetaType  BasicType
	PairType  PairType
	PrimType  PrimType
	SetType   SetType
	StrType   StrType
	SymType   SymType
	TimeType  TimeType
}

func (self *AbcLibT) Init(name string) *AbcLibT {
	self.BasicLib.Init(name)

	self.BindType(&self.BoolType, "Bool")
	self.BindType(&self.DequeType, "Deque")
	self.BindType(&self.IntType, "Int")
	self.BindType(&self.MacroType, "Macro")
	self.BindType(&self.MetaType, "Meta")
	self.BindType(&self.PairType, "Pair")
	self.BindType(&self.PrimType, "Prim")
	self.BindType(&self.SetType, "Set")
	self.BindType(&self.StrType, "Str")
	self.BindType(&self.SymType, "Sym")
	self.BindType(&self.TimeType, "Time")

	self.Bind("T", NewVal(&self.BoolType, true))
	self.Bind("F", NewVal(&self.BoolType, false))

	self.BindMacro("and",
		func(_ *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			andPc := vm.Emit(nil, true)

			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			vm.ops[andPc] = NewAndOp(pos, vm.EmitPc())
			return nil
		})

	self.BindMacro("bench",
		func(_ *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			vm.Emit(&BenchOp, true)

			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			vm.Emit(&StopOp, true)
			return nil
		})

	self.BindMacro("if",
		func(_ *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			ifPc := vm.Emit(nil, true)

			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			elsePc := vm.EmitPc()

			if args.Len() > 0 {
				next := args.PeekFront()

				if f, ok := next.(*IdForm); ok && f.name == "else" {
					args.PopFront()
					gotoPc := vm.Emit(nil, true)
					elsePc = vm.EmitPc()

					if err := args.PopFront().Emit(args, vm, env); err != nil {
						return err
					}

					vm.ops[gotoPc] = NewGotoOp(pos, vm.EmitPc())
				}
			}

			vm.ops[ifPc] = NewIfOp(pos, elsePc)
			return nil
		})

	self.BindMacro("or",
		func(_ *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			orPc := vm.Emit(nil, true)

			if err := args.PopFront().Emit(args, vm, env); err != nil {
				return err
			}

			vm.ops[orPc] = NewOrOp(pos, vm.EmitPc())
			return nil
		})

	self.BindPrim("=", 2,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			r := vm.Stack.PopBack()
			l := vm.Stack.PopBack()
			vm.Stack.PushBack(NewVal(&self.BoolType, l.Eq(r)))
			return pc, nil
		})

	self.BindPrim("<", 2,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			r := vm.Stack.PopBack()
			l := vm.Stack.PopBack()
			vm.Stack.PushBack(NewVal(&self.BoolType, l.Compare(r) == -1))
			return pc, nil
		})

	self.BindPrim(">", 2,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			r := vm.Stack.PopBack()
			l := vm.Stack.PopBack()
			vm.Stack.PushBack(NewVal(&self.BoolType, l.Compare(r) == 1))
			return pc, nil
		})

	self.BindPrim("+", 2,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			r := vm.Stack.PopBack()
			l := vm.Stack.PopBack()
			vm.Stack.PushBack(NewVal(&self.IntType, l.d.(int)+r.d.(int)))
			return pc, nil
		})

	self.BindPrim("-", 2,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			r := vm.Stack.PopBack()
			l := vm.Stack.PopBack()
			vm.Stack.PushBack(NewVal(&self.IntType, l.d.(int)-r.d.(int)))
			return pc, nil
		})

	self.BindPrim("milliseconds", 1,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			n := vm.Stack.PopBack().d.(int)
			vm.Stack.PushBack(NewVal(&self.TimeType, time.Duration(n)*time.Millisecond))
			return pc, nil
		})

	self.BindPrim("say", 1,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			vm.Stack.PopBack().Write(os.Stdout)
			io.WriteString(os.Stdout, "\n")
			return pc, nil
		})

	self.BindPrim("sleep", 1,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			time.Sleep(vm.Stack.PopBack().d.(time.Duration))
			return pc, nil
		})

	self.BindPrim("trace", 0,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			vm.Trace = !vm.Trace
			vm.Stack.PushBack(NewVal(&self.BoolType, vm.Trace))
			return pc, nil
		})

	self.BindPrim("type-of", 1,
		func(_ *Prim, vm *Vm, pos Pos, pc Pc) (Pc, error) {
			v := vm.Stack.PopBack()
			vm.Stack.PushBack(NewVal(&self.MetaType, v.t))
			return pc, nil
		})

	return self
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

type ValDeque = Deque[Val]

func NewValDeque(items []Val) *ValDeque {
	return NewDeque[Val](items)
}

type DequeType struct {
	BasicType
}

func (_ DequeType) IsTrue(v Val) bool {
	return v.d.(*ValDeque).Len() > 0
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

type MacroType struct {
	BasicType
}

func (_ MacroType) Emit(v Val, args *Forms, vm *Vm, env Env, pos Pos) error {
	return v.d.(*Macro).Emit(args, vm, env, pos)
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

func (_ PrimType) Emit(v Val, args *Forms, vm *Vm, env Env, pos Pos) error {
	p := v.d.(*Prim)

	for i := 0; i < p.arity; i++ {
		if err := args.PopFront().Emit(args, vm, env); err != nil {
			return err
		}
	}

	vm.Emit(NewPrimCallOp(pos, p), true)
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

func NewValSet(items []Val) *ValSet {
	return NewSet[Val](ValCompare, items)
}

type SetType struct {
	BasicType
}

func (_ SetType) IsTrue(v Val) bool {
	return v.d.(*ValSet).Len() > 0
}

type SymType struct {
	BasicType
}

func (_ SymType) Compare(l, r Val) int {
	return strings.Compare(l.d.(*Sym).name, r.d.(*Sym).name)
}

func (_ SymType) Dump(v Val, out io.Writer) error {
	_, err := fmt.Fprintf(out, "'%v", v.d.(*Sym).name)
	return err
}

type TimeType struct {
	BasicType
}

func (_ TimeType) Compare(l, r Val) int {
	return l.d.(time.Time).Compare(r.d.(time.Time))
}

func (_ TimeType) IsTrue(v Val) bool {
	return !v.d.(time.Time).IsZero()
}
