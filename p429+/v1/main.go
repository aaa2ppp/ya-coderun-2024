package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func solve(w, h int, b [][3]int, q [][2]int) [][]int {
	res := make([][]int, 0, len(q))

	for i := range q {
		r, s := q[i][0], q[i][1]
		matrix := calcDistances(w, h, b, r*2)
		if debugEnable {
			log.Println("matrix:")
			for _, row := range matrix {
				log.Println(toBin(row))
			}
		}

		var d [6]bool
		d[0] = checkReach(matrix, 0, 3)
		d[1] = checkReach(matrix, 0, 1)
		d[2] = checkReach(matrix, 1, 2)
		d[3] = checkReach(matrix, 2, 3)
		d[4] = checkReach(matrix, 0, 2)
		d[5] = checkReach(matrix, 1, 3)
		if debugEnable {
			log.Println("d:", toBin(d[:]))
		}

		corners := makeMatrix[bool](5, 5)
		if !d[0] && !d[1] && !d[4] {
			corners[0][1] = true
			corners[1][0] = true
		}
		if !d[0] && !d[2] && !d[4] && !d[5] {
			corners[0][2] = true
			corners[2][0] = true
		}
		if !d[0] && !d[3] && !d[5] {
			corners[0][3] = true
			corners[3][0] = true
		}
		if !d[1] && !d[2] && !d[5] {
			corners[1][2] = true
			corners[2][1] = true
		}
		if !d[1] && !d[3] && !d[4] && !d[5] {
			corners[1][3] = true
			corners[3][1] = true
		}
		if !d[2] && !d[3] && !d[4] {
			corners[2][3] = true
			corners[3][2] = true
		}

		res = append(res, getReach(corners, s))
	}

	return res
}

func toBin(a []bool) []int {
	b := make([]int, 0, len(a))
	for _, v := range a {
		if v {
			b = append(b, 1)
		} else {
			b = append(b, 0)
		}
	}
	return b
}

func getReach(matrix [][]bool, s int) []int {
	frontier := newQueue(len(matrix))
	visited := make([]bool, len(matrix))

	frontier.push(s)
	visited[s] = true

	for frontier.len() > 0 {
		node := frontier.pop()

		for neig := range matrix[node] {
			if matrix[node][neig] && !visited[neig] {
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

func checkReach(matrix [][]bool, a, b int) bool {
	frontier := newQueue(len(matrix))
	visited := make([]bool, len(matrix))

	frontier.push(a)
	visited[a] = true

	for frontier.len() > 0 {
		node := frontier.pop()
		if node == b {
			return true
		}

		for neig := range matrix[node] {
			if matrix[node][neig] && !visited[neig] {
				visited[neig] = true
				frontier.push(neig)
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

func calcDistances(w, h int, b [][3]int, d int) [][]bool {
	n := len(b)
	matrix := makeMatrix[bool](n+4, n+4)

	for i := 0; i < n; i++ {
		xi, yi, ri := b[i][0], b[i][1], b[i][2]

		if yi-ri < d {
			matrix[0][i+4] = true
			matrix[i+4][0] = true
		}

		if w-(xi+ri) < d {
			matrix[1][i+4] = true
			matrix[i+4][1] = true
		}

		if h-(yi+ri) < d {
			matrix[2][i+4] = true
			matrix[i+4][2] = true
		}

		if xi-ri < d {
			matrix[3][i+4] = true
			matrix[i+4][3] = true
		}
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			xi, yi, ri := b[i][0], b[i][1], b[i][2]
			xj, yj, rj := b[j][0], b[j][1], b[j][2]

			if pow2(xi-xj)+pow2(yi-yj) < pow2(ri+rj+d) {
				matrix[i+4][j+4] = true
				matrix[j+4][i+4] = true
			}
		}
	}

	return matrix
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

	b := make([][3]int, n)
	for i := range b {
		x, y, r, err := scanThreeInt(sc)
		if err != nil {
			return err
		}
		b[i] = [3]int{x, y, r}
	}

	q := make([][2]int, m)
	for i := range q {
		r, s, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		q[i] = [2]int{r, s - 1} // to 0-indexed s
	}

	res := solve(w, h, b, q)

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
