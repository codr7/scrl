package scrl

type Pair struct {
	left, right Val
}

func NewPair(left, right Val) Pair {
	return Pair{left, right}
}
