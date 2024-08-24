package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type point struct {
	x, y int
}

type segment struct {
	p0, p1 point
}

// возвращает true, если отрезки прересекаются. если конец одного отрезка пренадлежит другому,
// это не считается пересечением
func intersect(s1, s2 segment) bool {
	v1 := s1.p1.x - s1.p0.x
	w1 := s1.p1.y - s1.p0.y
	v2 := s2.p1.x - s2.p0.x
	w2 := s2.p1.y - s2.p0.y
	v3 := s2.p0.x - s1.p0.x
	w3 := s2.p0.y - s1.p0.y

	z1 := w1*v3 - v1*w3
	z2 := w2*v3 - v2*w3
	z3 := v1*w2 - w1*v2

	// if debugEnable {
		log.Println("z1,z2,z3:", z1, z2, z3)
	// }

	sign := func(a int) int {
		if a < 0 {
			return -1
		}
		if a > 0 {
			return 1
		}
		return 0
	}

	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	check := func(u1, u2 int) bool {
		if sign(u1) != sign(u2) {
			return false
		}
		return abs(u1) < abs(u2)
	}

	return check(z1, z3) && check(z2, z3)
}

func slowSolve(pp []point) bool {
	for i := range pp {
		p0 := pp[i]
		for j := i + 1; j < len(pp); j++ {
			p1 := pp[j]
			for k := j + 1; k < len(pp); k++ {
				p2 := pp[k]
				for l := k + 1; l < len(pp); l++ {
					p3 := pp[l]
					if intersect(segment{p0, p1}, segment{p2, p3}) ||
						intersect(segment{p0, p3}, segment{p1, p2}) ||
						intersect(segment{p0, p2}, segment{p1, p3}) {
						return true
					}
				}
			}
		}
	}
	return false
}

func run(in io.Reader, out io.Writer) {
	var n int
	_, err := fmt.Fscan(in, &n)
	if err != nil {
		panic(err)
	}

	pp := make([]point, n)
	for i := range pp {
		var x, y int
		_, err := fmt.Fscan(in, &x, &y)
		if err != nil {
			panic(err)
		}
		pp[i] = point{x, y}
	}

	res := slowSolve(pp)

	if res {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}
