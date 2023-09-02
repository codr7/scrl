package scrl

type Order = int

const (
	Lt = Order(-1)
	Eq = Order(0)
	Gt = Order(1)
)

type Compare[T any] func(l, r T) int
