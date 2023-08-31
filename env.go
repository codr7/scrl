package scrl

import (
	"fmt"
)

type EnvEachFunc func(id string, v Val)

type Env interface {
	Bind(id string, v Val)
	Find(id string) *Val
	Each(f EnvEachFunc)
	Import(source Env, names ...string) error
}

type BasicEnv struct {
	parent   Env
	bindings map[string]Val
}

func (self *BasicEnv) Init(parent Env) *BasicEnv {
	self.parent = parent
	self.bindings = make(map[string]Val)
	return self
}

func (self *BasicEnv) Bind(id string, v Val) {
	self.bindings[id] = v
}

func (self BasicEnv) Find(id string) *Val {
	v, ok := self.bindings[id]

	if !ok {
		if self.parent != nil {
			return self.parent.Find(id)
		}

		return nil
	}

	return &v
}

func (self BasicEnv) Each(f EnvEachFunc) {
	for k, v := range self.bindings {
		f(k, v)
	}
}

func (self *BasicEnv) Import(source Env, names ...string) error {
	if len(names) == 0 {
		source.Each(func(id string, v Val) {
			self.Bind(id, v)
		})
	} else {
		for _, id := range names {
			v := source.Find(id)

			if v == nil {
				return fmt.Errorf("Unknown identifier: %v", id)
			}

			self.Bind(id, *v)
		}
	}

	return nil
}
