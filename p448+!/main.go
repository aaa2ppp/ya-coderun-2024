package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

const modulo = 1e9 + 7

type border struct {
	ref *loop // reference to loop variable
	val int64
}

func (b border) String() string {
	if b.ref == nil {
		return strconv.Itoa(int(b.val))
	}
	return fmt.Sprintf("%d,%s(%d-%d)", b.val, b.ref.name(), b.ref.min(), b.ref.max())
}

type loop struct {
	id      int
	left    border
	right   border
	_counts []int64
}

func newLoop(id int, left, right border) *loop {
	lo := &loop{
		id:    id,
		left:  left,
		right: right,
	}
	if lo.left.ref != nil {
		lo.left.val = lo.left.ref.min()
	}
	if lo.right.ref != nil {
		lo.right.val = lo.right.ref.max()
	}
	return lo
}

func (lo *loop) String() string {
	return fmt.Sprintf("{name:%s left:%v right:%v}",
		lo.name(),
		lo.left,
		lo.right,
	)
}

func (lo *loop) name() string {
	return string(byte(lo.id + 'a'))
}

// (!) RE

// func (lo *loop) limitMin(val int64) {
// 	lo.left.val = max(lo.left.val, val)
// }

// func (lo *loop) limitMax(val int64) {
// 	lo.right.val = min(lo.right.val, val)
// }

// func (lo *loop) updateMinMax() {
// 	if lo.left.ref != nil {
// 		lo.left.val = max(lo.left.val, lo.left.ref.min())
// 	}
// 	if lo.right.ref != nil {
// 		lo.right.val = min(lo.right.val, lo.right.ref.max())
// 	}
// }

func (lo *loop) min() int64 {
	return lo.left.val
}

func (lo *loop) max() int64 {
	return lo.right.val
}

func (lo *loop) counts() []int64 {
	if lo._counts == nil {
		size := lo.max() - lo.min() + 1
		counts := make([]int64, size)
		for i := range counts {
			counts[i] = 1
		}
		lo._counts = counts
	}
	return lo._counts
}

func (lo *loop) calc() (int64, bool) {
	lo_counts := lo.counts()
	lo_min := lo.min()
	lo_max := lo.max()

	if lo.left.ref == nil && lo.right.ref == nil {
		var total int64
		for _, count := range lo_counts {
			total += count
			total %= modulo
		}
		return total, true
	}

	if lo.left.ref != nil {
		ref := lo.left.ref
		ref_counts := ref.counts()
		ref_min := ref.min()
		ref_max := ref.max()

		var count int64
		for i := lo_max; i > ref_max; i-- {
			count += lo_counts[i-lo_min]
			count %= modulo
		}

		for i := ref_max; i > lo_max; i-- {
			if debugEnable {
				log.Printf("%s: %s=%d = %d", lo.name(), ref.name(), i, count)
			}
			ref_counts[i-ref_min] = 0
		}

		for i := min(lo_max, ref_max); i >= max(lo_min, ref_min); i-- {
			count += lo_counts[i-lo_min]
			count %= modulo
			if debugEnable {
				log.Printf("%s: %s=%d = %d", lo.name(), ref.name(), i, count)
			}
			ref_counts[i-ref_min] *= count
			ref_counts[i-ref_min] %= modulo
		}

		for i := lo_min - 1; i >= ref_min; i-- {
			if debugEnable {
				log.Printf("%s: %s=%d = %d", lo.name(), ref.name(), i, count)
			}
			ref_counts[i-ref_min] *= count
			ref_counts[i-ref_min] %= modulo
		}

	} else if lo.right.ref != nil {
		ref := lo.right.ref
		ref_counts := ref.counts()
		ref_min := ref.min()
		ref_max := ref.max()

		var count int64
		for i := lo_min; i < ref_min; i++ {
			count += lo_counts[i-lo_min]
			count %= modulo
		}

		for i := ref_min; i < lo_min; i++ {
			if debugEnable {
				log.Printf("%s: %s=%d = %d", lo.name(), ref.name(), i, count)
			}
			ref_counts[i-ref_min] = 0
		}

		for i := max(lo_min, ref_min); i <= min(lo_max, ref_max); i++ {
			count += lo_counts[i-lo_min]
			count %= modulo
			if debugEnable {
				log.Printf("%s: %s=%d = %d", lo.name(), ref.name(), i, count)
			}
			ref_counts[i-ref_min] *= count
			ref_counts[i-ref_min] %= modulo
		}

		for i := lo_max + 1; i < ref_max; i++ {
			if debugEnable {
				log.Printf("%s: %s=%d = %d", lo.name(), ref.name(), i, count)
			}
			ref_counts[i-ref_min] *= count
			ref_counts[i-ref_min] %= modulo
		}
	}

	return 0, false
}

func scanBoreder(sc *bufio.Scanner, loops []*loop) (border, error) {
	if !sc.Scan() {
		if err := sc.Err(); err == nil {
			return border{}, io.EOF
		} else {
			return border{}, err
		}
	}
	return parseBorder(unsafeString(sc.Bytes()), loops)
}

func parseBorder(s string, loops []*loop) (border, error) {
	if c := s[0]; 'a' <= c && c <= 'z' {
		id := c - 'a'
		return border{ref: loops[id]}, nil // may be panic!

	} else {
		val, err := strconv.Atoi(s)
		if err != nil {
			return border{}, err
		}
		return border{val: int64(val)}, nil
	}
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	loops := make([]*loop, 0, 26)

	for id := 0; id < n; id++ {
		left, err := scanBoreder(sc, loops)
		if err != nil {
			return err
		}

		right, err := scanBoreder(sc, loops)
		if err != nil {
			return err
		}

		loops = append(loops, newLoop(id, left, right))
	}

	// (!) RE
	// // ограничить диапазоны циклов
	// for i := len(loops) - 1; i >= 0; i-- {
	// 	lo := loops[i]
	// 	if lo.left.ref != nil {
	// 		lo.left.ref.limitMax(lo.right.val)
	// 	}
	// 	if lo.right.ref != nil {
	// 		lo.right.ref.limitMin(lo.left.val)
	// 	}
	// }
	// for _, lo := range loops {
	// 	lo.updateMinMax()
	// }

	if debugEnable {
		for _, lo := range loops {
			log.Println(lo)
		}
	}

	for _, lo := range loops {
		if lo.min() > lo.max() {
			wo := writeOpts{end: '\n'}
			writeInt(bw, 0, wo)
			return nil
		}
	}

	// перемножить циклы
	total := int64(1)
	for i := len(loops) - 1; i >= 0; i-- {
		lo := loops[i]
		if count, ok := lo.calc(); ok {
			total *= count
			total %= modulo
		}
	}

	wo := writeOpts{end: '\n'}
	writeInt(bw, total, wo)

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
