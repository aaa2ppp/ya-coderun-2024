package main

import (
	"bytes"
	"io"
	"math"
	"math/rand"
	"slices"
	"sort"
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
			args{strings.NewReader(`1 1
0 5
0 4
1 3`)},
			`1`,
			false,
			false,
		},
		{
			"2",
			args{strings.NewReader(`1 2
0 5
0 4
0 1
2 5`)},
			`3`,
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`3 2
0 16
1 4
8 11
13 15
6 9
10 12`)},
			`2`,
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`4 2
0 16
1 4
8 11
7 10
13 15
6 9
10 12`)},
			``,
			true,
			false,
		},
		{
			"3",
			args{strings.NewReader(`4 2
0 16
1 4
8 11
9 10
13 15
6 9
10 12`)},
			``,
			true,
			false,
		},
		{
			"3",
			args{strings.NewReader(`4 2
0 16
1 4
8 11
9 12
13 15
6 9
10 12`)},
			``,
			true,
			false,
		},
		{
			"4",
			args{strings.NewReader(`6 7 0 9
0 1 1 2 3 3 4 5 5 5 5 6
0 0 2 2 2 3 4 6 6 7 7 7 8 9 9 9`)},
			`2`,
			false,
			false,
		},
		{
			"5",
			args{strings.NewReader(`3 3 0 9
1 3 3 3 5 6
2 2 6 6 6 8`)},
			`0`,
			false,
			false,
		},
		{
			"6",
			args{strings.NewReader(`5 5 0 10
8 9 4 6 7 8 1 3 6 6
8 9 1 1 9 10 3 3 5 6`)},
			`2`,
			false,
			false,
		},
		{
			"7",
			args{strings.NewReader(`3 4 0 10
0 1 1 4 5 5
2 4 4 5 5 5 10 10`)},
			`2`,
			false,
			false,
		},
		{
			"8",
			args{strings.NewReader(`2 5 0 10
4 6 7 9
0 1 1 1 1 2 5 6 7 9`)},
			`1`,
			false,
			false,
		},
		{
			"9",
			args{strings.NewReader(`5 3 0 10
10 10 6 7 8 10 0 6 6 6
7 8 9 10 4 7`)},
			`4`,
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

func Fuzz_solve_10(f *testing.F) {
	debugEnable = false
	for i := 0; i < 5; i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, a int) {
		t.Log(a)
		rnd := rand.New(rand.NewSource(int64(a)))

		n := (rnd.Intn(5) + 1)
		m := (rnd.Intn(5) + 1)

		l, r := 0, 10

		blue := make([]int, n*2)
		for i := range blue {
			blue[i] = l + rnd.Intn(r-l+1)
		}

		red := make([]int, m*2)
		for i := range red {
			red[i] = l + rnd.Intn(r-l+1)
		}

		sort.Ints(blue)
		t.Log(blue)
		sort.Ints(red)
		t.Log(red)

		blue_segs := make([][2]int, n)
		for i := range blue_segs {
			blue_segs[i] = [2]int{blue[i*2], blue[i*2+1]}
		}
		rnd.Shuffle(len(blue_segs), func(i, j int) {
			blue_segs[i], blue_segs[j] = blue_segs[j], blue_segs[i]
		})

		red_segs := make([][2]int, m)
		for i := range red_segs {
			red_segs[i] = [2]int{red[i*2], red[i*2+1]}
		}
		rnd.Shuffle(len(red_segs), func(i, j int) {
			red_segs[i], red_segs[j] = red_segs[j], red_segs[i]
		})

		t.Log(l, r)
		t.Log(blue_segs)
		t.Log(red_segs)
		want := bruteforceSolve(l, r, slices.Clone(blue), slices.Clone(red))
		if got := solve(l, r, blue_segs, red_segs); got != want {
			t.Errorf("solve() = %v, want %v", got, want)
		}
	})
}

func Fuzz_solve_1000(f *testing.F) {
	debugEnable = false
	for i := 0; i < 5; i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, a int) {
		t.Log(a)
		rnd := rand.New(rand.NewSource(int64(a)))

		n := (rnd.Intn(100) + 1)
		m := (rnd.Intn(100) + 1)

		l, r := 0, 1000

		blue := make([]int, n*2)
		for i := range blue {
			blue[i] = l + rnd.Intn(r-l+1)
		}

		red := make([]int, m*2)
		for i := range red {
			red[i] = l + rnd.Intn(r-l+1)
		}

		sort.Ints(blue)
		sort.Ints(red)

		t.Log(blue)
		t.Log(red)

		blue_segs := make([][2]int, n)
		for i := range blue_segs {
			blue_segs[i] = [2]int{blue[i*2], blue[i*2+1]}
		}
		rnd.Shuffle(len(blue_segs), func(i, j int) {
			blue_segs[i], blue_segs[j] = blue_segs[j], blue_segs[i]
		})

		red_segs := make([][2]int, m)
		for i := range red_segs {
			red_segs[i] = [2]int{red[i*2], red[i*2+1]}
		}
		rnd.Shuffle(len(red_segs), func(i, j int) {
			red_segs[i], red_segs[j] = red_segs[j], red_segs[i]
		})

		t.Log(l, r)
		t.Log(blue_segs)
		t.Log(red_segs)
		want := bruteforceSolve(l, r, slices.Clone(blue), slices.Clone(red))
		if got := solve(l, r, blue_segs, red_segs); got != want {
			t.Errorf("solve() = %v, want %v", got, want)
		}
	})
}

func bruteforceSolve(l, r int, blue, red []int) int {
	n, m := len(blue), len(red)
	_ = n

	min_intersection := math.MaxInt
	for x, end := l, r-(red[m-1]-red[0]); x <= end; x++ {
		moveTo(red, x)
		inter := calcIntesection(blue, red)
		// log.Println(x, red, inter)
		min_intersection = min(min_intersection, inter)
	}
	return min_intersection
}
