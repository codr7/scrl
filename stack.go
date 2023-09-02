package scrl

import (
	"io"
	"strings"
)

type Stack struct {
	Deque[Val]
}

func (self Stack) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "["); err != nil {
		return err
	}

	for i, v := range self.items {
		if i > 0 {
			if _, err := io.WriteString(out, " "); err != nil {
				return err
			}
		}

		if err := v.Dump(out); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, "]"); err != nil {
		return err
	}

	return nil
}

func (self Stack) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}
