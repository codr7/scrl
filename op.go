package scrl

import (
	"bufio"
)

type Op interface {
	Eval(vm *VM, pc PC) (PC, error)
	Dump(out *bufio.Writer) error
}

type PushOp struct {
	pos Pos
	val Val
}

func NewPushOp(pos Pos, val Val) *PushOp {
	return &PushOp{pos: pos, val: val}
}

func (self *PushOp) Eval(vm *VM, pc PC) (PC, error) {
	vm.task.Stack.Push(self.val)
	return vm.Eval(pc + 1)
}

func (self *PushOp) Dump(out *bufio.Writer) error {
	if _, err := out.WriteString("Push "); err != nil {
		return err
	}

	if err := self.val.Dump(out); err != nil {
		return err
	}

	return nil
}

type StopOpT struct{}

func (self *StopOpT) Eval(vm *VM, pc PC) (PC, error) {
	return pc, nil
}

func (self *StopOpT) Dump(out *bufio.Writer) error {
	_, err := out.WriteString("Stop")
	return err
}

var StopOp StopOpT
