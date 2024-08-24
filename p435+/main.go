package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"unsafe"
)

func solve(n int) (int, int) {
	if n == 1 {
		return 1, 1
	}
	cnt := make([]byte, n+1)
	for i := 2; i <= n/2; i++ {
		for j := i * 2; j <= n; j += i {
			cnt[j]++
		}
	}
	v := 0
	maximum := byte(0)
	for i := range cnt {
		if cnt[i] >= maximum {
			maximum = cnt[i]
			v = i
		}
	}
	return v, int(maximum) + 2
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	v, maximum := solve(n)
	buf := strconv.AppendInt(nil, int64(v), 10)
	buf = append(buf, '\n')
	buf = strconv.AppendInt(buf, int64(maximum), 10)
	bw.Write(buf)
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
