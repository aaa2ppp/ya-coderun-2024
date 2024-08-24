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

func solve(aa, bb [][]int) int {
	n := len(aa)
	m := len(aa[0])
	count := 0

	bb_pos := make([][2]int, n*m+1)
	bb_rows_contains_zero := make([]bool, n+1)
	bb_row_contains_zero := -1
	exists_zero := false
	for i, row := range bb {
		for j, v := range row {
			if v == 0 {
				bb_rows_contains_zero[i] = true
				bb_row_contains_zero = i
				exists_zero = true
			} else {
				bb_pos[v] = [2]int{i, j}
			}
		}
	}

	buf := make([]int, m) // буфер для перемещения внутри строки
	// aa_from_to_row1 := newFromTo(n + 1)
	aa_from_to_row2 := newFromTo(n + 1)
	aa_from_to_row3 := newFromTo(n + 1)
	aa_row_to_count := make([]int, n+1)
	aa_row_zero_count := make([]int, n+1)
	for i, row := range aa {
		buf = buf[:0]
		for _, v := range row {
			if v == 0 {
				aa_row_zero_count[i]++
				continue

			}
			if to := bb_pos[v][0]; to != i {
				// значение нужно переместить на другую строку
				if bb_rows_contains_zero[i] {
					aa_from_to_row3.push(i, to)
				} else {
					aa_from_to_row2.push(i, to)
				}
				aa_row_to_count[to]++
				continue

			}

			// Значение в нужной строке. Запоминим требуемую позицию (возможно нужно будет переместить).
			buf = append(buf, bb_pos[v][1])
		}

		// Перемещаем внутри строки
		// Оставляем на месте наибольшую возрастающую последовательность (по требуемой позиции).
		// Другие значения будем перемещать. Это требует наименьшего количества перемещений.
		if debugEnable {
			log.Println("buf:", buf)
		}
		if len(buf) == 0 {
			continue
		}
		need_move := len(buf) - findLISLength(buf)
		if need_move > 0 {
			if !exists_zero {
				// oops!.. для перемещения нужен хотябы один ноль
				return -1
			}
			if len(buf) == m {
				// одно перемешение будет через другую строку (нужно будет "одолжить" ноль)
				need_move++
			}
			count += need_move
		}
	}

	ready_for_reception := make(map[int]struct{}, n+1)
	for i := range aa_row_to_count {
		if aa_row_to_count[i] > 0 && aa_row_zero_count[i] > 0 {
			ready_for_reception[i] = struct{}{}
		}
	}

	// Будем перемещать между строк
	for aa_from_to_row2.len() > 0 || aa_from_to_row3.len() > 0 {
		if debugEnable {
			log.Println("from_to:", aa_from_to_row2)
			log.Println("from_to:", aa_from_to_row3)
		}
		if !exists_zero {
			// oops!.. для перемещения нужен хотябы один ноль
			return -1
		}

		to := -1
		for k := range ready_for_reception {
			to = k
			break
		}

		if to != -1 {
			from := aa_from_to_row2.popFrom(to)
			if from == -1 {
				from = aa_from_to_row3.popFrom(to)
			}
			if debugEnable {
				log.Printf("%d -> %d", from, to)
			}
	
			aa_row_zero_count[to]--
			aa_row_to_count[to]--
			if !(aa_row_to_count[to] > 0 && aa_row_zero_count[to] > 0) {
				delete(ready_for_reception, to)
			}

			aa_row_zero_count[from]++
			if aa_row_to_count[from] != 0 {
				ready_for_reception[from] = struct{}{}
			}
		} else {
			// need comments
			from, to := aa_from_to_row2.popAny()
			if from == -1 {
				from, to = aa_from_to_row3.popAny()
			}
			if debugEnable {
				log.Printf("%d -> %d", from, bb_row_contains_zero)
			}
			aa_from_to_row3.push(bb_row_contains_zero, to)
			aa_row_zero_count[bb_row_contains_zero]--
			aa_row_zero_count[from]++
			ready_for_reception[from] = struct{}{}
		}

		count++
	}

	return count
}

type fromTo struct {
	count int
	items map[int]stack
}

func newFromTo(n int) *fromTo {
	items := make(map[int]stack, n)
	return &fromTo{items: items}
}

func (p *fromTo) len() int {
	return p.count
}

func (p *fromTo) push(from, to int) {
	froms := p.items[to]
	froms.push(from)
	p.items[to] = froms
	p.count++
}

func (p *fromTo) popFrom(to int) int {
	if p.count == 0 {
		return -1
	}
	froms := p.items[to]
	if len(froms) == 0 {
		return -1
	}
	from := froms.pop()
	if len(froms) == 0 {
		delete(p.items, to)
	} else {
		p.items[to] = froms
	}
	p.count--
	return from
}

func (p *fromTo) popAny() (int, int) {
	from, to := -1, -1
	for k, v := range p.items {
		to = k
		from = v.pop()
		if len(v) == 0 {
			delete(p.items, to)
		} else {
			p.items[to] = v
		}
		p.count--
		break
	}
	return from, to
}

type stack []int

func (s *stack) push(v int) {
	*s = append(*s, v)
}

func (s *stack) pop() int {
	old := *s
	n := len(old)
	v := old[n-1]
	*s = old[:n-1]
	return v
}

func findLISLength(a []int) int {
	n := len(a)
	dp := make([]int, n+1)
	length := 0

	dp[0] = -math.MaxInt
	for i := 1; i < len(dp); i++ {
		dp[i] = math.MaxInt
	}

	for i := 0; i < n; i++ {
		j := sort.Search(n, func(j int) bool {
			return dp[j] >= a[i]
		})
		if dp[j-1] < a[i] && a[i] < dp[j] {
			dp[j] = a[i]
			length = max(length, j)
		}
	}

	return length
}

func scanIntMatrix(sc *bufio.Scanner, matrix [][]int) error {
	for _, row := range matrix {
		if err := scanInts(sc, row); err != nil {
			return err
		}
	}
	return nil
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

	aa := makeMatrix[int](n, m)
	if err := scanIntMatrix(sc, aa); err != nil {
		return err
	}

	bb := makeMatrix[int](n, m)
	if err := scanIntMatrix(sc, bb); err != nil {
		return err
	}

	writeInt(bw, solve(aa, bb), writeOpts{end: '\n'})
	return nil
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
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
