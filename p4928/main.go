package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

func factorize(n int) map[int]int {
	b := map[int]int{}
	for n%2 == 0 {
		b[2]++
		n /= 2
	}
	for n%3 == 0 {
		b[3]++
		n /= 3
	}
	i := 5
	inc := 2
	for i*i <= n {
		for n%i == 0 {
			b[i]++
			n /= i
		}
		i += inc // 7(4), 9(2), 11(2), 13(4), 17(2), 19(4), 21(2), 23(4)
		inc = 6 - inc
	}
	if n > 1 {
		b[i]++
	}
	return b
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, q, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	scanAA := func(aa []map[int]int) {
		for i := range aa {
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			f := factorize(v)
			if i > 0 {
				for k, v := range aa[i-1] {
					f[k] += v
				}
			}
			aa[i] = f
		}
	}

	aa := make([]map[int]int, n+1)
	scanAA(aa)
	if debugEnable {
		log.Printf("aa: %v", aa)
	}

	bb := make([]map[int]int, m+1)
	scanAA(bb)
	if debugEnable {
		log.Printf("bb: %v", bb)
	}

	scanCA := func(ca []int) {
		for i := range ca {
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			if i == 0 {
				ca[0] = v
			} else {
				ca[i] = min(ca[i-1], v)
			}
		}
	}

	ca := make([]int, n+1)
	scanCA(ca)
	if debugEnable {
		log.Printf("bb: %v", ca)
	}

	cb := make([]int, m+1)
	scanCA(cb)
	if debugEnable {
		log.Printf("bb: %v", cb)
	}

	for i := 0; i < q; i++ {
		k, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		xx := make([]int, k)
		if err := scanInts(sc, xx); err != nil {
			panic(err)
		}

		minmax_aa_i := math.MaxInt
		minmax_bb_i := math.MaxInt

		for _, x := range xx {
			f := factorize(x)
			if debugEnable {
				log.Printf("xf: %v", f)
			}

			get_max_i := func(aa []map[int]int) int {
				max_aa_i := 0
				for k, v := range f {
					i := sort.Search(len(aa), func(i int) bool {
						return aa[i][k] >= v
					})
					max_aa_i = max(max_aa_i, i)
				}
				return max_aa_i
			}

			minmax_aa_i = min(minmax_aa_i, get_max_i(aa))
			minmax_bb_i = min(minmax_bb_i, get_max_i(bb))
		}

		res := min(
			ca[minmax_aa_i]+cb[minmax_bb_i-1],
			ca[minmax_aa_i-1]+cb[minmax_bb_i],
		)

		writeInt(bw, res, writeOpts{sep: '\n', end: '\n'})
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
