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

type Op struct {
	l, r, x int
}

func intersect(ops []Op) [][]int {
	var res [][]int

	type event struct {
		pos int
		id  int
	}

	events := make([]event, 0, len(ops)*2)
	for idx := range ops {
		op := &ops[idx]
		id := idx + 1
		events = append(events, event{op.l, id}, event{op.r + 1, -id})
	}

	sort.Slice(events, func(i, j int) bool {
		ev1 := &events[i]
		ev2 := &events[j]
		return ev1.pos < ev2.pos || ev1.pos == ev2.pos && ev1.id < ev2.id
	})
	if debugEnable {
		log.Println("events:", events)
	}

	set := make([]bool, len(ops))
	count := 0
	for i := 0; i < len(events); {
		for ; i < len(events) && events[i].id > 0; i++ {
			idx := events[i].id-1
			set[idx] = true
			count++
		}

		nums := make([]int, 0, count)
		for i := range set {
			if set[i] {
				nums = append(nums, ops[i].x)
			}
		}
		sort.Ints(nums)
		res = append(res, nums)

		for ; i < len(events) && events[i].id < 0; i++ {
			idx := -events[i].id-1
			set[idx] = false
			count--
		}
	}

	return res
}

func permut(nums []int, buf []int, set []uint16, bit uint16) []int {
	buf = buf[:0]
	buf = append(buf, 0)
	for _, v := range nums {
		for i, n := 0, len(buf); i < n; i++ {
			if v2 := buf[i] + v; v2 < len(set) && set[v2]&bit == 0 {
				set[v2] |= bit
				buf = append(buf, v2)
			}
		}
	}
	return buf
}

func solve(n int, ops []Op) []int {
	res_set := make([]uint16, n+1)
	buf := make([]int, n+1)

	nums_set := intersect(ops)
	if debugEnable {
		log.Println("nums_set:", nums_set)
	}

	for i, nums := range nums_set {
		buf = permut(nums, buf, res_set, 1<<i)
	}

	k := 0
	for _, v := range res_set {
		if v != 0 {
			k++
		}
	}

	res := make([]int, 0, k)
	for i := range res_set {
		if res_set[i] != 0 {
			res = append(res, i)
		}
	}

	return res
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, q, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	ops := make([]Op, q)
	for i := range ops {
		l, r, x, err := scanThreeInt(sc)
		if err != nil {
			return err
		}
		ops[i] = Op{l, r, x}
	}

	maximums := solve(n, ops)

	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(bw, len(maximums), wo)
	writeInts(bw, maximums, wo)

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
