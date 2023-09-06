package scrl

import (
	"fmt"
)

type FunBody = func(self *Fun, vm *Vm, stack *Stack, pos Pos, pc Pc) (Pc, error)

type Fun struct {
	name  string
	arity int
	body  FunBody
}

func NewFun(name string, arity int, body FunBody) *Fun {
	return new(Fun).Init(name, arity, body)
}

func (self *Fun) Init(name string, arity int, body FunBody) *Fun {
	self.name = name
	self.arity = arity
	self.body = body
	return self
}

func (self *Fun) Call(vm *Vm, stack *Stack, pos Pos, pc Pc) (Pc, error) {
	return self.body(self, vm, stack, pos, pc)
}

func (self *Fun) String() string {
	return fmt.Sprintf("%v(%v)", self.name, self.arity)
}
