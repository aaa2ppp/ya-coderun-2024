package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

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

func solve(n, a, b, c int) int {
	n, a, b, c = prepare(n, a, b, c)
	if debugEnable {
		log.Println("prepared:", n, a, b, c)
	}

	row := bytes.Repeat([]byte{'.'}, a)
	row_count := 0

	prev_i := 0
	count := 0

	put := func(val int, let byte) bool {
		i := val / a
		j := val % a

		if row[j] != '.' {
			return false
		}

		row[j] = let
		if debugEnable {
			log.Printf("%d: %s", i, row)
		}

		count += (i-prev_i)*row_count + 1
		row_count++
		prev_i = i

		return true
	}

	var b_queue queue

	put(0, 'c')
	c_next := c
	b_queue.push(b)

	for row_count < len(row) && (!b_queue.empty() && b_queue.front() <= n || c_next < n) {
		if !b_queue.empty() && b_queue.front() <= c_next {
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

	count += (n/len(row) - prev_i) * row_count
	for j := n % len(row); j < len(row); j++ {
		if row[j] != '.' {
			count--
		}
	}

	return count
}

func prepare(n, a, b, c int) (int, int, int, int) {
	d := gcd(a, gcd(b, c))
	a /= d
	b /= d
	c /= d
	n = (n + d - 1) / d

	if a > b {
		a, b = b, a
	}
	if a > c {
		a, c = c, a
	}
	if b > c {
		b, c = c, b
	}

	return n, a, b, c
}

func gcd(a, b int) int {
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

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

func (q queue) empty() bool {
	return q.len() == 0
}

func (q queue) front() int {
	return q[0]
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
