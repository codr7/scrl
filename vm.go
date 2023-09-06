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

	ops []Op
}

func (self *Vm) Init() *Vm {
	self.Stack.Init(nil)
	self.Env.Init(nil)
	return self
}

func (self *Vm) EmitPc() Pc {
	return len(self.ops)
}

func (self *Vm) Emit(op Op, trace bool) Pc {
	if self.Trace && trace {
		self.Emit(&TraceOp, false)
	}

	pc := self.EmitPc()
	self.ops = append(self.ops, op)
	return pc
}

func (self *Vm) Eval(pc Pc) (Pc, error) {
	return self.ops[pc].Eval(self, pc)
}
