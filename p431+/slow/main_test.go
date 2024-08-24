package main

import (
	"bytes"
	"io"
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
		debug   bool
	}{
		{
			"0",
			args{strings.NewReader(`4
0 0
1 1
0 1
1 0`)},
			`Yes`,
			true,
		},
		{
			"1",
			args{strings.NewReader(`4
0 0
1 0
3 1
0 1`)},
			`Yes`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5
0 1
0 2
0 3
0 4
0 5`)},
			`No`,
			true,
		},
		{
			"2.2",
			args{strings.NewReader(`5
1 0
2 0
3 0
4 0
5 0`)},
			`No`,
			true,
		},
		{
			"2.3",
			args{strings.NewReader(`5
1 1
2 2
3 3
4 4
5 5`)},
			`No`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`4
0 0
3 1
4 4
6 1`)},
			`No`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`4
0 0
2 2
4 3
6 1`)},
			`Yes`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`6
0 0
0 4
2 3
3 0
4 5
5 3`)},
			`Yes`,
			true,
		},
		{
			"6",
			args{strings.NewReader(`7
0 0
1 1
2 2
3 3
4 2
5 1
6 0`)},
			`Yes`,
			true,
		},
		{
			"7",
			args{strings.NewReader(`5
0 0
4 0
4 2
4 4
8 0
`)},
			`No`,
			true,
		},
		{
			"7",
			args{strings.NewReader(`5
0 0
4 0
6 2
4 4
8 0
`)},
			`Yes`,
			true,
		},
		{
			"8",
			args{strings.NewReader(`5
0 0
1 1
1 2
1 3
2 0
`)},
			`No`,
			true,
		},
		{
			"9",
			args{strings.NewReader(`6
0 0
0 8
8 0
1 1 
2 2
3 3
`)},
			`No`,
			true,
		},
		{
			"10",
			args{strings.NewReader(`7
0 0
0 8
8 0
1 1 
2 2
3 3
4 4
`)},
			`No`,
			true,
		},
		{
			"11",
			args{strings.NewReader(`8
0 0
0 8
8 0
1 1 
2 2
3 3
4 4
4 0
`)},
			`Yes`,
			true,
		},
		{
			"12",
			args{strings.NewReader(`6
0 0
0 8
8 0
4 4 
6 2
2 2
`)},
			`Yes`,
			true,
		},
		{
			"12",
			args{strings.NewReader(`5
0 0
0 8
8 0
4 4 
6 2
`)},
			`No`,
			true,
		},
		{
			"13",
			args{strings.NewReader(`5
6 2
6 4
6 6
5 3 
8 5
`)},
			`Yes`,
			true,
		},
		{
			"14",
			args{strings.NewReader(`5
3 2
5 7
3 3
3 0 
1 1
`)},
			`No`,
			true,
		},
		{
			"14",
			args{strings.NewReader(`6
7 1
5 1
4 1
4 6 
8 1
4 5
`)},
			`Yes`,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

func Test_intersect(t *testing.T) {
	type args struct {
		s1 segment
		s2 segment
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"0",
			args{
				segment{point{0,0},point{1,1}},
				segment{point{0,1},point{1,0}},
			},
			true,
		},
		{
			"0",
			args{
				segment{point{-1000000000,0},point{1000000000,1}},
				segment{point{-1000000000,1},point{1000000000,0}},
			},
			true,
		},
		{
			"0",
			args{
				segment{point{0,0},point{3,0}},
				segment{point{1,0},point{4,0}},
			},
			false,
		},
		{
			"0",
			args{
				segment{point{1,0},point{1,2}},
				segment{point{0,1},point{2,1}},
			},
			true,
		},
		{
			"0",
			args{
				segment{point{1,0},point{1,2}},
				segment{point{0,2},point{2,2}},
			},
			false,
		},
		{
			"0",
			args{
				segment{point{1,0},point{1,2}},
				segment{point{0,0},point{2,0}},
			},
			false,
		},
		{
			"1",
			args{
				segment{point{0, 0}, point{3, 3}},
				segment{point{0, 3}, point{2, 1}},
			},
			true,
		},
		{
			"2",
			args{
				segment{point{0, 0}, point{1, 1}},
				segment{point{0, 3}, point{2, 1}},
			},
			false,
		},
		{
			"3",
			args{
				segment{point{0, 0}, point{3, 3}},
				segment{point{0, 2}, point{3, 2}},
			},
			true,
		},
		{
			"4",
			args{
				segment{point{0, 0}, point{3, 3}},
				segment{point{1, 0}, point{1, 2}},
			},
			true,
		},
		{
			"5",
			args{
				segment{point{0, 0}, point{2, 2}},
				segment{point{1, 1}, point{2, 0}},
			},
			false,
		},
		{
			"6",
			args{
				segment{point{0, 0}, point{2, 2}},
				segment{point{2, 0}, point{1, 1}},
			},
			false,
		},
		{
			"7",
			args{
				segment{point{0, 0}, point{2, 2}},
				segment{point{0, 2}, point{2, 2}},
			},
			false,
		},
		{
			"8",
			args{
				segment{point{0, 0}, point{2, 2}},
				segment{point{0, 2}, point{3, 2}},
			},
			false,
		},
		{
			"9",
			args{
				segment{point{1, 0}, point{1, 2}},
				segment{point{0, 1}, point{2, 1}},
			},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intersect(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
