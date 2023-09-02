package scrl

type Deque[T any] struct {
	items []T
}

func NewDeque[T any](items []T) *Deque[T] {
	return new(Deque[T]).Init(items)
}

func (self *Deque[T]) Init(items []T) *Deque[T] {
	self.items = items
	return self
}

func (self *Deque[T]) PushFront(it T) {
	self.items = append(self.items, it)
	copy(self.items[:len(self.items)-1], self.items[1:])
	self.items[0] = it
}

func (self *Deque[T]) PushBack(it T) {
	self.items = append(self.items, it)
}

func (self Deque[T]) PeekFront() T {
	return self.items[0]
}

func (self Deque[T]) PeekBack() T {
	return self.items[len(self.items)-1]
}

func (self *Deque[T]) PopFront() T {
	it := self.items[0]
	self.items = self.items[1:]
	return it
}

func (self *Deque[T]) PopBack() T {
	i := len(self.items) - 1
	it := self.items[i]
	self.items = self.items[:i]
	return it
}

func (self Deque[T]) Items() []T {
	return self.items
}

func (self Deque[T]) Len() int {
	return len(self.items)
}

func (self *Deque[T]) Clear() {
	self.items = nil
}

func (self *Deque[T]) Cut(n int) []T {
	out := make([]T, n)
	i := len(self.items) - n
	copy(out, self.items[i:])
	self.items = self.items[:i]
	return out
}
