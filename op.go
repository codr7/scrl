package scrl

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Op interface {
	Compile(next Bin) Bin
	Dump(out io.Writer) error
}

type AndOp struct {
	pos     Pos
	falsePc Pc
}

func NewAndOp(pos Pos, falsePc Pc) *AndOp {
	return &AndOp{pos: pos, falsePc: falsePc}
}

func (self AndOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		v := stack.PeekBack()

		if !v.IsTrue() {
			return vm.Eval(self.falsePc, stack)
		}

		stack.PopBack()
		return next(vm, stack, pc+1)
	}
}

func (self AndOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "And %v %v", self.pos, self.falsePc); err != nil {
		return err
	}

	return nil
}

var BenchOp BenchOpT

type BenchOpT struct{}

func (_ BenchOpT) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		reps := stack.PopBack().d.(int)
		startTime := time.Now()

		pc++
		startPc := pc

		for i := 0; i < reps; i++ {
			var err error

			if pc, err = vm.Eval(startPc, stack); err != nil {
				return pc, err
			}

			stack.Clear()
		}

		stack.PushBack(NewVal(&AbcLib.TimeType, time.Now().Sub(startTime)))
		return vm.Eval(pc, stack)
	}
}

func (_ BenchOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Bench")
	return err
}

type CallOp struct {
	pos    Pos
	target *Fun
}

func NewCallOp(pos Pos, target *Fun) *CallOp {
	return &CallOp{pos: pos, target: target}
}

func (self CallOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		pc, err := self.target.Call(vm, stack, self.pos, pc+1)

		if err != nil {
			return pc, err
		}

		return vm.Eval(pc, stack)
	}
}

func (self CallOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Call %v %v", self.pos, self.target); err != nil {
		return err
	}

	return nil
}

type DequeOp struct {
	pos       Pos
	itemCount int
}

func NewDequeOp(pos Pos, itemCount int) *DequeOp {
	return &DequeOp{pos: pos, itemCount: itemCount}
}

func (self DequeOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		d := NewValDeque(stack.Cut(self.itemCount))
		stack.PushBack(NewVal(&AbcLib.DequeType, d))
		return next(vm, stack, pc+1)
	}
}

func (self DequeOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Deque %v %v", self.pos, self.itemCount); err != nil {
		return err
	}

	return nil
}

type FunArgOp struct {
	pos   Pos
	index int
}

func NewFunArgOp(pos Pos, index int) *FunArgOp {
	return &FunArgOp{pos: pos, index: index}
}

func (self FunArgOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		stack.PushBack(vm.calls.PeekBack().args[self.index])
		return next(vm, stack, pc+1)
	}
}

func (self FunArgOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "FunArg %v %v", self.pos, self.index); err != nil {
		return err
	}

	return nil
}

type GotoOp struct {
	pos Pos
	pc  Pc
}

func NewGotoOp(pos Pos, pc Pc) *GotoOp {
	return &GotoOp{pos: pos, pc: pc}
}

func (self GotoOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		return vm.Eval(self.pc, stack)
	}
}

func (self GotoOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Goto %v %v", self.pos, self.pc); err != nil {
		return err
	}

	return nil
}

type IfOp struct {
	pos    Pos
	elsePc Pc
}

func NewIfOp(pos Pos, elsePc Pc) *IfOp {
	return &IfOp{pos: pos, elsePc: elsePc}
}

func (self IfOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		v := stack.PopBack()

		if v.IsTrue() {
			return next(vm, stack, pc+1)
		}

		return vm.Eval(self.elsePc, stack)
	}
}

func (self IfOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "If %v %v", self.pos, self.elsePc); err != nil {
		return err
	}

	return nil
}

type OrOp struct {
	pos    Pos
	truePc Pc
}

func NewOrOp(pos Pos, truePc Pc) *OrOp {
	return &OrOp{pos: pos, truePc: truePc}
}

func (self OrOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		v := stack.PeekBack()

		if v.IsTrue() {
			return vm.Eval(self.truePc, stack)
		}

		stack.PopBack()
		return next(vm, stack, pc+1)
	}
}

func (self OrOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Or %v %v", self.pos, self.truePc); err != nil {
		return err
	}

	return nil
}

var PairOp PairOpT

type PairOpT struct{}

func (_ PairOpT) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		r := stack.PopBack()
		l := stack.PopBack()
		stack.PushBack(NewVal(&AbcLib.PairType, NewPair(l, r)))
		return next(vm, stack, pc+1)
	}
}

func (_ PairOpT) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "Pair"); err != nil {
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

func (self PushOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		stack.PushBack(self.val)
		return next(vm, stack, pc+1)
	}
}

func (self PushOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Push %v %v", self.pos, self.val); err != nil {
		return err
	}

	return nil
}

var RetOp RetOpT

type RetOpT struct{}

func (_ RetOpT) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		return vm.Eval(vm.calls.PopBack().retPc, stack)
	}
}

func (_ RetOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Ret")
	return err
}

type SetOp struct {
	pos       Pos
	itemCount int
}

func NewSetOp(pos Pos, itemCount int) *SetOp {
	return &SetOp{pos: pos, itemCount: itemCount}
}

func (self SetOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		s := NewValSet(stack.Cut(self.itemCount))
		stack.PushBack(NewVal(&AbcLib.SetType, s))
		return next(vm, stack, pc+1)
	}
}

func (self SetOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Set %v", self.itemCount); err != nil {
		return err
	}

	return nil
}

var StopOp StopOpT

type StopOpT struct{}

func (_ StopOpT) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		return pc, nil
	}
}

func (_ StopOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Stop")
	return err
}

type TailCallOp struct {
	pos    Pos
	target *Fun
}

func NewTailCallOp(pos Pos, target *Fun) *TailCallOp {
	return &TailCallOp{pos: pos, target: target}
}

func (self TailCallOp) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		c := vm.calls.PeekBack()
		c.pos = self.pos
		c.target = self.target
		c.args = stack.Cut(self.target.Arity())
		return vm.bins[self.target.pc](vm, stack, self.target.pc)
	}
}

func (self TailCallOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "TailCall %v %v", self.pos, self.target); err != nil {
		return err
	}

	return nil
}

var TraceOp TraceOpT

type TraceOpT struct{}

func (_ TraceOpT) Compile(next Bin) Bin {
	return func(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
		pc++
		fmt.Fprintf(os.Stdout, "%v ", pc)
		vm.ops[pc].Dump(os.Stdout)
		io.WriteString(os.Stdout, "\n")
		return next(vm, stack, pc)
	}
}

func (_ TraceOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Trace")
	return err
}
