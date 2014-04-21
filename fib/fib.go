/*
 * Copyright 2014 Chris Cartland
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package fib

type Fibber interface {
	Fib(n int) int
}

func New() Fibber {
	return new(recursiveFibber)
}

type recursiveFibber struct{}

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
	f.fibs = make(map[int]int)
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
	memoized = f.Fib(n-1) + f.Fib(n-2)
	f.fibs[n] = memoized
	return
}

type memoizer struct {
	fibs map[int]int
}
