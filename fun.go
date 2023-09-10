package scrl

import (
	"fmt"
)

type FunBody = func(self *Fun, vm *Vm, stack *Stack, pos Pos, pc Pc) (Pc, error)

type Fun struct {
	name string
	args FunArgs
	ret  Type
	body FunBody
	pc   Pc
}

func NewFun(name string, args FunArgs, ret Type, body FunBody) *Fun {
	return new(Fun).Init(name, args, ret, body)
}

func (self *Fun) Init(name string, args FunArgs, ret Type, body FunBody) *Fun {
	self.name = name
	self.args = args
	self.ret = ret
	self.body = body
	self.pc = -1
	return self
}

func (self Fun) Arity() int {
	return len(self.args.items)
}

func (self *Fun) Call(vm *Vm, stack *Stack, pos Pos, pc Pc) (Pc, error) {
	return self.body(self, vm, stack, pos, pc)
}

func (self Fun) String() string {
	return fmt.Sprintf("%v()", self.name)
}

type FunArg struct {
	name string
	t    Type
}

type FunArgs struct {
	items []FunArg
}

func (self *FunArgs) Add(name string, t Type) *FunArgs {
	self.items = append(self.items, FunArg{name, t})
	return self
}
