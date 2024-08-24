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

var weekDays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// возвращает искомое время или -1, если решения нет
func solve(a_start, b_start, a_round, b_round int) int {
	d := gcd(a_round, b_round)

	if (a_start-b_start)%d != 0 {
		return -1
	}

	if a_round > b_round {
		a_start, b_start = b_start, a_start
		a_round, b_round = b_round, a_round
	}

	t := b_start - a_start
	if t < 0 {
		t = b_round - (a_start-b_start)%b_round
	}

	for t%a_round != 0 {
		t += b_round
	}

	return a_start + t
}

func gcd(a, b int) int {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func formatDayTime(t int) string {
	h := (t / 60) % 24
	m := t % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func parseDayTime(s string) (int, error) {
	h, err := strconv.Atoi(s[:2])
	if err != nil {
		return 0, err
	}
	m, err := strconv.Atoi(s[3:])
	if err != nil {
		return 0, err
	}
	return h*60 + m, nil
}

func scanDayTime(sc *bufio.Scanner) (int, error) {
	s, err := scanWord(sc)
	if err != nil {
		return 0, err
	}
	return parseDayTime(s)
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	a_start, err := scanDayTime(sc)
	if err != nil {
		return err
	}
	b_start, err := scanDayTime(sc)
	if err != nil {
		return err
	}
	a_round, err := scanDayTime(sc)
	if err != nil {
		return err
	}
	b_round, err := scanDayTime(sc)
	if err != nil {
		return err
	}

	t := solve(a_start, b_start, a_round, b_round)
	if t == -1 {
		bw.WriteString("Never")
		return nil
	}

	bw.WriteString(weekDays[(6+t/(24*60))%7])
	bw.WriteByte('\n')
	bw.WriteString(formatDayTime(t))
	bw.WriteByte('\n')

	return nil
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanWord(sc *bufio.Scanner) (string, error) {
	if !sc.Scan() {
		return "", io.EOF
	}
	return sc.Text(), nil
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
