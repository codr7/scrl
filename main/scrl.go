package main

import (
	"github.com/codr7/scrl"
)

func main() {
	var vm scrl.Vm
	vm.Init()
	env := scrl.NewEnv(nil)
	env.Import(&scrl.AbcLib)
	scrl.REPL(&vm, env, scrl.NewStack(nil))
}
