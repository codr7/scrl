package scrl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func REPL(vm *VM) {
	fmt.Printf("scrl v%v\n\n", VERSION)
	in := bufio.NewScanner(os.Stdin)
	out := os.Stdout
	var buf bytes.Buffer

	for {
		if _, err := out.WriteString("  "); err != nil {
			log.Fatal(err)
		}

		if !in.Scan() {
			if err := in.Err(); err != nil {
				log.Fatal(err)
			}

			break
		}

		line := in.Text()

		if line == "" {
			pos := NewPos("repl", 1, 1)
			pc := vm.EmitPC()
			var forms Forms

			if err := ReadForms(vm, bufio.NewReader(&buf), &forms, &pos); err != nil {
				fmt.Println(err)
				buf.Reset()
				goto NEXT
			}

			buf.Reset()

			if err := forms.Emit(vm, &vm.Env); err != nil {
				fmt.Println(err)
				goto NEXT
			}

			vm.Ops[vm.Emit(true)] = &StopOp

			if _, err := vm.Eval(pc); err != nil {
				fmt.Println(err)
				vm.Stack.Clear()
				goto NEXT
			}
		NEXT:
			if err := vm.Stack.Dump(out); err != nil {
				log.Fatal(err)
			}

			if _, err := io.WriteString(out, "\n"); err != nil {
				log.Fatal(err)
			}
		} else if _, err := fmt.Fprintln(&buf, line); err != nil {
			log.Fatal(err)
		}
	}
}
