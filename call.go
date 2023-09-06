package scrl

type Call struct {
	pos    Pos
	target *Fun
	args   []Val
	retPc  Pc
}

func NewCall(pos Pos, target *Fun, args []Val, retPc Pc) Call {
	return Call{pos, target, args, retPc}
}
