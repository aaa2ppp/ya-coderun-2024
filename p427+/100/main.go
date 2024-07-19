package main

import (
	"bytes"
	"fmt"
	"slices"
)

type queue []int

func (q *queue) push(v int) {
	*q = append(*q, v)
}

func (q *queue) pop() int {
	old := *q
	v := old[0]
	*q = old[1:]
	return v
}

func (q queue) len() int {
	return len(q)
}

func (q queue) front() int {
	return q[0]
}

func main() {
	a, b, c := 79, 89, 97
	// a, b, c := 89, 94, 99
	// a, b, c := 94, 97, 99
	// a, b, c := 97, 98, 99

	var matrix [][]byte
	matrix = append(matrix, bytes.Repeat([]byte{'.'}, a))
	count := 0
	last_val := 0
	iter_count := 0

	put := func(val int, let byte) bool {
		last_val = val
		i := val / a
		j := val % a
		for i >= len(matrix) {
			matrix = append(matrix, slices.Clone(matrix[len(matrix)-1]))
		}
		old_count := count
		if matrix[i][j] == '.' {			
			count++
		}
		matrix[i][j] = let
		iter_count++
		return old_count != count
	}

	axb := a * b
	var b_queue queue

	put(0, 'c')
	c_next := c
	b_queue.push(b)

	for count < a && (b_queue.front() < axb || c_next < axb) {
		if b_queue.front() <= c_next {
			v := b_queue.pop()
			if put(v, 'b') {
				b_queue.push(v + b)
			}
			continue
		}

		v := c_next
		c_next += c

		if put(v, 'c') {
			b_queue.push(v + b)
		}
	}

	for _, row := range matrix {
		fmt.Printf("%s\n", row)
	}
	fmt.Println(last_val, iter_count)
}
