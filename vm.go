package scrl

import ()

const (
	VERSION = 1
)

type PC = int
type Stack = Deque[Val]

type VM struct {
	Trace bool
	Stack Stack
	Env   BasicEnv
	Ops   []Op
}

func (self *VM) Init() *VM {
	self.Stack.Init(nil)
	self.Env.Init(nil)
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

func (self *VM) Eval(pc PC) (PC, error) {
	return self.Ops[pc].Eval(self, pc)
}
