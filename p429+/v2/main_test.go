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
			args{strings.NewReader(`3 4
6 6
3 1 1
3 3 1
3 5 1
1 1
1 2
1 3
1 4`)},
			`14
23
23
14`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1 2
22 22
11 11 5
3 1
4 1`)},
			`1234
1`,
			false,
			true,
		},
		{
			"11",
			args{strings.NewReader(`1 2
4 4
3 2 1
1 1
1 4`)},
			`14
14`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`1 2
4 4
1 2 1
1 2
1 3`)},
			`23
23`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`1 2
4 4
2 3 1
1 1
1 2`)},
			`12
12`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`1 2
4 4
2 1 1
1 4
1 3`)},
			`34
34`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`1 4
4 6
3 3 1
1 1
1 2
1 3
1 4`)},
			`1234
1234
1234
1234`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`2 4
6 6
1 3 1
5 3 1
1 1
1 2
1 3
1 4`)},
			`1234
1234
1234
1234`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`2 3
6 6
1 3 1
5 4 1
1 1
1 2
1 4`)},
			`124
124
124`,
			false,
			true,
		},
		{
			"12",
			args{strings.NewReader(`2 3
6 6
1 3 1
5 2 1
1 1
1 3
1 4`)},
			`134
134
134`,
			false,
			true,
		},
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

func Benchmark_solve(b *testing.B) {
	//rand := rand.New(rand.NewSource(1))
	n := 2000
	m := 100000

	h := int(1e9)
	w := int(1e9)
	d := min(w, h)

	bb := make([][3]int, n)
	for i := range bb {
		bb[i] = [3]int{rand.Intn(w) + 1, rand.Intn(h) + 1, rand.Intn((d+1)/2) + 1}
	}
	qq := make([][2]int, m)
	for i := range qq {
		qq[i] = [2]int{rand.Intn(d) + 1, rand.Intn(4)}
	}

	b.Run("xxx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			solve(w, h, bb, qq)
		}
	})
}
