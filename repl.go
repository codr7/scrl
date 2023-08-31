package scrl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func REPL(vm *VM) {
	fmt.Printf("scrl v%v\n\n", VERSION)
	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	var buf bytes.Buffer

	for {
	NEXT:
		if _, err := out.WriteString("  "); err != nil {
			log.Fatal(err)
		}

		out.Flush()

		if !in.Scan() {
			if err := in.Err(); err != nil {
				log.Fatal(err)
			}

			break
		}

		line := in.Text()

		if line == "" {
			pos := NewPos("repl", 1, 1)
			var forms Forms

			if err := ReadForms(vm, bufio.NewReader(&buf), &forms, &pos); err != nil {
				fmt.Println(err)
				buf.Reset()
				goto NEXT
			}

			buf.Reset()
			pc := vm.Emit(0)

			if err := forms.Emit(vm, &vm.task.Env); err != nil {
				fmt.Println(err)
				goto NEXT
			}

			vm.Ops[vm.Emit(1)] = &StopOp

			if _, err := vm.Eval(pc); err != nil {
				fmt.Println(err)
				vm.task.Stack.Clear()
				goto NEXT
			}

			if err := vm.task.Stack.Dump(out); err != nil {
				log.Fatal(err)
			}

			if _, err := out.WriteRune('\n'); err != nil {
				log.Fatal(err)
			}
		} else if _, err := fmt.Fprintln(&buf, line); err != nil {
			log.Fatal(err)
		}
	}
}
