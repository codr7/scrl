package scrl

type Sym struct {
	name string
}

func NewSym(name string) *Sym {
	return &Sym{name}
}
