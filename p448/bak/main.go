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

const modulo = 1e9 + 7

type border struct {
	ref      *loop // reference to loop variable
	val      uint32
	is_const bool
}

func (b border) min() uint32 {
	if b.is_const {
		return b.val
	}
	return b.ref.min()
}

func (b border) max() uint32 {
	if b.is_const {
		return b.val
	}
	return b.ref.max()
}

func (b border) valueCount(v uint32) uint32 {
	if b.is_const {
		return 1
	}
	return b.ref.valueCount(v)
}

type loop struct {
	id          byte
	parent      *loop
	left        border
	right       border
	value_count []uint32 // count of uses of each value of loop variable
	min_value   uint32   // min variable value
	max_value   uint32   // max variable value
	total_count uint32
	done        bool
}

func (lo *loop) name() string {
	return string(lo.id + 'a')
}

func (lo *loop) do() {
	lo.min_value = lo.left.min()
	lo.max_value = lo.right.max()

	if lo.min_value > lo.max_value {
		lo.min_value = math.MaxUint32
		lo.max_value = 0
		lo.done = true
		return
	}

	lo.value_count = make([]uint32, lo.max_value+1-lo.min_value+1)

	if lo.left.is_const && lo.right.is_const {
		count := uint32(1)
		if lo.parent != nil {
			count = lo.parent.totalCount()
		}
		for i := range lo.value_count {
			lo.value_count[i] = count
		}

	} else if lo.left.is_const {
		l := lo.min_value
		r, r_min := lo.max_value, max(l, lo.right.min())
		i := len(lo.value_count) - 2 
		for ; r >= r_min ; r, i = r-1, i-1 {
			count := lo.parent.valueCount(r)
			lo.value_count[i] = (lo.value_count[i+1] + count) % modulo
		}
		for ; r >= l; r, i = r-1, i-1 {
			lo.value_count[i] = lo.value_count[i+1]
		}
		lo.value_count = lo.value_count[:len(lo.value_count)-1]

	} else { //if lo.right.is_const {
		r := lo.max_value
		l, l_max := lo.min_value, min(r, lo.left.max())
		i := 1
		for ; l <= l_max; l, i = l+1, i+1 {
			count := lo.parent.valueCount(l)
			lo.value_count[i] = (lo.value_count[i-1] + count)%modulo
		}
		for ; l <= r; l, i = l+1, i+1 {
			lo.value_count[i] = lo.value_count[i-1]
		}
		lo.value_count = lo.value_count[1:]
	}

	for _, count := range lo.value_count {
		lo.total_count += count
		lo.total_count %= modulo
	}

	log.Println(lo.name(), lo.min_value, lo.max_value, lo.value_count)

	lo.done = true
}

func (lo *loop) min() uint32 {
	if debugEnable {
		log.Printf("%s.min()", lo.name())
	}
	if !lo.done {
		lo.do()
	}
	return lo.min_value
}

func (lo *loop) max() uint32 {
	if debugEnable {
		log.Printf("%s.max()", lo.name())
	}
	if !lo.done {
		lo.do()
	}
	return lo.max_value
}

func (lo *loop) valueCount(val uint32) uint32 {
	if debugEnable {
		log.Printf("%s.valueCount(%d)", lo.name(), val)
	}
	if !lo.done {
		lo.do()
	}
	i := val - lo.min_value
	if i < uint32(len(lo.value_count)) {
		return lo.value_count[i]
	}
	return 0
}

func (lo *loop) totalCount() uint32 {
	if debugEnable {
		log.Printf("%s.totalCount()", lo.name())
	}
	if !lo.done {
		lo.do()
	}
	return lo.total_count
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
		return border{ref: loops[c-'a']}, nil // may be panic!

	} else {
		val, err := strconv.Atoi(s)
		if err != nil {
			return border{}, err
		}
		return border{val: uint32(val), is_const: true}, nil
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
	for i := 0; i < n; i++ {
		left, err := scanBoreder(sc, loops)
		if err != nil {
			return err
		}

		right, err := scanBoreder(sc, loops)
		if err != nil {
			return err
		}

		lo := &loop{
			id:     byte(i),
			parent: parent,
			left:   left,
			right:  right,
		}

		loops = append(loops, lo)
		parent = lo
	}

	res := loops[n-1].totalCount()
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
	~int | ~int64 | ~uint32 | ~int16 | ~int8
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
