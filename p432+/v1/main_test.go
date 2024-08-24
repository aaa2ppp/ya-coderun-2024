package main_v1

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
 0  1 -1
 1  0  1
-1  1  0`)},
			`2`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2
0 1
1 0`)},
			`1`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`3
 0 -1 -1
-1  0  1
-1  1  0`)},
			`0`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`5
 0 16 12  1 12
16  0 12 13 -1
12 12  0  5  2
1  13  5  0  2
12 -1  2  2  0`)},
			`4`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`9
0 66 31 87 70 24 55 74 17
66 0 9 3 30 70 71 16 -1
31 9 0 20 58 87 64 92 23
87 3 20 0 61 37 12 47 42
70 30 58 61 0 60 65 7 84
24 70 87 37 60 0 49 34 12
55 71 64 12 65 49 0 98 42
74 16 92 47 7 34 98 0 64
17 -1 23 42 84 12 42 64 0
`)},
			`5`,
			false,
			true,
		},
		// {
		// 	"5",
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
	rnd := rand.New(rand.NewSource(1))
	n := 300
	matrix := makeMatrix(n, n)
	for i := range matrix {
		for j := i + 1; j < n; j++ {
			v := rnd.Intn(100)
			matrix[i][j] = v
			matrix[j][i] = v
		}
	}

	b.Run("xxx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Solve(matrix)
		}
	})
}
