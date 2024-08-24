package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unsafe"
)

func main() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	v := make([]int, n)
	for i := range v {
		v[i] = rand.Intn(1000000)
	}

	t := rand.Intn(500000) + 500000
	s := rand.Intn(500000) + 500000

	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	wo := writeOpts{sep: ' ', end: '\n'}
	writeInts(bw, []int{n, t, s}, wo)
	writeInts(bw, v, wo)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
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
