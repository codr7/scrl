package scrl

import (
	"fmt"
	"io"
	"os"
	"time"
)

func Eval(vm *Vm, stack *Stack, pc Pc) (Pc, error) {
NEXT:
	switch op := vm.ops[pc].(type) {
	case *AndOp:
		v := stack.PeekBack()

		if v.IsTrue() {
			stack.PopBack()
			pc++
		} else {
			pc = op.falsePc
		}

		goto NEXT
	case *BenchOpT:
		reps := stack.PopBack().d.(int)
		startTime := time.Now()

		pc++
		startPc := pc

		for i := 0; i < reps; i++ {
			var err error

			if pc, err = Eval(vm, stack, startPc); err != nil {
				return pc, err
			}

			stack.Clear()
		}

		stack.PushBack(NewVal(&AbcLib.TimeType, time.Now().Sub(startTime)))
		goto NEXT
	case *CallOp:
		var err error
		pc, err = op.target.Call(vm, stack, op.pos, pc+1)

		if err != nil {
			return pc, err
		}

		goto NEXT
	case *DequeOp:
		d := NewValDeque(stack.Cut(op.itemCount))
		stack.PushBack(NewVal(&AbcLib.DequeType, d))
		pc++
		goto NEXT
	case *FunArgOp:
		stack.PushBack(vm.calls.PeekBack().args[op.index])
		pc++
		goto NEXT
	case *GotoOp:
		pc = op.pc
		goto NEXT
	case *IfOp:
		v := stack.PopBack()

		if v.IsTrue() {
			pc++
		} else {
			pc = op.elsePc
		}

		goto NEXT
	case *OrOp:
		v := stack.PeekBack()

		if v.IsTrue() {
			pc = op.truePc
		} else {
			stack.PopBack()
			pc++
		}

		goto NEXT
	case *PairOpT:
		r := stack.PopBack()
		l := stack.PopBack()
		stack.PushBack(NewVal(&AbcLib.PairType, NewPair(l, r)))
		pc++
		goto NEXT
	case *PushOp:
		stack.PushBack(op.val)
		pc++
		goto NEXT
	case *RetOpT:
		pc = vm.calls.PopBack().retPc
		goto NEXT
	case *SetOp:
		s := NewValSet(stack.Cut(op.itemCount))
		stack.PushBack(NewVal(&AbcLib.SetType, s))
		pc++
		goto NEXT
	case *StopOpT:
		//Exit
	case *TraceOpT:
		pc++
		fmt.Fprintf(os.Stdout, "%v ", pc)
		vm.ops[pc].Dump(os.Stdout)
		io.WriteString(os.Stdout, "\n")
		goto NEXT
	default:
		return pc, fmt.Errorf("Invalid op: %v %v", pc, op)
	}

	return pc, nil
}
