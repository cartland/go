package fib

type Fibber interface {
	Fib(n int) int
}

func New() Fibber {
	return new(recursiveFibber)
}

type recursiveFibber struct {}

func (f *recursiveFibber) Fib(n int) int {
	return fib(n)
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func NewMemoizer() Fibber {
  f := new(memoizer)
  f.fibs = make(map[int] int)
  return f
}

func (f *memoizer) Fib(n int) (memoized int) {
  if n <= 1 {
    return n
  }
  memoized, ok := f.fibs[n]
  if ok {
    return
  }
  memoized = f.Fib(n - 1) + f.Fib(n - 2)
  f.fibs[n] = memoized
  return
}

type memoizer struct {
  fibs map[int] int
}
