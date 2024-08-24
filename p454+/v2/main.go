package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func solve(n, m int, ab [][2]int) ([]int, bool) {

	people := make([][][2]int, n+1)
	toppings := make([][2][]int, m+1)

	for i := range ab {
		p, t := ab[i][0], ab[i][1]
		if debugEnable {
			log.Println(p, t)
		}

		var want int
		if t > 0 {
			want = 1
		} else {
			want = 0
			t = -t
		}

		toppings[t][want] = append(toppings[t][want], p)
		people[p] = append(people[p], [2]int{t, want})
	}

	if debugEnable {
		log.Println("toppings:", toppings)
		log.Println("people:", people)
	}

	// мрак и тихий ужас

	var (
		visited = make([]bool, len(toppings))
		nodes   = make([]int, len(toppings))
		states  = make([]int, len(toppings))
	)

	var dfs func(node int, nodes []int) bool

	dfs = func(node int, nodes []int) bool {
		var dfs2 func(node int, want int) bool

		dfs2 = func(node int, want int) bool {
			if visited[node] {
				return states[node] == want
			}

			visited[node] = true
			nodes = append(nodes, node)
			states[node] = want

			if debugEnable {
				log.Println("set node", node, want)
			}

			for _, i := range toppings[node][(want+1)&1] {
				for _, p := range people[i] {
					neig, want := p[0], p[1]

					if neig == node {
						continue
					}

					if !dfs2(neig, want) {
						return false
					}
				}
			}

			return true
		}

		for node < len(toppings) && visited[node] {
			node++
		}
		if node == len(toppings) {
			return true
		}

		if debugEnable {
			log.Println("try node", node, 0)
		}
		if dfs2(node, 0) && dfs(node+1, nodes[len(nodes):]) {
			return true
		}
		for _, v := range nodes {
			visited[v] = false
		}
		nodes = nodes[:0]

		if debugEnable {
			log.Println("try node", node, 1)
		}
		if dfs2(node, 1) && dfs(node+1, nodes[len(nodes):]) {
			return true
		}
		for _, v := range nodes {
			visited[v] = false
		}

		return false
	}

	dfs(1, nodes[:0])

	count := 0
	for i := range states {
		count += states[i]
	}

	res := make([]int, 0, count)
	for i := 1; i < len(states); i++ {
		if states[i] == 1 {
			res = append(res, i)
		}
	}

	return res, true
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, k, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	ab := make([][2]int, 0, k)
	for i := 0; i < k; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		ab = append(ab, [2]int{a, b})
	}

	res, ok := solve(n, m, ab)

	if !ok {
		bw.WriteString("-1\n")
	} else {
		wo := writeOpts{sep: ' ', end: '\n'}
		writeInt(bw, len(res), wo)
		writeInts(bw, res, wo)
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

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}