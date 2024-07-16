// v4
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type item struct {
	l, r int
	used int
}

func (it *item) count() int {
	return it.r - it.l
}

type queue struct {
	first int
	last  int
	size  int
	count int
	items []item
}

const max_used = 3

func newQueue(n int) *queue {
	items := make([]item, (n+1)/2+1)
	items[0] = item{l: 0, r: 1, used: max_used}
	return &queue{
		size:  1,
		items: items,
		count: 1,
	}
}

func (q *queue) push(l, r int) int {
	if q.size == len(q.items) {
		panic("queue full")
	}

	last := &q.items[q.last]
	if l < last.l {
		panic(fmt.Errorf("l <= last.l: l,r = [%d, %d) last = [%d, %d)", l, r, last.l, last.r))
	}

	if l <= last.r {
		if last.r < r {
			q.count += r - last.r
			last.r = r
		}
	} else {
		q.count += r - l
		q.last = q.nextIdx(q.last)
		q.items[q.last] = item{l, r, max_used}
		q.size++
	}

	return q.last
}

func (q *queue) nextIdx(idx int) int {
	idx++
	if idx == len(q.items) {
		idx = 0
	}
	return idx
}

func (q *queue) pop() {
	if q.size == 0 {
		panic("queue empty")
	}

	q.first = q.nextIdx(q.first)
	q.size--
}

func (q *queue) front() *item {
	if q.size == 0 {
		panic("queue empty")
	}

	return &q.items[q.first]
}

func (q *queue) back() *item {
	if q.size == 0 {
		panic("queue empty")
	}

	return &q.items[q.last]
}

func solve(n int, a, b, c int) int {
	// prepare
	d := gcd(a, gcd(b, c))
	a /= d
	b /= d
	c /= d
	n = (n + d - 1) / d

	q := newQueue(max(a, b, c))
	if debugEnable {
		log.Println(n, a, b, c)
		log.Println(*q)
	}

	i, j, k := 0, 0, 0
	a_next, b_next, c_next := [2]int{a, a + 1}, [2]int{b, b + 1}, [2]int{c, c + 1}

	count := 1
	end := min(a, b, c)
	for a_next[0] < n || b_next[0] < n || c_next[0] < n {
		if it := q.back(); it.count() >= end {
			count = q.count + n - q.back().r
			break
		}

		min_next := min(a_next[0], b_next[0], c_next[0])
		if debugEnable {
			log.Println(">:", []int{i, j, k}, [][2]int{a_next, b_next, c_next}, min_next)
		}

		if a_next[0] == min_next {
			q.push(a_next[0], a_next[1])
			q.items[i].used--
			if i != q.last {
				i = q.nextIdx(i)
			}
			it := &q.items[i]
			a_next = [2]int{it.l + a, min(it.r+a, n)}
		}

		if b_next[0] == min_next {
			q.push(b_next[0], b_next[1])
			q.items[j].used--
			if j != q.last {
				j = q.nextIdx(j)
			}
			it := &q.items[j]
			b_next = [2]int{it.l + b, min(it.r+b, n)}
		}

		if c_next[0] == min_next {
			q.push(c_next[0], c_next[1])
			q.items[k].used--
			if k != q.last {
				k = q.nextIdx(k)
			}
			it := &q.items[k]
			c_next = [2]int{it.l + c, min(it.r+c, n)}
		}

		for q.front().used == 0 {
			q.pop()
		}

		if debugEnable {
			log.Printf("=: %v", *q)
		}

		count = q.count // xxx
	}

	return count
}

func gcd(a, b int) int {
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, a, b, c, err := scanFourInt(sc)
	if err != nil {
		return err
	}

	count := solve(n, a, b, c)
	writeInt(bw, count, writeOpts{end: '\n'})

	return nil
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

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
