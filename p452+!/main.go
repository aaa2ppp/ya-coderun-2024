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

	bb_pos := make([][2]int, n*m+1)
	exists_zero := false
	for i, row := range bb {
		for j, v := range row {
			if v == 0 {
				exists_zero = true
			} else {
				bb_pos[v] = [2]int{i, j}
			}
		}
	}

	count := 0
	row_count := make([][3]int, n) // считаем количество дырок - 0 (0), входящих(1) и исходящих(2) ребер (относительно движения дырок)
	graph := make([][]int, n)      // неорентированный граф для выделения компонент
	var edges [][2]int
	if debugEnable {
		edges = make([][2]int, 0, n)
	}

	buf := make([]int, m) // буфер для перемещения внутри строки by LIS
	for i, row := range aa {
		buf = buf[:0]
		for _, v := range row {
			if v == 0 {
				row_count[i][0]++
				continue
			}

			if a, b := i, bb_pos[v][0]; a != b {
				// перемещения между строками
				count++

				// ноль (дырка) движется в противоположную сторону перемещению из b в a
				row_count[a][1]++ // входящие ребра
				row_count[b][2]++ // исходящие ребра

				if debugEnable {
					edges = append(edges, [2]int{b, a})
				}
				graph[a] = append(graph[a], b)
				graph[b] = append(graph[b], a)
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
			if len(buf) == m {
				// одно перемешение будет через другую строку (нужно будет "одолжить" ноль)
				need_move++
			}
			count += need_move
		}
	}

	for _, edge := range edges {
		b, a := edge[0], edge[1]
		if debugEnable {
			log.Println(b, "->", a, row_count[b], row_count[a])
		}

	}

	// обходим граф добавляем дифицит дырок и +1 на каждую компонету, если у нее нет источника дырки 
	visited := make([]bool, n)
	var dfs func(node int) (int, int, bool)

	dfs = func(node int) (int, int, bool) { // недостаток 0, в компоненте есть источник 0
		if visited[node] {
			return 0, 0, false
		}
		visited[node] = true
		count := 1
		exists_src := row_count[node][0] > 0 && row_count[node][2] > 0
		shortage := 0
		if n := row_count[node][0] + row_count[node][1] - row_count[node][2]; n < 0 {
			shortage = -n
		}
		for _, neig := range graph[node] {
			co, sh, ex := dfs(neig)
			count += co
			shortage += sh
			exists_src = exists_src || ex
		}
		return count, shortage, exists_src
	}

	for node := range graph {
		if !visited[node] {
			node_count, shortage, exists_src := dfs(node)
			if debugEnable {
				log.Println(node_count, shortage, exists_src)
			}
			if node_count == 1 {
				continue
			}
			count += shortage
			_ = exists_src
			if !exists_src {
				count++
			}
		}
	}

	// это можно проверить раньше
	if count > 0 && !exists_zero {
		return -1
	}

	return count
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
