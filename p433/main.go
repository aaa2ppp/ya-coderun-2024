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

type warehouse struct {
	x     int
	price int
}

type order struct {
	x      int
	banned []int
}

func slowSolve(w []warehouse, o []order) []int {
	res := make([]int, len(o))

	for i := range o {
		banned := make(map[int]struct{}, len(o[i].banned))
		for _, v := range o[i].banned {
			banned[v-1] = struct{}{} // [v-1] to 0-indexing
		}
		minimum := math.MaxInt
		for j := range w {
			if _, ok := banned[j]; ok {
				if debugEnable {
					log.Println(i+1, j+1, "skip banned")
				}
				continue
			}
			price := w[j].price
			dist := o[i].x - w[j].x
			cost := price + dist*dist
			minimum = min(minimum, cost)
			if debugEnable {
				log.Println(i+1, j+1, "price:", price, "dist:", dist, "cost:", cost, "min:", minimum)
			}
		}
		res[i] = minimum
	}

	return res
}

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

func reverse[T any](aa []T) {
	for i, j := 0, len(aa)-1; i < j; i, j = i+1, j-1 {
		aa[i], aa[j] = aa[j], aa[i]
	}
}

func solve(w []warehouse, o []order) []int {
	iterate := iterate2

	res := make([]int, len(o))
	for i := range res {
		res[i] = math.MaxInt
	}

	// left -> right
	iterate(w, o, res)

	// reverse(w)
	for i := range w {
		w[i].x *= -1
	}

	// reverse(o)
	for i := range o {
		o[i].x *= -1
	}

	// left <- right
	iterate(w, o, res)

	return res
}

func intersect(w1, w2 warehouse, x int) fraction {
	return fraction{
		n: (w1.price + w1.x*w1.x) - (w2.price + w2.x*w2.x),
		m: 2 * (w1.x - w2.x),
	}
}

func cost(w warehouse, x int) int {
	dist := x - w.x
	return w.price + dist*dist
}

func iterate2(w []warehouse, o []order, res []int) {
	wx := make([]int, len(w))
	for i := range wx {
		wx[i] = i
	}
	sort.Slice(wx, func(i, j int) bool {
		i = wx[i]
		j = wx[j]
		return w[i].x < w[j].x
	})

	ox := make([]int, len(o))
	for i := range ox {
		ox[i] = i
	}
	sort.Slice(ox, func(i, j int) bool {
		i = ox[i]
		j = ox[j]
		return o[i].x < o[j].x
	})

	fix := func(k int, x int) int {
		for j := k; j > 0; j-- {
			xj := intersect(w[wx[j]], w[wx[j-1]], x)
			xf := fraction{x, 1}
			if !xf.less(xj) {
				return j
			}
			wx[j], wx[j-1] = wx[j-1], wx[j]
		}
		return 0
	}
	_ = fix

	i, l, r := 0, 0, 0
	for i < len(o) {
		_ = l
		if debugEnable {
			log.Println(i, r)
		}
		if r < len(wx) && w[wx[r]].x <= o[ox[i]].x { // TODO: дофига квадратных скобок
			r++
			continue
		}

		x := o[ox[i]].x

		if debugEnable {
			log.Println(wx[:r])
		}

		banned := make(map[int]struct{}, len(o[ox[i]].banned))
		for _, v := range o[ox[i]].banned {
			banned[v-1] = struct{}{} // [v-1] to 0-indexing
		}

		for j := l + 1; j < r; j++ {
			for jj := j - 1; jj >= 0; jj-- {
				if cost(w[wx[jj]], x) >= cost(w[wx[jj+1]], x) {
					break
				}
				wx[jj], wx[jj+1] = wx[jj+1], wx[jj]
			}
		}

		// sort.Slice(wx[l:r], func(i, j int) bool {
		// 	i += l
		// 	j += l
		// 	return cost(w[wx[i]], x) > cost(w[wx[j]], x)
		// })

		xf := fraction{x, 1}
		for j := l + 1; j < r; j++ {
			xj := intersect(w[wx[j]], w[wx[j-1]], x)
			if xf.less(xj) {
				l = j - 1
				break
			}
		}

		// for l := 0; l < r; l++ {
		// 	fix(l, x)
		// }

		for j := r - 1; j >= 0; j-- {
			// fix(j, x)
			if _, ok := banned[wx[j]]; !ok {
				res[ox[i]] = min(res[ox[i]], cost(w[wx[j]], x))
				break
			}
		}

		i++
	}
}

func iterate1(w []warehouse, o []order, res []int) {
	// TODO это можно вынести наружу из функции
	wx := make([]int, len(w))
	for i := range wx {
		wx[i] = i
	}
	sort.Slice(wx, func(i, j int) bool {
		i = wx[i]
		j = wx[j]
		return w[i].x < w[j].x
	})

	ox := make([]int, len(o))
	for i := range ox {
		ox[i] = i
	}
	sort.Slice(ox, func(i, j int) bool {
		i = ox[i]
		j = ox[j]
		return o[i].x < o[j].x
	})

	i := 0
	r := 0
	for i < len(ox) {
		if debugEnable {
			log.Println(i, r)
		}
		if r < len(wx) && w[wx[r]].x <= o[ox[i]].x { // TODO: дофига квадратных скобок
			r++
			continue
		}

		x := o[ox[i]].x

		if debugEnable {
			log.Println(wx[:r])
		}

		banned := make(map[int]struct{}, len(o[ox[i]].banned))
		for _, v := range o[ox[i]].banned {
			banned[v-1] = struct{}{} // [v-1] to 0-indexing
		}

		for j := r - 1; j >= 0; j-- {
			dist := x - w[wx[j]].x
			if dist*dist > res[ox[i]] {
				break
			}
			if _, ok := banned[wx[j]]; ok {
				continue
			}
			c := dist*dist + w[wx[j]].price
			if debugEnable {
				log.Println("cost:", res[ox[i]], c)
			}
			res[ox[i]] = min(res[ox[i]], c)
		}

		i++
	}
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	w := make([]warehouse, n)
	o := make([]order, m)
	for i := range w {
		s, p, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		w[i] = warehouse{s, p}
	}

	for i := range o {
		c, k, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		dd := make([]int, k)
		if err := scanInts(sc, dd); err != nil {
			panic(err)
		}
		o[i] = order{c, dd}
	}

	res := solve(w, o)

	writeInts(bw, res, writeOpts{sep: '\n', end: '\n'})
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
