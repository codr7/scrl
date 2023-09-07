package scrl

import (
	"fmt"
	"io"
)

type Op interface {
	Dump(out io.Writer) error
}

type AndOp struct {
	pos     Pos
	falsePc Pc
}

func NewAndOp(pos Pos, falsePc Pc) *AndOp {
	return &AndOp{pos: pos, falsePc: falsePc}
}

func (self AndOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "And %v %v", self.pos, self.falsePc); err != nil {
		return err
	}

	return nil
}

var BenchOp BenchOpT

type BenchOpT struct{}

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

func (self OrOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Or %v %v", self.pos, self.truePc); err != nil {
		return err
	}

	return nil
}

var PairOp PairOpT

type PairOpT struct{}

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

func (self PushOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Push %v %v", self.pos, self.val); err != nil {
		return err
	}

	return nil
}

var RetOp RetOpT

type RetOpT struct{}

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

func (self SetOp) Dump(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Set %v", self.itemCount); err != nil {
		return err
	}

	return nil
}

var StopOp StopOpT

type StopOpT struct{}

func (_ StopOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Stop")
	return err
}

var TraceOp TraceOpT

type TraceOpT struct{}

func (_ TraceOpT) Dump(out io.Writer) error {
	_, err := io.WriteString(out, "Trace")
	return err
}
