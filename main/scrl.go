package main

import (
	"github.com/codr7/scrl"
)

func main() {
	var vm scrl.VM
	vm.Init()
	vm.Env.Import(&scrl.AbcLib)
	scrl.REPL(&vm)
}
