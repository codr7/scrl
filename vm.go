package scrl

import ()

const (
	VERSION = 1
)

type Pc = int
type Stack = Deque[Val]

type Vm struct {
	Trace bool
	Stack Stack
	Env   BasicEnv
	Ops   []Op
}

func (self *Vm) Init() *Vm {
	self.Stack.Init(nil)
	self.Env.Init(nil)
	return self
}

func (self *Vm) EmitPc() Pc {
	return len(self.Ops)
}

func (self *Vm) Emit(trace bool) Pc {
	if self.Trace && trace {
		self.Ops[self.Emit(false)] = &TraceOp
	}

	pc := self.EmitPc()
	self.Ops = append(self.Ops, nil)
	return pc
}

func (self *Vm) Eval(pc Pc) (Pc, error) {
	return self.Ops[pc].Eval(self, pc)
}
