package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
	"unsafe"
)

func run(in io.Reader, out io.Writer) error {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	t, err := br.ReadBytes('\n')
	if err != nil {
		return err
	}
	t = bytes.TrimRightFunc(t, unicode.IsSpace)

	s, err := br.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return err
	}
	s = bytes.TrimRightFunc(s, unicode.IsSpace)

	res := solve(t, s)

	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(bw, len(res), wo)
	writeInts(bw, res, wo)

	return nil
}

const X = 269

func solve(t, s []byte) []int {
	var res []int

	pow := make([]uint64, len(s)+1)
	pow[0] = 1
	for i := 1; i < len(pow); i++ {
		pow[i] = pow[i-1] * X
	}

	s_pow := makePattern(s, pow)
	if debugEnable {
		log.Println("s_pow:", s_pow)
	}

	// TODO: можно ужаться по пямяти, если использовать кольцевую очередь размера len(s)+1
	t_hashe := make([]uint64, len(t)+1)
	for i := 1; i < len(t_hashe); i++ {
		t_hashe[i] = t_hashe[i-1]*X + uint64(t[i-1])
	}
	if debugEnable {
		log.Println("t_hash:", t_hashe)
	}

main_loop:
	for i := 0; i <= len(t)-len(s); i++ {
		var set ASCIISet
		var s_hash uint64
		for _, p := range s_pow {
			if set.Contains(t[i+p.pos]) {
				// oops!..
				continue main_loop
			}
			set.Add(t[i+p.pos])
			s_hash += uint64(t[i+p.pos]) * p.pow
		}

		t_h1 := t_hashe[i] * pow[len(s)]
		t_h2 := t_hashe[i+len(s)]

		if t_h2 == s_hash+t_h1 {
			// bingo!
			res = append(res, i+1) // +1 to 1-indexed result
		}
	}

	return res
}

type patternItem struct {
	pos int
	pow uint64
}

func makePattern(s []byte, pow []uint64) []patternItem {
	s_pow_set := make([]patternItem, 128)
	count := 0
	var set ASCIISet
	for i, c := range s {
		if !set.Contains(c) {
			set.Add(c)
			s_pow_set[c] = patternItem{i, 0}
			count++
		}
		s_pow_set[c].pow += pow[len(s)-i-1]
	}

	s_pow := make([]patternItem, 0, count)
	for i := range s_pow_set {
		if set.Contains(byte(i)) {
			s_pow = append(s_pow, s_pow_set[i])
		}
	}

	return s_pow
}

type ASCIISet struct {
	two_word [2]uint64
}

func (s *ASCIISet) Add(c byte) {
	// any ASCII <= 127, max bit 7
	s.two_word[(c>>6)&1] |= 1 << (c & 63)
}

func (s *ASCIISet) Contains(c byte) bool {
	return s.two_word[(c>>6)&1]&(1<<(c&63)) != 0
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
