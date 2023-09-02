package scrl

import (
	"fmt"
	"io"
	"os"
)

type Op interface {
	Eval(vm *VM, pc PC) (PC, error)
	Dump(out io.Writer) error
}

type PairOp struct {
	pos Pos
}

func NewPairOp(pos Pos) *PairOp {
	return &PairOp{pos: pos}
}

func (self *PairOp) Eval(vm *VM, pc PC) (PC, error) {
	r := vm.task.Stack.PopBack()
	l := vm.task.Stack.PopBack()
	vm.task.Stack.PushBack(NewVal(&AbcLib.PairType, NewPair(l, r)))
	return vm.Eval(pc + 1)
}

func (self *PairOp) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "Pair"); err != nil {
		return err
	}

	return nil
}

type PrimCallOp struct {
	pos    Pos
	target *Prim
}

func NewPrimCallOp(pos Pos, target *Prim) *PrimCallOp {
	return &PrimCallOp{pos: pos, target: target}
}

func (self *PrimCallOp) Eval(vm *VM, pc PC) (PC, error) {
	return self.target.Call(vm, self.pos, pc)
}

func (self *PrimCallOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Prim %v", self.target); err != nil {
		return err
	}

	return nil
}

type PushOp struct {
	pos Pos
	val Val
}

func NewPushOp(pos Pos, val Val) *PushOp {
	return &PushOp{pos: pos, val: val}
}

func (self *PushOp) Eval(vm *VM, pc PC) (PC, error) {
	vm.task.Stack.PushBack(self.val)
	return vm.Eval(pc + 1)
}

func (self *PushOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Push %v", self.val); err != nil {
		return err
	}

	return nil
}

type SetOp struct {
	pos       Pos
	itemCount int
}

func NewSetOp(pos Pos, itemCount int) *SetOp {
	return &SetOp{pos: pos, itemCount: itemCount}
}

func (self *SetOp) Eval(vm *VM, pc PC) (PC, error) {
	s := NewValSet(ValCompare)

	for _, v := range vm.task.Stack.Cut(self.itemCount) {
		s.Add(v)
	}

	vm.task.Stack.PushBack(NewVal(&AbcLib.SetType, s))
	return vm.Eval(pc + 1)
}

func (self *SetOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Set %v", self.itemCount); err != nil {
		return err
	}

	return nil
}

var StopOp StopOpT

type StopOpT struct{}

func (self *StopOpT) Eval(vm *VM, pc PC) (PC, error) {
	return pc, nil
}

func (self *StopOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Stop")
	return err
}

var TraceOp TraceOpT

type TraceOpT struct{}

func (self *TraceOpT) Eval(vm *VM, pc PC) (PC, error) {
	pc++
	fmt.Fprintf(os.Stdout, "%v ", pc)
	vm.Ops[pc].Dump(os.Stdout)
	io.WriteString(os.Stdout, "\n")
	return vm.Eval(pc)
}

func (self *TraceOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Trace")
	return err
}
