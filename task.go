package scrl

type PC = int

type Task struct {
	Stack Stack
	Env   BasicEnv

	prev, next *Task
	pc         PC
}

func (self *Task) Init(prev *Task) *Task {
	self.prev = prev
	self.Stack.Init(nil)

	var prevEnv Env

	if prev != nil {
		prevEnv = &prev.Env
		prev.next = self
	}

	self.Env.Init(prevEnv)
	return self
}
