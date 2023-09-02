package scrl

type Forms struct {
	items []Form
}

func (self *Forms) Init(items []Form) *Forms {
	self.items = items
	return self
}

func (self *Forms) PushFront(form Form) {
	self.items = append(self.items, nil)
	copy(self.items[:len(self.items)-1], self.items[1:])
	self.items[0] = form
}

func (self *Forms) PushBack(form Form) {
	self.items = append(self.items, form)
}

func (self Forms) PeekFront() Form {
	i := len(self.items)

	if i == 0 {
		return nil
	}

	return self.items[0]
}

func (self *Forms) PopFront() Form {
	i := len(self.items)

	if i == 0 {
		return nil
	}

	f := self.items[0]
	self.items = self.items[1:]
	return f
}

func (self *Forms) PopBack() Form {
	i := len(self.items) - 1
	f := self.items[i]
	self.items = self.items[:i]
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
		if err := self.PopFront().Emit(self, vm, env); err != nil {
			return err
		}
	}

	return nil
}
