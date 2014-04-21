package main

import "fmt"
import "github.com/cartland/hello-go/fib"

func main() {
	fmt.Printf("Hello, world.\n")
  m := fib.NewMemoizer()
  for i := 0; i < 10; i++ {
    fmt.Printf("f.Fib($v) $v", i, f.Fib(i))
  }
}
