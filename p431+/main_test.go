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
			"13",
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
			"14",
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
			"15",
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
			"16",
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
		{
			"17",
			args{strings.NewReader(`6
2 2
3 3
5 5
3 9
6 3
1 1
`)},
			`No`,
			true,
		},
		{
			"18",
			args{strings.NewReader(`4
-1000000000 0
1000000000 1
-1000000000 1
1000000000 0
`)},
			`Yes`,
			true,
		},

		// -1,0),(0,0),(1,0),(0,-1),(0,-2),(2,1)
		{
			"19",
			args{strings.NewReader(`6
-1 0
0 0
1 0
0 -1
0 -2
2 1
`)},
			`No`,
			true,
		},
		// (-1000000000, -1000000000),
		// (1000000000, 999999999),
		// (999999999, 999999998),
		// (-999999999, -999999999)
		{
			"19",
			args{strings.NewReader(`4
-1000000000 -1000000000
1000000000 999999999
999999999 999999998
-999999999 -999999999
`)},
			`Yes`,
			true,
		},
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
			"1",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 3, y: 3}},
				segment{p0: point{x: 0, y: 3}, p1: point{x: 2, y: 1}},
			},
			true,
		},
		{
			"2",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 1, y: 1}},
				segment{p0: point{x: 0, y: 3}, p1: point{x: 2, y: 1}},
			},
			false,
		},
		{
			"3",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 3, y: 3}},
				segment{p0: point{x: 0, y: 2}, p1: point{x: 3, y: 2}},
			},
			true,
		},
		{
			"4",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 3, y: 3}},
				segment{p0: point{x: 1, y: 0}, p1: point{x: 1, y: 2}},
			},
			true,
		},
		{
			"5",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 2, y: 2}},
				segment{p0: point{x: 1, y: 1}, p1: point{x: 2, y: 0}},
			},
			false,
		},
		{
			"6",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 2, y: 2}},
				segment{p0: point{x: 2, y: 0}, p1: point{x: 1, y: 1}},
			},
			false,
		},
		{
			"7",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 2, y: 2}},
				segment{p0: point{x: 0, y: 2}, p1: point{x: 2, y: 2}},
			},
			false,
		},
		{
			"8",
			args{
				segment{p0: point{x: 0, y: 0}, p1: point{x: 2, y: 2}},
				segment{p0: point{x: 0, y: 2}, p1: point{x: 3, y: 2}},
			},
			false,
		},
		{
			"9",
			args{
				segment{p0: point{x: 1, y: 0}, p1: point{x: 1, y: 2}},
				segment{p0: point{x: 0, y: 1}, p1: point{x: 2, y: 1}},
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

func Fuzz_solve(f *testing.F) {
	tests := [][]point{
		{{0, 0}, {1, 0}, {3, 1}, {0, 1}},
		{{0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}},
		{{-1, 0}, {0, 0}, {1, 0}, {0, -1}, {0, -2}, {2, 1}, {3, 2}},
	}

	for i := range tests {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, a int) {
		var pp []point

		if 0 <= a && a < len(tests) {
			defer func(v bool) {debugEnable=v}(debugEnable)
			debugEnable = true 
			pp = tests[a]
		} else {
			rand := rand.New(rand.NewSource(int64(a)))
			n := rand.Intn(6) + 4
			pp = make([]point, 0, n)
			pm := make(map[point]struct{}, n)
			for len(pp) < n {
				p := point{rand.Intn(9), rand.Intn(9)}
				if _, ok := pm[p]; ok {
					continue
				}
				pp = append(pp, p)
				pm[p] = struct{}{}
			}
		}

		pp2 := make([]point, len(pp))
		copy(pp2, pp)
		want := slowSolve(pp)
		got := solve(pp2)
		if got != want {
			t.Log("pp:", pp)
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
