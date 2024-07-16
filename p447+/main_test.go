package main

import (
	"bytes"
	"io"
	"math/rand"
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

func Fuzz_solve(t *testing.F) {
	for i := int64(1); i <= 5; i++ {
		t.Add(i)
	}

	t.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))

		n := r.Intn(9) + 1
		m := r.Intn(9) + 1
		a := genRandInts(r, 100, n)
		b := genRandInts(r, 100, m)

		want := calc(a, b) % 1e9

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
	a := genRandInts(r, 1e9, n)
	b := genRandInts(r, 1e9, m)

	for i := 0; i < t.N; i++ {
		solve(a, b)
	}
}
