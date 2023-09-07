package scrl

import ()

type MacroBody = func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos, ret bool) error

type Macro struct {
	name string
	body MacroBody
}

func NewMacro(name string, body MacroBody) *Macro {
	return new(Macro).Init(name, body)
}

func (self *Macro) Init(name string, body MacroBody) *Macro {
	self.name = name
	self.body = body
	return self
}

func (self *Macro) Emit(args *Forms, vm *Vm, env Env, pos Pos, ret bool) error {
	return self.body(self, args, vm, env, pos, ret)
}

func (self *Macro) String() string {
	return self.name
}
