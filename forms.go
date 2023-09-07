package scrl

type Forms struct {
	Deque[Form]
}

func (self *Forms) Emit(vm *Vm, env Env) error {
	for len(self.items) > 0 {
		if err := self.PopFront().Emit(self, vm, env); err != nil {
			return err
		}
	}

	return nil
}
