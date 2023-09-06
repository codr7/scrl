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

func (self *BasicLib) BindMacro(name string, body MacroBody) {
	self.Bind(name, NewVal(&AbcLib.MacroType, NewMacro(name, body)))
}

func (self *BasicLib) BindFun(name string, arity int, body FunBody) {
	self.Bind(name, NewVal(&AbcLib.FunType, NewFun(name, arity, body)))
}

func (self *BasicLib) BindType(t Type, name string) {
	t.Init(name)
	self.Bind(name, NewVal(&AbcLib.MetaType, t))
}
