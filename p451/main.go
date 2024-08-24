package main

import (
	"bufio"
	"container/heap"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type fraction struct {
	n, m int
}

func (f fraction) normalise() fraction {
	d := gcd(f.n, f.m)
	return fraction{f.n / d, f.m / d}
}

func (f fraction) less(other fraction) bool {
	return f.n*other.m < other.n*f.m
}

func (f fraction) asFloat() float64 {
	return float64(f.n) / float64(f.m)
}

func (f fraction) String() string {
	return strconv.FormatFloat(f.asFloat(), 'f', 6, 64)
}

type Item struct {
	l   int
	r   int
	avg fraction
}

type Heap []*Item

func (h Heap) Len() int { return len(h) }
func (h Heap) Less(i, j int) bool {
	return h[i].avg.less(h[j].avg) || !h[j].avg.less(h[i].avg) && (h[i].r-h[i].l) < (h[j].r-h[j].l)
}
func (h Heap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Item))
}

func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func slowSolve(k int, aa []int) float64 {
	n := len(aa)
	_ = n

	if debugEnable {
		log.Printf("k=%d aa=%d", k, aa)
	}

	sum := make([]int, len(aa)+1)
	for i := range aa {
		sum[i+1] = sum[i] + aa[i]
	}

	var h Heap

	buf := make([]int, 0, len(aa))
	res := fraction{0, 1}

	for d := k; d <= len(aa); d++ {
		buf = buf[:0]
		cur_max := 0
		cur_l := -1
		for i := 0; i+d <= len(aa); i++ {
			v := sum[i+d] - sum[i]
			if v > cur_max {
				cur_max = v
				cur_l = i
			}
			buf = append(buf, v)
		}
		cur_res := fraction{cur_max, d}
		if res.less(cur_res) {
			res = cur_res
			heap.Push(&h, &Item{cur_l, cur_l + d, cur_res})
		}
		// if debugEnable {
		// 	log.Printf("d=%3d: %3d, %.6f", d, buf, cur_res.asFloat())
		// }
	}

	if debugEnable {
		for h.Len() > 0 {
			it := heap.Pop(&h).(*Item)
			l, r, avg := it.l, it.r, it.avg
			log.Printf("[%d:%d](%d) %v %.6f", l, r, r-l, aa[l:r], avg.asFloat())

		}
	}
	return res.asFloat()
}

func solve(k int, aa []int) float64 {
	sum := make([]int, len(aa)+1)
	for i := range aa {
		sum[i+1] = sum[i] + aa[i]
	}

	res := fraction{sum[len(aa)] - sum[0], len(aa)}

	for len(aa) > k {
		if debugEnable {
			log.Printf("aa=%v %v", aa, res)
		}
		left := res
		left_i := len(aa)

		for i := 1; i+k <= len(aa); i++ {
			v := fraction{sum[i] - sum[0], i}
			if v.less(left) {
				left = v
				left_i = i
			}
		}
		if debugEnable {
			log.Printf("left: %d %v", left_i, left)
		}

		right := res
		right_i := 0

		for i := len(aa); i >= 0+k; i-- {
			v := fraction{sum[len(aa)] - sum[i], len(aa) - i}
			if v.less(right) {
				right = v
				right_i = i
			}
		}
		if debugEnable {
			log.Printf("right: %d %v", right_i, right)
		}

		if !left.less(res) && !right.less(res) {
			break
		}

		if left.less(right) {
			aa = aa[left_i:]
			sum = sum[left_i:]
		} else {
			aa = aa[:right_i]
			sum = sum[:right_i+1]
		}

		res = fraction{sum[len(aa)] - sum[0], len(aa)}
	}

	if debugEnable {
		log.Printf("aa=%v %v", aa, res)
	}

	return res.asFloat()
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, k, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	aa := make([]int, n)
	if err := scanInts(sc, aa); err != nil {
		panic(err)
	}

	res := solve(k, aa)
	bw.WriteString(strconv.FormatFloat(res, 'f', 6, 64))
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
