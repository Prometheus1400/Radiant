package environment

type Environment[T any] struct {
	parentEnvironment *Environment[T]
	table             map[string]T
	accessCount       map[string]int64
}

func NewEnvironment[T any](parent *Environment[T]) *Environment[T] {
	return &Environment[T]{table: make(map[string]T, 0), parentEnvironment: parent}
}

func (s *Environment[T]) Define(name string) {
	var zero T
	s.table[name] = zero
	// s.accessCount[name] = 0
}

func (s *Environment[T]) Get(name string) (T, bool) {
	val, exists := s.table[name]

	if !exists && s.parentEnvironment != nil {
		return s.parentEnvironment.Get(name)
	}
	return val, exists
}

func (s *Environment[T]) Set(name string, value T) {
	s.table[name] = value
}
