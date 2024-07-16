package main

import "fmt"

const modulo = 1e9 + 7

func main() {
	var res int
	for a := 1; a <= 10; a++ {
		fmt.Printf("%2d\n", a)
		for b := a; b <= 12; b++ {
			fmt.Printf("%2d %2d\n", a, b)
			for c := 5; c <= a; c++ {
				fmt.Printf("%2d %2d %2d\n", a, b, c)
				res += 1
				res %= modulo
			}
			fmt.Println("-- -- --")
		}
		fmt.Println("-- --")
	}
	fmt.Println("--")
	fmt.Println(res)
}
