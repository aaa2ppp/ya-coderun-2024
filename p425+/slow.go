package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strings"
)

type Row []uint64

func (r Row) String() string {
	var sb strings.Builder
	for _, v := range r {
		sb.WriteString(fmt.Sprintf("%064b", v))
	}
	return sb.String()
}

func (r Row) InitMask(op Op) (Row, int) {
	offset := (op.l - 1) / 64
	size := (op.r+63)/64 - offset
	r = r[:size]
	r.Clear()
	r[0] = math.MaxUint64 << (64 - (op.l-1)%64)

	if rshift := op.r % 64; rshift != 0 {
		r[size-1] |= math.MaxUint64 >> rshift
	}
	return r, offset
}

func (r Row) Clear() {
	for i := range r {
		r[i] = 0
	}
}

func (r Row) Reduce() int {
	for _, v := range r {
		if v != 0 {
			return 1
		}
	}
	return 0
}

// func (r Row) And(a, b Row) {
// 	for i := range r {
// 		r[i] = a[i] & b[i]
// 	}
// }

// func (r Row) AndNot(a, b Row) {
// 	for i := range r {
// 		r[i] = a[i] &^ b[i]
// 	}
// }

// func (r Row) Or(a, b Row) {
// 	for i := range r {
// 		r[i] = a[i] | b[i]
// 	}
// }

func (dst Row) Merge(src, mask Row) {
	// dst[0] |= src[0] &^ mask[0]
	// if n := len(dst) - 1; n > 0 {
	// 	dst[n] |= src[n] &^ mask[n]
	// 	for i := 1; i < n; i++ {
	// 		dst[i] |= src[i]
	// 	}
	// }
	for i := 0; i < len(dst); i++ {
		dst[i] |= src[i] &^ mask[i]
	}
}

type Desk []Row

func NewDesk(n int) Desk {
	row_size := (n + 63) / 64
	desk := make([]Row, n+1)
	for i := range desk {
		desk[i] = make([]uint64, row_size)
	}
	bottom := desk[0]
	for i := range bottom {
		bottom[i] = math.MaxUint64
	}

	if lbits := n % 64; n != 0 {
		bottom[row_size-1] <<= 64 - lbits
	}
	return Desk(desk)
}

func slowSolve(n int, ops []Op) []int {
	panic("slowSolve: BROKEN что-то со сдигами в функции NewDesk и Merge")
	if debugEnable {
		log.Println("slowSolve:", n, ops)
	}

	desk := NewDesk(n)
	if debugEnable {
		log.Printf("%2d    : %v", 0, desk[0])
	}
	buf1 := make(Row, len(desk[0]))

	sort.Slice(ops, func(i, j int) bool {
		return ops[i].x < ops[j].x
	})

	top := 0
	for _, op := range ops {
		mask, j := buf1.InitMask(op)
		j2 := j + len(mask)
		if debugEnable {
			log.Println("op    :", op)
			log.Println("mask  :", mask)
			log.Println("offset:", j)
			log.Println("top   :", top)
		}

		new_top := min(top+op.x, len(desk)-1)
		for i, i2 := new_top-op.x, new_top; i >= 0; i, i2 = i-1, i2-1 {
			src := desk[i][j:j2]
			dst := desk[i2][j:j2]
			Row.Merge(dst, src, mask)
			// if i2 > top && dst.Reduce() != 0 {
			// 	top = i2
			// }
			if debugEnable {
				log.Printf("%2d(%2d): %v", i2, i, dst)
				log.Printf("%2d(%2d): %v", i2, i, desk[i2])
			}
		}
		top = new_top
	}

	var res []int
	for i := 1; i <= top; i++ {
		row := desk[i]
		// log.Printf("%2d: %v", i, row)
		if row.Reduce() != 0 {
			res = append(res, i)
		}
	}

	return res
}

func slowRun(in io.Reader, out io.Writer) error {
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

	maximums := slowSolve(n, ops)

	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(bw, len(maximums), wo)
	writeInts(bw, maximums, wo)

	return nil
}
