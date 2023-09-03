package scrl

import ()

type MacroBody = func(self *Macro, args *Forms, vm *VM, env Env, pos Pos) error

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

func (self *Macro) Emit(args *Forms, vm *VM, env Env, pos Pos) error {
	return self.body(self, args, vm, env, pos)
}

func (self *Macro) String() string {
	return self.name
}
