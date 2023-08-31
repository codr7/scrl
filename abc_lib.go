package scrl

var AbcLib AbcLibT

func init() {
	AbcLib.Init("abc")
}

type IntType struct {
	BasicType
}

func (_ IntType) IsTrue(v Val) bool {
	return v.d != 0
}

type StrType struct {
	BasicType
}

func (_ StrType) IsTrue(v Val) bool {
	return len(v.d.(string)) > 0
}

type AbcLibT struct {
	BasicLib
	IntType IntType
	StrType StrType
}

func (self *AbcLibT) Init(name string) *AbcLibT {
	self.BasicLib.Init(name)
	self.IntType.Init("Int")
	self.StrType.Init("Str")
	return self
}
