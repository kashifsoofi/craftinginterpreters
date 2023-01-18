package lox

type scope map[string]bool

func newScope() scope {
	return scope{}
}

type stack struct {
	values []scope
}

func newStack() *stack {
	return &stack{
		values: make([]scope, 0),
	}
}

func (s *stack) push(value map[string]bool) {
	s.values = append(s.values, value)
}

func (s *stack) pop() map[string]bool {
	l := len(s.values)
	if l == 0 {
		panic("cannot pop value from empty stack")
	}

	value := s.values[l-1]
	s.values = s.values[:l-1]
	return value
}

func (s *stack) empty() bool {
	return len(s.values) == 0
}

func (s *stack) peek() map[string]bool {
	l := len(s.values)
	return s.values[l-1]
}
