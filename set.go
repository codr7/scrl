package scrl

import (
	"fmt"
	"io"
	"strings"
)

type Set[T any] struct {
	Deque[T]
	compare Compare[T]
}

func NewSet[T any](compare Compare[T], items []T) *Set[T] {
	return new(Set[T]).Init(compare, items)
}

func (self *Set[T]) Init(compare Compare[T], items []T) *Set[T] {
	self.compare = compare
	self.Deque.Init(items)
	return self
}

func (self Set[T]) Index(v T) (int, *T) {
	min, max := 0, len(self.items)

	for min < max {
		i := (min + max) / 2
		it := self.items[i]

		switch self.compare(v, it) {
		case Lt:
			max = i
		case Eq:
			return i, &it
		case Gt:
			min = i + 1
		}
	}

	return min, nil
}

func (self Set[T]) Find(v T) *T {
	_, found := self.Index(v)
	return found
}

func (self *Set[T]) Add(v T) bool {
	i, found := self.Index(v)

	if found != nil {
		return false
	}

	self.items = append(self.items, v)
	copy(self.items[i+1:], self.items[i:])
	self.items[i] = v
	return true
}

func (self *Set[T]) Remove(v T) *T {
	i, found := self.Index(v)

	if found != nil {
		self.items = self.items[:i+copy(self.items[i:], self.items[i+1:])]
	}

	return found
}

func (self Set[T]) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "{"); err != nil {
		return err
	}

	for i, it := range self.items {
		if i > 0 {
			if _, err := io.WriteString(out, " "); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprint(out, it); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, "}"); err != nil {
		return err
	}

	return nil
}

func (self Set[T]) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}
