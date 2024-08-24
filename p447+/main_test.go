package main

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`3
2 3 5
2
4 5`)},
			`10`,
			false,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			if err := run(tt.args.in, out); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); strings.TrimRight(gotOut, "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func saveTest(fn string, a, b []int, res int, truncated bool) {
	f, _ := os.Create(fn)
	o := bufio.NewWriter(f)
	wo := writeOpts{sep: ' ', end: '\n'}
	writeInt(o, len(a), wo)
	writeInts(o, a, wo)
	writeInt(o, len(b), wo)
	writeInts(o, b, wo)
	o.Flush()
	f.Close()
	f, _ = os.Create(fn + ".a")
	o = bufio.NewWriter(f)
	// if res >= 1e9 {
	// 	truncated = true
	// 	res %= 1e9
	// }
	// if truncated {
	// fmt.Fprintf(o, "%09d\n", res)
	// } else {
	writeInt(o, res, wo)
	// }
	o.Flush()
	f.Close()
}

func Fuzz_solve(t *testing.F) {
	for i := int64(1); i <= 5; i++ {
		t.Add(i)
	}

	count := 0

	t.Fuzz(func(t *testing.T, seed int64) {
		count++

		r := rand.New(rand.NewSource(seed))

		n := r.Intn(9) + 1
		m := r.Intn(9) + 1
		a := genRandInts(r, 100, n)
		b := genRandInts(r, 100, m)

		want := calc(a, b) % 1e9

		// if count < 100 {
		// 	fn := fmt.Sprintf("./test_data/%d", count+1)
		// 	saveTest(fn, a, b, want, false)
		// }

		t.Log(a)
		t.Log(b)
		if got, _ := solve(a, b); got != want {
			t.Errorf("solve() = %v, want %v", got, want)
		} else {
			t.Logf("solve() = %v", got)
		}
	})
}

func calc(a, b []int) int {
	aa := 1
	for _, v := range a {
		aa *= v
	}
	bb := 1
	for _, v := range b {
		bb *= v
	}
	return gcd(aa, bb)
}

func genRandInts(r *rand.Rand, maxV int, n int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = r.Intn(maxV) + 1
	}
	return a
}

func Benchmark_solve(t *testing.B) {
	r := rand.New(rand.NewSource(1))

	n := 977
	m := 997
	a := genRandInts(r, 1e9-1, n)
	for i := range a {
		if a[i]%10 == 0 || a[i]%5 == 0 {
			a[i] = a[i] + 1
		}
	}
	b := genRandInts(r, 1e9, m)
	for i := range b {
		if b[i]%10 == 0 || b[i]%5 == 0 {
			b[i] = b[i] + 1
		}
	}
	aa := make([]int, len(a))
	bb := make([]int, len(b))

	// copy(aa, a)
	// copy(bb, b)
	// res, truncated := solve(aa, bb)
	// saveTest("./test_data/977", a, b, res, truncated)

	for i := 0; i < t.N; i++ {
		copy(aa, a)
		copy(bb, b)
		_, _ = solve(aa, bb)
	}
}
