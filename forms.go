package scrl

type Forms struct {
	items []Form
}

func (self *Forms) Init(items []Form) *Forms {
	self.items = items
	return self
}

func (self *Forms) Push(form Form) {
	self.items = append(self.items, form)
}

func (self Forms) Peek() Form {
	i := len(self.items)

	if i == 0 {
		return nil
	}

	return self.items[0]
}

func (self *Forms) Pop() Form {
	i := len(self.items)

	if i == 0 {
		return nil
	}

	f := self.items[0]
	self.items = self.items[1:]
	return f
}

func (self Forms) Items() []Form {
	return self.items
}

func (self Forms) Len() int {
	return len(self.items)
}

func (self *Forms) Emit(vm *VM, env Env) error {
	for len(self.items) > 0 {
		if err := self.Pop().Emit(self, vm, env); err != nil {
			return err
		}
	}

	return nil
}
