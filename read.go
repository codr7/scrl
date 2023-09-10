package scrl

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

func ReadForms(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	for {
		if err := ReadForm(vm, in, out, pos); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func SkipWhitespace(vm *Vm, in *bufio.Reader, pos *Pos) error {
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

func ReadForm(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
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
	case '[':
		return ReadDeque(vm, in, out, pos)
	case '{':
		return ReadSet(vm, in, out, pos)
	case '\'':
		return ReadQuote(vm, in, out, pos)
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

	return NewError(*pos, "Invalid syntax: %v", c)
}

func ReadBody(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos, closingChar rune) error {
	for {
		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				return NewError(*pos, "Missing %v", closingChar)
			}

			return err
		}

		if c == closingChar {
			pos.column++
			break
		} else {
			in.UnreadRune()
		}

		if err := ReadForm(vm, in, out, pos); err != nil {
			if err == io.EOF {
				return NewError(*pos, "Missing %v", closingChar)
			}

			return err
		}
	}

	return nil
}

func ReadDeque(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var body Forms

	if err := ReadBody(vm, in, &body, pos, ']'); err != nil {
		return err
	}

	out.PushBack(NewDequeForm(fpos, body.items...))
	return nil
}

func readId(vm *Vm, in *bufio.Reader, pos *Pos) (string, error) {
	var buf strings.Builder

	for {
		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}

			return "", err
		}

		if c == '(' || c == ')' || c == '{' || c == '}' || c == '[' || c == ']' ||
			unicode.IsSpace(c) || unicode.IsControl(c) {
			in.UnreadRune()
			break
		}

		buf.WriteRune(c)
		pos.column++
	}

	return buf.String(), nil
}

func ReadId(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	s, err := readId(vm, in, pos)

	if err != nil {
		return err
	}

	out.PushBack(NewIdForm(fpos, s))
	return nil
}

func ReadInt(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
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

func ReadList(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var body Forms

	if err := ReadBody(vm, in, &body, pos, ')'); err != nil {
		return err
	}

	out.PushBack(NewListForm(fpos, body.items...))
	return nil
}

func ReadPair(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	left := out.PopBack()

	if err := ReadForm(vm, in, out, pos); err != nil {
		if err == io.EOF {
			return NewError(*pos, "Invalid pair")
		}

		return err
	}

	right := out.PopBack()
	out.PushBack(NewPairForm(fpos, left, right))
	return nil
}

func ReadQuote(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++

	if err := ReadForm(vm, in, out, pos); err != nil {
		return err
	}

	out.PushBack(NewLitForm(fpos, out.PopBack().Quote(vm)))
	return nil
}

func ReadSet(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var body Forms

	if err := ReadBody(vm, in, &body, pos, '}'); err != nil {
		return err
	}

	out.PushBack(NewSetForm(fpos, body.items...))
	return nil
}

func ReadStr(vm *Vm, in *bufio.Reader, out *Forms, pos *Pos) error {
	fpos := *pos
	pos.column++
	var buf strings.Builder

	for {

		c, _, err := in.ReadRune()

		if err != nil {
			if err == io.EOF {
				return NewError(*pos, "Open string")
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
