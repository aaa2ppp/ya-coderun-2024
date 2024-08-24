package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type topping struct {
	state  int
	people [2][]int
}

func solve(n, m int, ab [][2]int) ([]int, bool) {
	people := make([]int, n+1)
	toppings := make([]topping, m+1)

	for i := range ab {
		p, t := ab[i][0], ab[i][1]
		if debugEnable {
			log.Println(p, t)
		}

		var want int
		if t > 0 {
			want = 1
		} else {
			want = 0
			t = -t
		}
		toppings[t].people[want] = append(toppings[t].people[want], p)

		if debugEnable {
			log.Println(toppings, people)
		}

		for want != toppings[t].state {
			old := people[p]

			if old == 0 {
				people[p] = t
				break
			}

			people[p] = t
			if tolgeTopping(toppings, people, old) {
				break
			}

			people[p] = old
			if tolgeTopping(toppings, people, t) {
				break
			}

			return nil, false
		}
	}

	count := 0
	for i := range toppings {
		if toppings[i].state == 1 {
			count++
		}
	}

	res := make([]int, 0, count)
	for i := 1; i < len(toppings); i++ {
		if toppings[i].state == 1 {
			res = append(res, i)
		}
	}

	return res, true
}

// XXX
var bak_states []int
var bak_people [][2]int

func tolgeTopping(toppings []topping, people []int, t int) (ok bool) {

	// XXX backup
	if bak_states == nil {
		bak_states = make([]int, 0, len(toppings))
		bak_people = make([][2]int, 0, len(people))
	} else {
		bak_states = bak_states[:0]
		bak_people = bak_people[:0]
	}

	restore := func() {
		for i := len(bak_states) - 1; i >= 0; i-- {
			node := bak_states[i]
			toppings[node].state += 1
			toppings[node].state &= 1
		}
		for i := len(bak_people) - 1; i >= 0; i-- {
			p, node := bak_people[i][0], bak_people[i][1]
			people[p] = node
		}
	}

	n := len(toppings)
	visited := make([]bool, n)
	frontier := make(queue[int], 0, n)

	frontier.push(t)
	visited[t] = true

	for !frontier.empty() {
		node := frontier.pop()
		for _, p := range toppings[node].people[toppings[node].state] {
			if people[p] == 0 {
				bak_people = append(bak_people, [2]int{p, people[p]})
				people[p] = node
				continue
			}
			if visited[people[p]] {
				restore() // XXX
				return false
			}
			frontier.push(people[p])
			visited[people[p]] = true

			bak_people = append(bak_people, [2]int{p, people[p]})
			people[p] = node
		}
		bak_states = append(bak_states, node)
		toppings[node].state += 1
		toppings[node].state &= 1
		for _, p := range toppings[node].people[toppings[node].state] {
			if people[p] == node {
				bak_people = append(bak_people, [2]int{p, people[p]})
				people[p] = 0
			}
		}
	}
	return true
}

type queue[T any] []T

func (q queue[T]) empty() bool {
	return len(q) == 0
}

func (q *queue[T]) push(v T) {
	*q = append(*q, v)
}

func (q *queue[T]) pop() T {
	old := *q
	v := old[0]
	*q = old[1:]
	return v
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, k, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	ab := make([][2]int, 0, k)
	for i := 0; i < k; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		ab = append(ab, [2]int{a, b})
	}

	res, ok := solve(n, m, ab)

	if !ok {
		bw.WriteString("-1\n")
	} else {
		wo := writeOpts{sep: ' ', end: '\n'}
		writeInt(bw, len(res), wo)
		writeInts(bw, res, wo)
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

func scanThreeInt(sc *bufio.Scanner) (v1, v2, v3 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	return v1, v2, v3, err
}

func scanFourInt(sc *bufio.Scanner) (v1, v2, v3, v4 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	if err == nil {
		v4, err = scanInt(sc)
	}
	return v1, v2, v3, v4, err
}

func scanInts(sc *bufio.Scanner, a []int) error {
	for i := range a {
		v, err := scanInt(sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type writeOpts struct {
	sep byte
	end byte
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))

	if err == nil && opts.end != 0 {
		bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

// ----------------------------------------------------------------------------

func gcd(a, b int) int {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
