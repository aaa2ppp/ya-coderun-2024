package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unsafe"
)

const (
	infinity = math.MaxInt
)

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	matrix := makeMatrix(n, n)
	for i := range matrix {
		if err := scanInts(sc, matrix[i]); err != nil {
			return err
		}
	}

	res := solve(matrix)
	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(bw, res, wo)
	_ = n

	return nil
}

func makeMatrix(n, m int) [][]int {
	buf := make([]int, n*m)
	matrix := make([][]int, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func solve(matrix [][]int) int {
	n := len(matrix)

	distances := make([]int, n)
	previous := make([]int, n)
	calcDistances(matrix, distances, previous)
	if debugEnable {
		log.Println("distances:", distances)
		log.Println("previous :", previous)
	}

	// Восстанавливаем дерево путей
	path_tree := make([][]int, len(previous))
	for b, a := range previous {
		if a == -1 {
			continue
		}
		path_tree[a] = append(path_tree[a], b)
	}

	order := topologicalSort(path_tree)

	// Cчитаем сколько раз вошли в узел (вес входящего ребра)
	frequency := make([]int, n)
	for _, i := range order {
		frequency[i] = 1
	}
	for _, i := range order {
		if prev := previous[i]; prev != -1 {
			frequency[prev] += frequency[i]
		}
	}
	if debugEnable {
		log.Println("frequency:", frequency)
	}

	// Ищем ребро вносящее наибольшее изменение начиная самого тяжелого.
	// Прекращаем поиск на ребре, чей вес меньше или равен (которое не может вносить больше изменений),
	// чем текущий результат.

	passed := make([]bool, n)
	passed[0] = true

	max_frequency := func() (int, int) {
		maximum := -1
		max_idx := -1
		for i, v := range frequency {
			if passed[i] {
				continue
			}
			if v > maximum {
				maximum = v
				max_idx = i
			}
		}
		return maximum, max_idx
	}

	res := 0
	distances2 := make([]int, n)
	previous2 := make([]int, n) // stub

	for {
		freq, idx := max_frequency()
		if idx == -1 || freq <= res {
			break
		}
		passed[idx] = true

		a, b := idx, previous[idx]
		bak := matrix[a][b]

		matrix[a][b] = -1
		matrix[b][a] = -1

		calcDistances(matrix, distances2, previous2)
		diff := countDiffs(distances, distances2)
		res = max(res, diff)

		matrix[a][b] = bak
		matrix[b][a] = bak
	}

	return res
}

func countDiffs(a, b []int) int {
	count := 0
	for i := range a {
		if a[i] != b[i] {
			count++
		}
	}
	return count
}

func topologicalSort(paths [][]int) []int {
	order := make([]int, len(paths))

	var dfs func(node int) 
	dfs = func(node int) {
		for _, neig := range paths[node] {
			dfs(neig)
		}
		order = append(order, node)
	}

	dfs(0)
	return order
}

func calcDistances(matrix [][]int, distances, previous []int) {
	n := len(matrix)

	for i := range distances {
		distances[i] = infinity
	}
	for i := range previous {
		previous[i] = -1
	}

	passed := make([]bool, n)

	get_node := func() int {
		min_dist := infinity
		min_node := -1
		for node, dist := range distances {
			if passed[node] {
				continue
			}
			if dist < min_dist {
				min_dist = dist
				min_node = node
			}
		}
		return min_node
	}

	distances[0] = 0
	for {
		node := get_node()
		if node == -1 {
			break
		}
		cur_dist := distances[node]
		if cur_dist == infinity {
			break
		}
		passed[node] = true

		for neig, dist := range matrix[node] {
			if dist == -1 || neig == node || passed[neig] {
				continue
			}
			if next_dist := cur_dist + dist; next_dist < distances[neig] {
				distances[neig] = next_dist
				previous[neig] = node
			}
		}
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

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
