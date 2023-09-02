package scrl

import ()

const (
	VERSION = 1
)

type VM struct {
	Trace bool
	Ops   []Op

	main Task
	task *Task
}

func (self *VM) Init() *VM {
	self.main.Init(nil)
	self.task = &self.main
	return self
}

func (self *VM) EmitPC() PC {
	return len(self.Ops)
}

func (self *VM) Emit(trace bool) PC {
	if self.Trace && trace {
		self.Ops[self.Emit(false)] = &TraceOp
	}

	pc := self.EmitPC()
	self.Ops = append(self.Ops, nil)
	return pc
}

func (self *VM) Task() *Task {
	return self.task
}

func (self *VM) Eval(pc PC) (PC, error) {
	return self.Ops[pc].Eval(self, pc)
}
