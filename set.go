package scrl

import (
	"fmt"
	"io"
	"strings"
)

type SetItem[K, V any] struct {
	key   K
	value V
}

type Set[K, V any] struct {
	compare Compare[K]
	items   []SetItem[K, V]
}

func NewSet[K, V any](compare Compare[K]) *Set[K, V] {
	return new(Set[K, V]).Init(compare)
}

func (self *Set[K, V]) Init(compare Compare[K]) *Set[K, V] {
	self.compare = compare
	return self
}

func (self Set[K, V]) Index(key K) (int, *V) {
	min, max := 0, len(self.items)

	for min < max {
		i := (min + max) / 2
		it := self.items[i]

		switch self.compare(key, it.key) {
		case Lt:
			max = i
		case Eq:
			return i, &it.value
		case Gt:
			min = i + 1
		}
	}

	return min, nil
}

func (self Set[K, V]) Clone() *Set[K, V] {
	dst := NewSet[K, V](self.compare)
	dst.items = make([]SetItem[K, V], len(self.items))
	copy(dst.items, self.items)
	return dst
}

func (self Set[K, V]) Find(key K) *V {
	_, found := self.Index(key)
	return found
}

func (self Set[K, V]) Each(f func(key, val interface{}) bool) bool {
	for _, it := range self.items {
		if !f(it.key, it.value) {
			return false
		}
	}

	return true
}

func (self *Set[K, V]) Add(key K, val V) bool {
	i, found := self.Index(key)

	if found != nil {
		return false
	}

	self.items = append(self.items, SetItem[K, V]{})
	copy(self.items[i+1:], self.items[i:])
	self.items[i] = SetItem[K, V]{key, val}
	return true
}

func (self *Set[K, V]) Remove(key K) *V {
	i, found := self.Index(key)

	if found != nil {
		self.items = self.items[:i+copy(self.items[i:], self.items[i+1:])]
	}

	return found
}

func (self Set[K, V]) Keys() []K {
	out := make([]K, len(self.items))

	for i, it := range self.items {
		out[i] = it.key
	}

	return out
}

func (self Set[K, V]) Values() []V {
	out := make([]V, len(self.items))

	for i, it := range self.items {
		out[i] = it.value
	}

	return out
}

func (self Set[K, V]) Len() int {
	return len(self.items)
}

func (self Set[K, V]) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "{"); err != nil {
		return err
	}

	for i, it := range self.items {
		if i > 0 {
			if _, err := io.WriteString(out, " "); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprint(out, "%v %v", it.key, it.value); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, "}"); err != nil {
		return err
	}

	return nil
}

func (self Set[K, V]) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}
