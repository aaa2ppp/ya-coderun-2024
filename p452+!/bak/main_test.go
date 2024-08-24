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
		wantErr bool
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`2 4
1 0 2 0
3 5 4 0
2 1 0 0
3 0 4 5`)},
			`2`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`3 3
1 2 3
4 5 6
7 8 0
4 2 3
6 5 1
0 7 8`)},
			`4`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`2 2
1 2
3 4
2 3
4 1`)},
			`-1`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`1 8
7 6 5 4 3 2 1 0
1 2 3 4 5 6 7 0`)},
			`6`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`1 8
0 6 1 7 2 4 5 3
0 7 6 5 4 3 2 1`)},
			`4`,
			false,
			true,
		},
		{
			"6",
			args{strings.NewReader(`4 4
1 2 3 4
5 6 7 8
9 10 11 12
13 14 15 0
5 2 3 4
1 6 7 8
13 10 11 12
9 14 15 0`)},
			`5`,
			false,
			true,
		},
		{
			"7",
			args{strings.NewReader(`5 4
1 2 3 4
5 6 7 8
9 10 11 12
13 14 15 0
0 0 0 0
5 2 3 10
1 6 7 8
13 0 11 12
9 14 15 0
4 0 0 0`)},
			`6`,
			false,
			true,
		},
		{
			"8",
			args{strings.NewReader(`5 4
1 2 3 4
5 6 7 8
9 10 11 12
13 14 15 0
0 0 0 0
5 2 3 10
1 6 7 8
13 4 11 12
9 14 15 0
0 0 0 0`)},
			`6`,
			false,
			true,
		},
		{
			"9",
			args{strings.NewReader(`8 3
1 2 3
4 5 6
7 8 9
10 11 12
13 14 15
16 17 18
19 20 21
22 23 0
4 2 3
7 5 11
1 22 9
13 6 12
16 14 15
10 17 18
19 20 23
8 21 0`)},
			`12`,
			false,
			true,
		},
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

func Test_findLISLength(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"1234",
			args{[]int{1, 2, 3, 4}},
			4,
		},
		{
			"4321",
			args{[]int{4, 3, 2, 1}},
			1,
		},
		{
			"4132",
			args{[]int{4, 1, 3, 2}},
			2,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findLISLength(tt.args.a); got != tt.want {
				t.Errorf("findLISLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
