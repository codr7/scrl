package scrl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func ReadForms(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	for {
		if err := ReadForm(vm, in, out, pos); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func SkipWhitespace(vm *VM, in *bufio.Reader, pos *Pos) error {
	for {
		c, _, err := in.ReadRune()

		if err != nil {
			return err
		}

		switch c {
		case '\n':
			pos.line++
			pos.column = 1
		case ' ', '\t':
			pos.column++
		default:
			in.UnreadRune()
			goto EXIT
		}
	}
EXIT:
	return nil
}

func ReadForm(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	if err := SkipWhitespace(vm, in, pos); err != nil {
		return err
	}

	c, _, err := in.ReadRune()

	if err != nil {
		return err
	}

	switch c {
	case '(':
		return ReadList(vm, in, out, pos)
	case '{':
		return ReadSet(vm, in, out, pos)
	case '"':
		return ReadStr(vm, in, out, pos)
	case ':':
		return ReadPair(vm, in, out, pos)
	default:
		if unicode.IsDigit(c) {
			in.UnreadRune()
			return ReadInt(vm, in, out, pos)
		} else if !unicode.IsSpace(c) && !unicode.IsControl(c) {
			in.UnreadRune()
			return ReadId(vm, in, out, pos)
		}
	}

	return fmt.Errorf("Invalid syntax: %v", c)
}

func ReadId(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	var buf strings.Builder

	for {
		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if c == '(' || c == ')' || c == '{' || c == '}' || c == '[' || c == ']' ||
			unicode.IsSpace(c) || unicode.IsControl(c) {
			in.UnreadRune()
			break
		}

		buf.WriteRune(c)
		pos.column++
	}

	out.PushBack(NewIdForm(fpos, buf.String()))
	return nil
}

func ReadInt(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	var v int
	base := 10
	fpos := *pos

	for {

		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if !unicode.IsDigit(c) && (base != 16 || c < 'a' || c > 'f') {
			if err = in.UnreadRune(); err != nil {
				return err
			}

			break
		}

		var dv int

		if base == 16 && c >= 'a' && c <= 'f' {
			dv = 10 + int(c) - int('a')
		} else {
			dv = int(c) - int('0')
		}

		v = v*base + dv
		pos.column++
	}

	out.PushBack(NewLitForm(fpos, NewVal(&AbcLib.IntType, v)))
	return nil
}

func ReadList(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var body Forms

	for {
		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("Open list")
			}

			return err
		}

		if c == ')' {
			pos.column++
			break
		} else {
			in.UnreadRune()
		}

		if err := ReadForm(vm, in, &body, pos); err != nil {
			if err == io.EOF {
				return fmt.Errorf("Open list")
			}

			return err
		}
	}

	out.PushBack(NewListForm(fpos, body.items...))
	return nil
}

func ReadPair(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	left := out.PopBack()

	if err := ReadForm(vm, in, out, pos); err != nil {
		if err == io.EOF {
			return fmt.Errorf("Invalid pair")
		}

		return err
	}

	right := out.PopBack()
	out.PushBack(NewPairForm(fpos, left, right))
	return nil
}

func ReadSet(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var body Forms

	for {
		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("Open set")
			}

			return err
		}

		if c == '}' {
			pos.column++
			break
		} else {
			in.UnreadRune()
		}

		if err := ReadForm(vm, in, &body, pos); err != nil {
			if err == io.EOF {
				return fmt.Errorf("Open set")
			}

			return err
		}
	}

	out.PushBack(NewSetForm(fpos, body.items...))
	return nil
}

func ReadStr(vm *VM, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var buf strings.Builder

	for {

		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("Open string")
			}

			return err
		}

		if c == '"' {
			break
		}

		buf.WriteRune(c)
		pos.column++
	}

	out.PushBack(NewLitForm(fpos, NewVal(&AbcLib.StrType, buf.String())))
	return nil
}
