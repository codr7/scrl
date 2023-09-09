package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/codr7/scrl"
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)

		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var vm scrl.Vm
	vm.Init()
	env := scrl.NewEnv(nil)
	env.Import(&scrl.AbcLib)
	scrl.REPL(&vm, env, scrl.NewStack(nil))
}
