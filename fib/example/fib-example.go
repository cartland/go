package main

import (
	"fmt"
	"github.com/cartland/go/fib"
)

func main() {
	f := fib.NewMemoizer()
	for i := 0; i < 20; i++ {
		fmt.Printf("f.Fib(%v) %v\n", i, f.Fib(i))
	}
}
