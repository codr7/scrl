package scrl

import ()

const (
	VERSION = 1
)

type Pc = int
type Syms = map[string]*Sym
type Bin = func(vm *Vm, stack *Stack, pc Pc) (Pc, error)

type Vm struct {
	Trace bool

	syms  Syms
	ops   []Op
	bins  []Bin
	calls Deque[Call]
}

func (self *Vm) Init() *Vm {
	self.syms = make(Syms)
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

func (self *Vm) Eval(pc Pc, stack *Stack) (Pc, error) {
	return self.ops[pc].Eval(self, stack, pc)
}

func (self *Vm) Compile(pc Pc) error {
	var next Bin
	var err error

	for i := len(self.ops) - 1; i >= pc; pc-- {
		next, err = self.ops[i].Compile(next)

		if err != nil {
			return err
		}

		self.ops[i] = next
	}

	return nil
}

func (self *Vm) Sym(name string) *Sym {
	s := self.syms[name]

	if s == nil {
		s = NewSym(name)
		self.syms[name] = s
	}

	return s
}
