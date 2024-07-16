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
	val uint64
}

func (b border) String() string {
	if b.is_const() {
		return strconv.Itoa(int(b.val))
	}
	return fmt.Sprintf("%s(%d-%d)", b.ref.name(), b.ref.min(), b.ref.max())
}

func (b border) is_const() bool {
	return b.ref == nil
}

func (b border) min() uint64 {
	if b.ref != nil {
		return b.ref.min()
	}
	return b.val
}

func (b border) max() uint64 {
	if b.ref != nil {
		return b.ref.max()
	}
	return b.val
}

type loop struct {
	id     int
	parent *loop
	left   border
	right  border
	value  uint64
}

func newLoop(id int, parent *loop, left, right border) *loop {
	// if !left.is_const() {
	// 	left.ref.limitMax(right.val)
	// }
	// if !right.is_const() {
	// 	right.ref.limitMin(left.val)
	// }
	return &loop{
		id:     id,
		parent: parent,
		left:   left,
		right:  right,
	}
}

func (lo *loop) limitMin(v uint64) {
	if lo.left.is_const() {
		lo.left.val = max(lo.left.val, v)
	} else {
		lo.left.ref.limitMin(v)
	}
}

func (lo *loop) limitMax(v uint64) {
	if lo.right.is_const() {
		lo.right.val = min(lo.right.val, v)
	} else {
		lo.right.ref.limitMax(v)
	}
}

func (lo *loop) String() string {
	return fmt.Sprintf("{name:%s parent:%s left:%v right:%v value:%v}",
		lo.name(),
		lo.parent.name(),
		lo.left,
		lo.right,
		lo.value,
	)
}

func (lo *loop) name() string {
	if lo == nil {
		return "nil"
	}
	return string(byte(lo.id + 'a'))
}

func (lo *loop) min() uint64 {
	return lo.left.min()
}

func (lo *loop) max() uint64 {
	return lo.right.max()
}

func (lo *loop) exec() (count uint64) {
	if lo == nil {
		return 1
	}

	if debugEnable {
		log.Printf("%v.exec()", lo)
		defer func() {
			log.Printf("%s.exec() -> %d", lo.name(), count)
		}()
	}

	if lo.value != 0 && lo.min() <= lo.value && lo.value <= lo.max() {
		log.Println(1)
		if lo.left.ref != nil {
			log.Println(11)
			old := lo.left.ref.value
			defer func() { lo.left.ref.value = old }()
			lo.left.ref.value = lo.value
		}
		if lo.right.ref != nil {
			log.Println(12)
			old := lo.right.ref.value
			defer func() { lo.right.ref.value = old }()
			lo.right.ref.value = lo.value
		}
		return lo.parent.exec()
	}

	if lo.left.ref != nil {
		r := lo.min()
		l_min := lo.left.min()
		l_max := min(lo.left.max(), r)
		var count uint64
		old := lo.left.ref.value
		for l := l_min; l <= l_max; l++ {
			lo.left.ref.value = l
			count += lo.parent.exec()
			count %= modulo
		}
		lo.left.ref.value = old
		return count
	}

	if lo.right.ref != nil {
		l := lo.min()
		r_min := max(lo.right.min(), l)
		r_max := lo.right.max()
		var count uint64
		old := lo.right.ref.value
		for r := r_min; r <= r_max; r++ {
			lo.right.ref.value = r
			count += lo.parent.exec()
			count %= modulo
		}
		lo.right.ref.value = old
		return count
	}

	return (lo.max() + 1 - lo.min()) * lo.parent.exec()
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
	if c := s[0]; 'a' <= c && c <= 'b' {
		id := c - 'a'
		return border{ref: loops[id]}, nil // may be panic!

	} else {
		val, err := strconv.Atoi(s)
		if err != nil {
			return border{}, err
		}
		return border{val: uint64(val)}, nil
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

	loops := make([]*loop, 0, n)

	var parent *loop
	for id := 0; id < n; id++ {
		left, err := scanBoreder(sc, loops)
		if err != nil {
			return err
		}

		right, err := scanBoreder(sc, loops)
		if err != nil {
			return err
		}

		lo := newLoop(id, parent, left, right)
		loops = append(loops, lo)
		parent = lo
	}

	res := loops[n-1].exec()
	wo := writeOpts{end: '\n'}
	writeInt(bw, res, wo)

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
	~int | ~int64 | ~uint64 | ~int16 | ~int8
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
