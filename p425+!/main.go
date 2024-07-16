package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

type Segment struct {
	l, r int
}

type Op struct {
	l, r int
	x    int
}

type SRow []Segment

func (z SRow) And(r SRow, s Segment) SRow {
	z = z[:0]
	for _, s2 := range r {
		if s.r >= s2.l && s2.r >= s.l {
			z = append(z, Segment{max(s.l, s2.l), min(s.r, s2.r)})
		}
	}
	return z
}

func (z SRow) Or(s Segment) SRow {
	for i := len(z) - 1; i >= 0; i-- {
		s2 := z[i]
		if s.r+1 >= s2.l && s2.r+1 >= s.l {
			n := len(z)
			z[i] = z[n-1]
			z = z[:n-1]
			s.l = min(s.l, s2.l)
			s.r = max(s.r, s2.r)
		}
	}
	return append(z, s)
}

func solve(n int, ops []Op) []int {
	desk := make([]SRow, n+1)
	desk[0] = SRow{{1, n}}
	var buf SRow

	sort.Slice(ops, func(i, j int) bool {
		return ops[i].x < ops[j].x
	})

	top := 0
	for _, op := range ops {
		if debugEnable {
			log.Println("op:", op)
		}

		s := Segment{op.l, op.r}
		new_top := min(top+op.x, len(desk)-1)
		for i, i2 := new_top-op.x, new_top; i >= 0; i, i2 = i-1, i2-1 {

			buf = buf.And(desk[i], s)
			if debugEnable {
				log.Printf("%2d: and: %v", i, buf)
			}

			if len(buf) != 0 {
				if i2 > top {
					top = i2
				}
				for _, s := range buf {
					desk[i2] = desk[i2].Or(s)
				}
				if debugEnable {
					log.Printf("%2d: or: %v", i2, desk[i2])
				}
			}
		}
	}

	var res []int
	for i := 1; i <= top; i++ {
		row := desk[i]
		if len(row) != 0 {
			res = append(res, i)
		}
	}

	return res
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, q, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	ops := make([]Op, q)
	for i := range ops {
		l, r, x, err := scanThreeInt(sc)
		if err != nil {
			return err
		}
		ops[i] = Op{l, r, x}
	}

	maximums := solve(n, ops)

	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(bw, len(maximums), wo)
	writeInts(bw, maximums, wo)

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
