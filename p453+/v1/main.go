package main_v1

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	// scan edges
	edges := make([][2]int, m)
	for i := range edges {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		edges[i] = newEdge(a, b)
	}
	if debugEnable {
		log.Println("edges:", edges)
	}

	res := solve(n, edges)
	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(bw, res, wo)

	return nil
}

func newEdge(a, b int) [2]int {
	if a > b {
		return [2]int{b, a}
	}
	return [2]int{a, b}
}

func solve(n int, edges [][2]int) int {
	graph := makeGraph(n, edges)
	if debugEnable {
		log.Println("graph:", graph)
	}

	loop_edges := findLoopEdges(graph)
	if debugEnable {
		log.Println("loop_edges:", loop_edges)
	}

	graph_without_loop := makeGraphExclude(n, edges, loop_edges)
	if debugEnable {
		log.Println("graph_without_loop:", graph_without_loop)
	}

	// резделим граф на компоненты смежности (размером >= 3, т.к. для компонент <=2 невозможно создать цикл добавляя ребра)
	components := devideByComponents(graph_without_loop)
	if debugEnable {
		log.Println("components:", components)
	}

	// обратим внимание, что каждая компонента так же не содержит циклов, и не воможно создатьтолько цикл соединя компоненты
	// для каждой компоненты считаем количество ребер которые можно добавить для возникновения цикла
	total := 0
	for _, comp := range components {
		node := comp[0]
		size := comp[1]
		count := countEdgesToLooping(graph_without_loop, node, size)
		if debugEnable {
			log.Printf("count edges for [%d, %d]: %d", node, size, count)
		}
		total += count
	}

	return total
}

func makeGraph(n int, edges [][2]int) [][]int {
	graph := make([][]int, n+1)
	for _, edge := range edges {
		a, b := edge[0], edge[1]
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}
	return graph
}

func makeGraphExclude(n int, edges, exclude [][2]int) [][]int {
	exclude_set := make(map[[2]int]struct{}, len(exclude))
	for _, edge := range exclude {
		exclude_set[edge] = struct{}{}
	}
	graph := make([][]int, n+1)
	for _, edge := range edges {
		if _, ok := exclude_set[edge]; ok {
			continue
		}
		a, b := edge[0], edge[1]
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}
	return graph
}

// находит все ребра графа которые участвуют в чиклах
func findLoopEdges(graph [][]int) [][2]int {
	const (
		white = 0
		grey  = 1
		black = 2
	)

	var edges [][2]int
	visited := make([]byte, len(graph))

	var dfs func(node, prev int) map[int]struct{}

	dfs = func(node, prev int) (set map[int]struct{}) {
		if visited[node] == black {
			return nil
		}
		if visited[node] == grey {
			return map[int]struct{}{node: {}}
		}
		visited[node] = grey

		for _, neig := range graph[node] {
			if neig == prev {
				continue
			}

			if neig_set := dfs(neig, node); len(neig_set) > 0 {
				edges = append(edges, newEdge(node, neig))

				if set == nil {
					set = neig_set
				} else {
					// O(N)!
					for k := range neig_set {
						set[k] = struct{}{}
					}
				}
			}
		}

		visited[node] = black
		delete(set, node)
		return set
	}

	for node := range graph {
		if visited[node] == white {
			dfs(node, -1)
		}
	}

	return edges
}

// возвращает cписок [2]int{node, size}, где node один из узов компоненты,
// а size - количество узлов в компоненте, для size >= 3
func devideByComponents(graph [][]int) [][2]int {
	visited := make([]bool, len(graph))

	var dfs1 func(node int) int
	dfs1 = func(node int) int {
		visited[node] = true
		size := 1
		for _, neig := range graph[node] {
			if !visited[neig] {
				size += dfs1(neig)
			}
		}
		return size
	}

	var comps [][2]int
	for node := range graph {
		if !visited[node] {
			size := dfs1(node)
			// игнорируем компоненты, для которых гарантированно не возможно создать цикл
			if size >= 3 {
				comps = append(comps, [2]int{node, size})
			}
		}
	}

	return comps
}

// подсчитывает количество способов добавления ребера в компоненту для возникновения цикла
func countEdgesToLooping(graph [][]int, node int, size int) int {
	var dfs2 func(node, prev int) int
	dfs2 = func(node, prev int) int {
		count := size - len(graph[node]) - 1
		for _, neig := range graph[node] {
			if neig != prev {
				count += dfs2(neig, node)
			}
		}
		return count
	}
	return dfs2(node, -1) / 2
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
