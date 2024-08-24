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

func solve(w, h int, bb [][3]int, qq [][2]int) [][]int {
	var d [6]int
	max_d := min(w, h)

	d[0] = sort.Search(max_d, func(d int) bool {
		return checkReach(w, h, bb, 3, 0, d)
	})
	d[1] = sort.Search(max_d, func(d int) bool {
		return checkReach(w, h, bb, 0, 1, d)
	})
	d[2] = sort.Search(max_d, func(d int) bool {
		return checkReach(w, h, bb, 1, 2, d)
	})
	d[3] = sort.Search(max_d, func(d int) bool {
		return checkReach(w, h, bb, 2, 3, d)
	})
	d[4] = sort.Search(max_d, func(d int) bool {
		return checkReach(w, h, bb, 0, 2, d)
	})
	d[5] = sort.Search(max_d, func(d int) bool {
		return checkReach(w, h, bb, 1, 3, d)
	})

	res := make([][]int, 0, len(qq))
	for i := range qq {
		r, s := qq[i][0], qq[i][1]
		r2 := r * 2

		corners := makeMatrix[bool](5, 5)
		if r2 < d[0] && r2 < d[1] && r2 < d[4] {
			corners[0][1] = true
			corners[1][0] = true
		}
		if r2 < d[0] && r2 < d[2] && r2 < d[4] && r2 < d[5] {
			corners[0][2] = true
			corners[2][0] = true
		}
		if r2 < d[0] && r2 < d[3] && r2 < d[5] {
			corners[0][3] = true
			corners[3][0] = true
		}
		if r2 < d[1] && r2 < d[2] && r2 < d[5] {
			corners[1][2] = true
			corners[2][1] = true
		}
		if r2 < d[1] && r2 < d[3] && r2 < d[4] && r2 < d[5] {
			corners[1][3] = true
			corners[3][1] = true
		}
		if r2 < d[2] && r2 < d[3] && r2 < d[4] {
			corners[2][3] = true
			corners[3][2] = true
		}

		res = append(res, getReach(corners, s))
	}

	return res
}

func getReach(matrix [][]bool, s int) []int {
	frontier := newQueue(len(matrix))
	visited := make([]bool, len(matrix))

	frontier.push(s)
	visited[s] = true

	for frontier.len() > 0 {
		node := frontier.pop()

		for neig := range matrix {
			if !visited[neig] && matrix[node][neig] {
				visited[neig] = true
				frontier.push(neig)
			}
		}
	}

	res := make([]int, 0, len(matrix))
	for i := range visited {
		if visited[i] {
			res = append(res, i)
		}
	}

	return res
}

func checkReach(w, h int, bb [][3]int, a, b int, d int) bool {
	frontier_a := newQueue(len(bb))
	frontier_b := newQueue(len(bb))
	visited_a := make([]bool, len(bb))
	visited_b := make([]bool, len(bb))

	for i := 0; i < len(bb); i++ {
		xi, yi, ri := bb[i][0], bb[i][1], bb[i][2]

		if yi-ri < d {
			if a == 0 {
				frontier_a.push(i)
			}
			if b == 0 {
				frontier_b.push(i)
			}
		}

		if w-(xi+ri) < d {
			if a == 1 {
				frontier_a.push(i)
			}
			if b == 1 {
				frontier_b.push(i)
			}
		}

		if h-(yi+ri) < d {
			if a == 2 {
				frontier_a.push(i)
			}
			if b == 2 {
				frontier_b.push(i)
			}
		}

		if xi-ri < d {
			if a == 3 {
				frontier_a.push(i)
			}
			if b == 3 {
				frontier_b.push(i)
			}
		}
	}

	for frontier_a.len() > 0 && frontier_b.len() > 0 {
		node_a := frontier_a.pop()
		if visited_b[node_a] {
			return true
		}
		node_b := frontier_b.pop()
		if visited_a[node_b] {
			return true
		}

		for neig := range bb {
			i, j := node_a, neig
			xi, yi, ri := bb[i][0], bb[i][1], bb[i][2]
			xj, yj, rj := bb[j][0], bb[j][1], bb[j][2]
			if !visited_a[neig] && pow2(xi-xj)+pow2(yi-yj) < pow2(ri+rj+d) {
				visited_a[neig] = true
				frontier_a.push(neig)
			}
		}

		for neig := range bb {
			i, j := node_b, neig
			xi, yi, ri := bb[i][0], bb[i][1], bb[i][2]
			xj, yj, rj := bb[j][0], bb[j][1], bb[j][2]

			if !visited_b[neig] && pow2(xi-xj)+pow2(yi-yj) < pow2(ri+rj+d) {
				visited_b[neig] = true
				frontier_b.push(neig)
			}
		}
	}

	return false
}

// TODO queue
type queue []int

func newQueue(n int) *queue {
	q := make(queue, 0, n)
	return &q
}

func (q queue) len() int {
	return len(q)
}

func (q *queue) push(v int) {
	*q = append(*q, v)
}

func (q *queue) pop() int {
	old := *q
	v := old[0]
	*q = old[1:]
	return v
}

func pow2(a int) int {
	return a * a
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	w, h, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	bb := make([][3]int, n)
	for i := range bb {
		x, y, r, err := scanThreeInt(sc)
		if err != nil {
			return err
		}
		bb[i] = [3]int{x, y, r}
	}

	qq := make([][2]int, m)
	for i := range qq {
		r, s, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		qq[i] = [2]int{r, s - 1} // to 0-indexed s
	}

	res := solve(w, h, bb, qq)

	buf := make([]byte, 0, 16)
	for i := range res {
		buf = buf[:0]
		for _, v := range res[i] {
			buf = strconv.AppendInt(buf, int64(v+1), 10) // to 1-indexed
		}
		buf = append(buf, '\n')
		bw.Write(buf)
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