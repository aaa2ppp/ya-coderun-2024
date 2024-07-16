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
1 2 3
3
1 1 3
0 1 2
1 1 2`)},
			`6
6`,
			false,
			true,
		},
		{
			"1",
			args{strings.NewReader(`3
3 1 2
3
1 1 3
0 2 3
1 2 3`)},
			`6
6`,
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

func Benchmark_calcMults(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	n := 100000
	a := make([]int, n)
	for i := range a {
		a[i] = r.Intn(1e7) + 1
	}
	mults := make([][2]int, len(a)+1)
	mults[0] = [2]int{1, 1}
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000000; i++ {
			k := r.Intn(100000)
			if k < 400 {
				op0(a, mults, 0, len(a)-1)
			} else {
				op1(a, mults, 0, len(a)-1)
				_ = r
			}
		}
	}
}
