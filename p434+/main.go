package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"unsafe"
)

type solver struct {
	bw  *bufio.Writer
	buf []byte
	cnt []int
}

func newSolver(bw *bufio.Writer, n int) *solver {
	return &solver{
		bw:  bw,
		buf: make([]byte, n*5),
		cnt: make([]int, n),
	}
}

func (s *solver) solve(i, nn int) {
	n := len(s.cnt)
	v := n - i
	if v == 1 {
		s.cnt[i] = nn
		s.output()
		return
	}
	for j := 0; j <= nn/v; j++ {
		s.cnt[i] = j
		s.solve(i+1, nn-v*j)
	}
}

func (s *solver) output() {
	s.buf = s.buf[:0]
	n := len(s.cnt)
	for i, v := 0, n; i < n; i, v = i+1, v-1 {
		for j := 0; j < s.cnt[i]; j++ {
			s.buf = strconv.AppendInt(s.buf, int64(v), 10)
			s.buf = append(s.buf, " + "...)
		}
	}
	s.buf = s.buf[:len(s.buf)-3]
	s.buf = append(s.buf, '\n')
	s.bw.Write(s.buf)
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

	newSolver(bw, n).solve(0, n)
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
