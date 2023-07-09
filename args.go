package clear

import "strconv"

type Arg interface {
	Names() []string
	Positional() bool
	Saturated() bool
	Consume(Tokenizer) error
}

type Converter[T any] interface {
	From(string) (T, error)
}

type floatConv struct {
}

func (f floatConv) From(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

type intConv struct {
}

func (f intConv) From(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

type uintConv struct {
}

func (f uintConv) From(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

type strConv struct{}

func (f strConv) From(s string) (string, error) {
	return s, nil
}

var FloatArg = NewArg[float64, floatConv]
var IntArg = NewArg[int64, intConv]
var UintArg = NewArg[uint64, uintConv]
var StringArg = NewArg[string, strConv]

type argImpl[T any] struct {
	names      []string
	min        int
	max        int
	positional bool
	values     []T
	converter  Converter[T]
	eachaction func(int, T) error
	allaction  func([]T) error
}

type ArgOption[T any] func(*argImpl[T])

func Min[T any](min int) ArgOption[T] {
	return func(ai *argImpl[T]) {
		ai.min = min
	}
}

func Max[T any](max int) ArgOption[T] {
	return func(ai *argImpl[T]) {
		ai.max = max
	}
}

func EachValAction[T any](f func(index int, val T) error) ArgOption[T] {
	return func(ai *argImpl[T]) {
		ai.eachaction = f
	}
}

func AliasOption[T any](name string) ArgOption[T] {
	return func(ai *argImpl[T]) {
		ai.names = append(ai.names, name)
	}
}

func NewArg[T any, C Converter[T]](name string, options ...ArgOption[T]) Arg {
	positional := false
	for name[0] == '-' {
		name = name[1:]
		positional = true
	}
	var c C
	ai := &argImpl[T]{
		names:      []string{name},
		converter:  c,
		values:     make([]T, 0),
		positional: positional,
	}

	for _, option := range options {
		option(ai)
	}

	return ai
}

func (na *argImpl[T]) Names() []string {
	return na.names
}

func (na *argImpl[T]) Positional() bool {
	return na.positional
}

func (na *argImpl[T]) Saturated() bool {
	return len(na.values) == na.max && na.max != 0
}

func (na *argImpl[T]) Satisfied() bool {
	return len(na.values) >= na.min
}

func (na *argImpl[T]) Consume(tzer Tokenizer) error {
	na.values = make([]T, 0)

	for tzer.HasNext() {
		token := tzer.Next()

		if token == "" || token == "--" {
			if !na.Satisfied() {
				return &insufficientArg{
					min:     na.min,
					current: len(na.values),
				}
			}
			return nil
		}

		if val, err := na.converter.From(token); err != nil {
			return err
		} else {
			index := len(na.values)
			na.values = append(na.values, val)
			if na.eachaction != nil {
				if err := na.eachaction(index, val); err != nil {
					return err
				}
			}
		}

		if na.Saturated() {
			break
		}
	}

	if !na.Satisfied() {
		return &insufficientArg{
			min:     na.min,
			current: len(na.values),
		}
	}

	if na.allaction != nil {
		return na.allaction(na.values)
	}
	return nil
}

func (na *argImpl[T]) Get() any {
	return na.values
}
