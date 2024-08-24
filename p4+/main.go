package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"math/big"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

func solve(t, s int, v []int) *big.Int {
	count := big.NewInt(0)

	sort.Ints(v)

	s0 := int64(v[0]) * int64(t) // max: 1e6 * 1e6 = 1e12
	s0_i := 0
	s0_div_s := s0 / int64(s) // max: 1e12
	s0_mod_s := s0 % int64(s) // max: 1e12

	var sum_div_s, si_count int64

	nodes := make([]Node[int64], len(v))
	new_node := func() *Node[int64] {
		node := &nodes[0]
		nodes = nodes[1:]
		return node
	}

	var tree *Node[int64]
	tree = tree.Insert(s0_mod_s, new_node)

	for i := 1; i < len(v); i++ {
		si := int64(v[i]) * int64(t) // max: 1e6 * 1e6 = 1e12

		if si != s0 {
			sum_div_s += s0_div_s * int64(i-s0_i)        // max: 1e12 * 1e6 = 1e18
			si_div_s := si / int64(s)                    // max: 1e12
			si_count = (si_div_s-1)*int64(i) - sum_div_s // must be > 0, max: (1e12 * 1e6) - 1e18 = 1e18
			s0 = si
			s0_i = i
			s0_div_s = si_div_s
		}

		si_mod_s := si % int64(s)                // max: 1e12
		si_last := int64(tree.FindIdx(si_mod_s)) // max: 1e6
		tree = tree.Insert(si_mod_s, new_node)
		s0_mod_s = si_mod_s

		// count += si_count + si_last // max: (1e18 + 1e6)+... !!! may be int64 overflow !!!
		count.Add(count, big.NewInt(si_count+si_last))
	}

	return count
}

// ----------------------------------------------------------------------------

type Node[K cmp.Ordered] struct {
	key   K
	left  *Node[K]
	right *Node[K]
	size  int
	hght  int
}

func (n *Node[K]) Key() K {
	if n == nil {
		return *new(K)
	}
	return n.key
}

func (n *Node[K]) Size() int {
	if n == nil {
		return 0
	}
	return n.size
}

func (n *Node[K]) height() int {
	if n == nil {
		return 0
	}
	return n.hght
}

// Возвращает индекс (начиная с 0) первого не меньшиго по ключу узла.
// Если такого ключа не существует, возвращает n.Size()
func (n *Node[K]) FindIdx(key K) int {
	if n == nil {
		return 0
	}
	if key <= n.key {
		return n.left.FindIdx(key)
	}
	return n.right.FindIdx(key) + n.left.Size() + 1
}

// (!) Как ни странно, это медленние рекурсии
// // Возвращает индекс (начиная с 0) первого не меньшиго по ключу узла.
// // Если такого ключа не существует, возвращает n.Size()
// func (n *Node[K]) FindIdx(key K) int {
// 	idx := 0
// 	for n != nil {
// 		if key <= n.key {
// 			n = n.left
// 			continue
// 		}
// 		idx += n.left.Size() + 1
// 		n = n.right
// 	}
// 	return idx
// }

func (n *Node[K]) Insert(key K, new_node func() *Node[K]) *Node[K] {
	if n == nil {
		n := new_node()
		*n = Node[K]{
			key:  key,
			size: 1,
			hght: 1,
		}
		return n
	}

	if key <= n.key {
		new_left := n.left.Insert(key, new_node)
		n.left = new_left
		return n.repair()
	}

	new_right := n.right.Insert(key, new_node)
	n.right = new_right
	return n.repair()
}

func (n *Node[K]) update() {
	n.size = n.left.Size() + n.right.Size() + 1
	n.hght = max(n.left.height(), n.right.height()) + 1
}

func (n *Node[K]) repair() *Node[K] {
	d := n.left.height() - n.right.height()
	if d < -1 {
		return n.leftRotate()
	}
	if d > 1 {
		return n.rightRotate()
	}
	n.update()
	return n
}

func (n *Node[K]) leftRotate() *Node[K] {
	al := n
	bt := al.right

	if bt.right.height()-bt.left.height() > 0 {
		al.right = bt.left
		al.update()
		bt.left = al
		bt.update()
		return bt
	}

	ga := bt.left
	al.right = ga.left
	al.update()
	bt.left = ga.right
	bt.update()
	ga.left = al
	ga.right = bt
	ga.update()
	return ga
}

func (n *Node[K]) rightRotate() *Node[K] {
	al := n
	bt := al.left

	if bt.left.height()-bt.right.height() > 0 {
		al.left = bt.right
		al.update()
		bt.right = al
		bt.update()
		return bt
	}

	ga := bt.right
	al.left = ga.right
	al.update()
	bt.right = ga.left
	bt.update()
	ga.right = al
	ga.left = bt
	ga.update()
	return ga
}

// ----------------------------------------------------------------------------

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, t, s, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	v := make([]int, n)
	if err := scanInts(sc, v); err != nil {
		panic(err)
	}

	res := solve(t, s, v)

	fmt.Fprintf(bw, "%d\n", res)
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
