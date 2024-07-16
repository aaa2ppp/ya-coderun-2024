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

const modulo int = 1e9 + 7

func add(a []int, l, r int) {
	for i := l; i <= r; i++ {
		a[i]++
	}
}

func calcMults(a []int, mults [][2]int) {
	for i, v := range a {
		m := mults[i][0] * v
		m %= modulo
		mults[i+1] = [2]int{m, 0}
	}
}

func gcdex(a int, b int) (d, x, y int) {
	if a == 0 {
		return b, 0, 1
	}
	d, x1, y1 := gcdex(b%a, a)
	x = y1 - (b/a)*x1
	y = x1
	return d, x, y
}

func op0(a []int, mults [][2]int, l, r int) {
	add(a, l, r)
	calcMults(a[l:], mults[l:])
}

func op1(_ []int, mults [][2]int, l, r int) int {
	m := mults[l][1]
	if m == 0 {
		_, m, _ = gcdex(mults[l][0], modulo)
		// m = (m + modulo) % modulo
		if m < 0 {
			m += modulo
		}
		mults[l][1] = m
	}
	return m * mults[r+1][0] % modulo
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		return err
	}
	_ = n

	a := make([]int, n+1)
	a[0] = 1
	if err := scanInts(sc, a[1:]); err != nil {
		return err
	}

	q, err := scanInt(sc)
	if err != nil {
		return err
	}

	mults := make([][2]int, len(a)+1)
	mults[0] = [2]int{1, 1}
	calcMults(a, mults)

	// TL ?
	for i := 0; i < q; i++ {
		t, l, r, err := scanThreeInt(sc)
		if err != nil {
			return err
		}
		switch t {
		case 0:
			op0(a, mults, l, r)
		case 1:
			res := op1(a, mults, l, r)
			writeInt(bw, res, writeOpts{end: '\n'})
		default:
			return fmt.Errorf("unknown query type %d", t)
		}
	}

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
