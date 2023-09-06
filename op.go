package scrl

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Op interface {
	Eval(vm *Vm, pc PC) (PC, error)
	Dump(out io.Writer) error
}

type AndOp struct {
	pos     Pos
	falsePC PC
}

func NewAndOp(pos Pos, falsePC PC) *AndOp {
	return &AndOp{pos: pos, falsePC: falsePC}
}

func (self *AndOp) Eval(vm *Vm, pc PC) (PC, error) {
	v := vm.Stack.PeekBack()

	if !v.IsTrue() {
		return vm.Eval(self.falsePC)
	}

	vm.Stack.PopBack()
	return vm.Eval(pc + 1)
}

func (self *AndOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "And %v", self.falsePC); err != nil {
		return err
	}

	return nil
}

var BenchOp BenchOpT

type BenchOpT struct{}

func (self *BenchOpT) Eval(vm *Vm, pc PC) (PC, error) {
	reps := vm.Stack.PopBack().d.(int)
	startTime := time.Now()

	pc++
	startPC := pc

	for i := 0; i < reps; i++ {
		var err error

		if pc, err = vm.Eval(startPC); err != nil {
			return pc, err
		}

		vm.Stack.Clear()
	}

	vm.Stack.PushBack(NewVal(&AbcLib.TimeType, time.Now().Sub(startTime)))
	return pc, nil
}

func (self *BenchOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Bench")
	return err
}

type DequeOp struct {
	pos       Pos
	itemCount int
}

func NewDequeOp(pos Pos, itemCount int) *DequeOp {
	return &DequeOp{pos: pos, itemCount: itemCount}
}

func (self *DequeOp) Eval(vm *Vm, pc PC) (PC, error) {
	d := NewValDeque(vm.Stack.Cut(self.itemCount))
	vm.Stack.PushBack(NewVal(&AbcLib.DequeType, d))
	return vm.Eval(pc + 1)
}

func (self *DequeOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Deque %v", self.itemCount); err != nil {
		return err
	}

	return nil
}

type GotoOp struct {
	pos Pos
	pc  PC
}

func NewGotoOp(pos Pos, pc PC) *GotoOp {
	return &GotoOp{pos: pos, pc: pc}
}

func (self *GotoOp) Eval(vm *Vm, pc PC) (PC, error) {
	return vm.Eval(self.pc)
}

func (self *GotoOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Goto %v", self.pc); err != nil {
		return err
	}

	return nil
}

type IfOp struct {
	pos    Pos
	elsePC PC
}

func NewIfOp(pos Pos, elsePC PC) *IfOp {
	return &IfOp{pos: pos, elsePC: elsePC}
}

func (self *IfOp) Eval(vm *Vm, pc PC) (PC, error) {
	v := vm.Stack.PopBack()

	if v.IsTrue() {
		return vm.Eval(pc + 1)
	}

	return vm.Eval(self.elsePC)
}

func (self *IfOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "If %v", self.elsePC); err != nil {
		return err
	}

	return nil
}

type OrOp struct {
	pos    Pos
	truePC PC
}

func NewOrOp(pos Pos, truePC PC) *OrOp {
	return &OrOp{pos: pos, truePC: truePC}
}

func (self *OrOp) Eval(vm *Vm, pc PC) (PC, error) {
	v := vm.Stack.PeekBack()

	if v.IsTrue() {
		return vm.Eval(self.truePC)
	}

	vm.Stack.PopBack()
	return vm.Eval(pc + 1)
}

func (self *OrOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Or %v", self.truePC); err != nil {
		return err
	}

	return nil
}

type PairOp struct {
	pos Pos
}

func NewPairOp(pos Pos) *PairOp {
	return &PairOp{pos: pos}
}

func (self *PairOp) Eval(vm *Vm, pc PC) (PC, error) {
	r := vm.Stack.PopBack()
	l := vm.Stack.PopBack()
	vm.Stack.PushBack(NewVal(&AbcLib.PairType, NewPair(l, r)))
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

func (self *PrimCallOp) Eval(vm *Vm, pc PC) (PC, error) {
	pc, err := self.target.Call(vm, self.pos, pc+1)

	if err != nil {
		return pc, err
	}

	return vm.Eval(pc)
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

func (self *PushOp) Eval(vm *Vm, pc PC) (PC, error) {
	vm.Stack.PushBack(self.val)
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

func (self *SetOp) Eval(vm *Vm, pc PC) (PC, error) {
	s := NewValSet(vm.Stack.Cut(self.itemCount))
	vm.Stack.PushBack(NewVal(&AbcLib.SetType, s))
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

func (self *StopOpT) Eval(vm *Vm, pc PC) (PC, error) {
	return pc, nil
}

func (self *StopOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Stop")
	return err
}

var TraceOp TraceOpT

type TraceOpT struct{}

func (self *TraceOpT) Eval(vm *Vm, pc PC) (PC, error) {
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
