package scrl

type Stack = Deque[Val]

func NewStack(items []Val) *Stack {
	s := new(Stack)
	s.Init(items)
	return s
}
