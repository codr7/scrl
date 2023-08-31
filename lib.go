package scrl

type Lib interface {
	Env
	Name() string
}

type BasicLib struct {
	BasicEnv
	name string
}

func (self *BasicLib) Init(name string) *BasicLib {
	self.BasicEnv.Init(nil)
	self.name = name
	return self
}

func (self BasicLib) Name() string {
	return self.name
}
