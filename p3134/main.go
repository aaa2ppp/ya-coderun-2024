package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func run(in io.Reader, out io.Writer) error {
	var n int
	fmt.Fscan(in, &n)
	fmt.Fprintln(out, phi(n))
	return nil
}

// get from http://e-maxx.ru/algo/euler_function
func phi(n int) int {
	result := n
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			for n%i == 0 {
				n /= i
			}
			result -= result / i
		}
	}
	if n > 1 {
		result -= result / n
	}
	return result
}
