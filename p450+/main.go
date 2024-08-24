package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type stack []int

func (s stack) peek() int {
	return s[len(s)-1]
}

func (s *stack) push(v int) {
	*s = append(*s, v)
}

func (s *stack) pop() {
	old := *s
	*s = old[:len(old)-1]
}

func reversed(aa []int) []int {
	bb := make([]int, len(aa))
	for i, j := 0, len(aa)-1; i < len(aa); i, j = i+1, j-1 {
		bb[j] = aa[i]
	}
	return bb
}

func solve(aa, bb, cc []int) int {
	if debugEnable {
		log.Println("------------")
		log.Printf("bb: %3d", bb)
		log.Printf("cc: %3d", cc)
	}

	n := len(aa)

	// TODO: достаточно простого стека, в который складываем минимумы (мы нигде не используем values кроме дебага)
	ss := stack(make([]int, 0, n))

	sum := cc[len(cc)-1]
	for i := 1; i < len(bb)-1; i++ {
		sum += bb[i]
	}

	ss.push(sum)

	for i := len(bb) - 2; i > 1; i-- {
		sum -= bb[i]
		sum += cc[i]
		ss.push(min(sum, ss.peek()))
	}

	if debugEnable {
		log.Printf("ss: --- %3d", reversed(ss))
		log.Printf("aa: %3d", aa)
	}

	minimum := aa[0] + ss.peek()

	aa_sum := aa[0]
	bb_sum := 0
	for i := 1; i < len(aa)-2; i++ {
		aa_sum += aa[i]
		bb_sum += bb[i]
		ss.pop()
		minimum = min(minimum, aa_sum+ss.peek()-bb_sum)
	}

	if debugEnable {
		log.Println("minimum:", minimum)
	}

	return minimum
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	aa := make([]int, n)
	bb := make([]int, n)
	cc := make([]int, n)

	if err := scanInts(sc, aa); err != nil {
		panic(err)
	}
	if err := scanInts(sc, bb); err != nil {
		panic(err)
	}
	if err := scanInts(sc, cc); err != nil {
		panic(err)
	}

	r1 := solve(aa, bb, cc)
	r2 := solve(aa, cc, bb)
	r3 := solve(bb, aa, cc)
	r4 := solve(bb, cc, aa)
	r5 := solve(cc, aa, bb)
	r6 := solve(cc, bb, aa)

	res := min(r1, r2, r3, r4, r5, r6)

	writeInt(bw, res, writeOpts{end: '\n'})
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
