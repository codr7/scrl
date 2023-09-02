package scrl

import (
	"fmt"
)

type PrimBody = func(self *Prim, vm *VM, pos Pos, pc PC) (PC, error)

type Prim struct {
	name  string
	arity int
	body  PrimBody
}

func NewPrim(name string, arity int, body PrimBody) *Prim {
	return new(Prim).Init(name, arity, body)
}

func (self *Prim) Init(name string, arity int, body PrimBody) *Prim {
	self.name = name
	self.arity = arity
	self.body = body
	return self
}

func (self *Prim) Call(vm *VM, pos Pos, pc PC) (PC, error) {
	return self.body(self, vm, pos, pc+1)
}

func (self *Prim) String() string {
	return fmt.Sprintf("%v(%v)", self.name, self.arity)
}
