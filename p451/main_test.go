package main

import (
	"bytes"
	"io"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
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
			"1",
			args{strings.NewReader(`4 1
1 2 3 4`)},
			`4.000000`,
			false,
		},
		{
			"2",
			args{strings.NewReader(`4 2
2 4 3 4`)},
			`3.666667`, // 3.666666
			false,
		},
		{
			"3",
			args{strings.NewReader(`6 3
7 1 2 1 3 6`)},
			`3.333333`,
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
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

var x float64

func Test_solve(t *testing.T) {
	debugEnable = false

	solve := solve
	// solve := slowSolve

	type args struct {
		k  int
		aa []int
	}

	tests := []struct {
		args
		want  float64
		debug bool
	}{
		{
			args{1, []int{1, 2, 3, 4}},
			4.000000,
			true,
		},
		{
			args{2, []int{2, 4, 3, 4}},
			3.666666,
			true,
		},
		{
			args{3, []int{7, 1, 2, 1, 3, 6}},
			3.333333,
			true,
		},
		{
			args{36, []int{67, 86, 68, 80, 70, 61, 71, 52, 21, 68, 74, 44, 45, 71, 76, 1, 13, 31, 17, 69, 37, 15, 42, 46, 75, 5, 20, 56, 88, 44, 81, 94, 68, 90, 13, 90, 57, 24, 69, 50, 89, 9, 2, 50, 4, 19, 68, 12, 39}},
			54.585366,
			true,
		},
		{
			args{61, []int{25, 25, 75, 5, 25, 93, 42, 22, 13, 91, 11, 10, 36, 99, 24, 7, 65, 84, 86, 56, 96, 56, 52, 86, 62, 10, 40, 28, 15, 11, 91, 55, 13, 3, 93, 90, 11, 90, 13, 22, 68, 80, 84, 81, 3, 71, 48, 19, 17, 40, 25, 38, 72, 76, 1, 100, 8, 9, 72, 3, 15, 96, 100, 52, 79, 80, 12, 53, 50, 12, 82, 13, 88, 43, 25, 39, 91, 57, 5, 84, 35, 63, 69, 18, 97, 24, 19, 15, 22, 11, 84, 35}},
			50.739130,
			true,
		},
		{
			args{17, []int{40, 46, 20, 87, 77, 46, 22, 12, 100, 41, 22, 68, 42, 99, 79, 74, 24, 40, 78, 80, 99, 66, 9, 27, 64, 11, 19, 5, 10, 54, 21, 7, 8, 23, 56, 4, 6, 72, 46, 38, 17, 68, 67, 52, 100, 46, 92, 70, 44, 55, 3, 11, 75}},
			60.842105,
			true,
		},
		{
			args{9, []int{35,66,100,29,100,74,98,20,58,41,87,61,16,57,13,100,39,66,96,51,41,63,18,91,51,50,21,26,78,93,13,35,36,36,39,64,11,72}},
			67.444444,
			true,
		},
		{
			args{6, []int{13,12,19,15,62,29,89,93,85,7,43,10,3,94,97,63,1,2,67,40,57,29,23}},
			62.166667,
			true,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			if got := solve(tt.args.k, tt.args.aa); math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("got %.6f, want %.6f", got, tt.want)
			}
		})
	}
}

func Test_1(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixMilli()))
	n := 10
	aa := make([]int, n)
	debugEnable = true
	for i := 0; i < 10; i++ {
		log.Println("-------------------")
		for i := range aa {
			aa[i] = rand.Intn(9) + 1
		}
		x = slowSolve(1, aa)
	}
}

func Fuzz_solve(f *testing.F) {
	type args struct {
		k  int
		aa []int
	}
	tests := []args{
		{
			1,
			[]int{1, 2, 3, 4},
		},
		{
			2,
			[]int{2, 4, 3, 4},
		},
		{
			3,
			[]int{7, 1, 2, 1, 3, 6},
		},
	}
	for i := range tests {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		var tt args
		if 0 <= i && i < len(tests) {
			tt = tests[i]
		} else {
			n := rand.Intn(100) + 1
			k := rand.Intn(n) + 1
			aa := make([]int, n)
			for i := range aa {
				aa[i] = rand.Intn(100) + 1
			}
			tt = args{k, aa}
		}
		want := slowSolve(tt.k, tt.aa)
		got := solve(tt.k, tt.aa)
		if math.Abs(got-want) > 1e-6 {
			t.Logf("%d %d %v", len(tt.aa), tt.k, tt.aa)
			t.Errorf("got %.6f, want %.6f", got, want)
		}
	})
}
