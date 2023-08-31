package scrl

import ()

const (
	VERSION = 1
)

type VM struct {
	Trace bool

	Ops  []Op
	main Task
	task *Task
}

func (self *VM) Init() *VM {
	self.main.Init(nil)
	self.task = &self.main
	return self
}

func (self *VM) Emit(n int) PC {
	pc := len(self.Ops)

	for i := 0; i < n; i++ {
		self.Ops = append(self.Ops, nil)
	}

	return pc
}

func (self *VM) Eval(pc PC) (PC, error) {
	return self.Ops[pc].Eval(self, pc)
}
