package scrl

import (
	"fmt"
)

type Pos struct {
	source       string
	line, column int
}

func NewPos(source string, line, column int) Pos {
	var p Pos
	p.Init(source, line, column)
	return p
}

func (self *Pos) Init(source string, line, column int) {
	self.source = source
	self.line = line
	self.column = column
}

func (self *Pos) Source() string {
	return self.source
}

func (self Pos) String() string {
	return fmt.Sprintf("%v@%v:%v", self.source, self.line, self.column)
}
