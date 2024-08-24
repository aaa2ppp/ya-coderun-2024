package main

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	debugEnable = false
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
			args{strings.NewReader(`4 3
1 3 1
2 4 2
3 4 4`)},
			`4
1 2 3 4 `,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`7 2
1 5 1
3 7 2`)},
			`3
1 2 3 `,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`10 3
1 1 2
1 1 3
1 1 6`)},
			`6
2 3 5 6 8 9 `,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`10 3
1 1 1
2 2 2
3 3 3`)},
			`3
1 2 3`,
			false,
			true,
		},
		{
			"4.1",
			args{strings.NewReader(`74 3
65 65 1
66 66 2
67 67 3`)},
			`3
1 2 3`,
			false,
			true,
		},
		{
			"4.2",
			args{strings.NewReader(`140 3
129 129 1
130 130 2
131 131 3`)},
			`3
1 2 3`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`100 3
16 21 3
12 59 4
15 64 5`)},
			`7
3 4 5 7 8 9 12`,
			false,
			true,
		},
		{
			"8",
			args{strings.NewReader(`100 4
3 5 1 21 59 1 55 62 1 29 63 2`)},
			`4
1 2 3 4`,
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
			if gotOut := out.String(); strings.TrimRight(gotOut, " \r\n") != strings.TrimRight(tt.wantOut, " \r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Benchmark_solve_100(t *testing.B) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 100
	q := 100
	ops := make([]Op, q)

	for i := range ops {
		l := rd.Intn(n) + 1
		r := rd.Intn(n) + 1
		x := rd.Intn(n) + 1

		if l > r {
			l, r = r, l
		}
		ops[i] = Op{l, r, x}
	}

	for i := 0; i < t.N; i++ {
		solve(n, ops)
	}
}

func Benchmark_solve_1000(t *testing.B) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1000
	q := 1000
	ops := make([]Op, q)

	for i := range ops {
		l := rd.Intn(n) + 1
		r := rd.Intn(n) + 1
		x := rd.Intn(n) + 1

		if l > r {
			l, r = r, l
		}
		ops[i] = Op{l, r, x}
	}

	for i := 0; i < t.N; i++ {
		solve(n, ops)
	}
}

func Benchmark_solve_10000(t *testing.B) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 10000
	q := 10000
	ops := make([]Op, q)

	for i := range ops {
		l := rd.Intn(n) + 1
		r := rd.Intn(n) + 1
		x := rd.Intn(n) + 1

		if l > r {
			l, r = r, l
		}
		ops[i] = Op{l, r, x}
	}

	for i := 0; i < t.N; i++ {
		solve(n, ops)
	}
}

// func Benchmark_solve_100000(t *testing.B) {
// 	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	n := 100000
// 	q := 100000
// 	ops := make([]Op, q)

// 	for i := range ops {
// 		l := rd.Intn(n) + 1
// 		r := rd.Intn(n) + 1
// 		x := rd.Intn(n) + 1

// 		if l > r {
// 			l, r = r, l
// 		}
// 		ops[i] = Op{l, r, x}
// 	}

// 	for i := 0; i < t.N; i++ {
// 		solve(n, ops)
// 	}
// }

func Fuzz_slove(f *testing.F) {
	debugEnable = false
	type args struct {
		n   int
		ops []Op
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"1",
			args{4, []Op{
				{1, 3, 1},
				{2, 4, 2},
				{3, 4, 4},
			}},
			[]int{1, 2, 3, 4},
		},
		{
			"2",
			args{7, []Op{
				{1, 5, 1},
				{3, 7, 2},
			}},
			[]int{1, 2, 3},
		},
		{
			"3",
			args{10, []Op{
				{1, 1, 2},
				{1, 1, 3},
				{1, 1, 6},
			}},
			[]int{2, 3, 5, 6, 8, 9},
		},
		{
			"4",
			args{10, []Op{
				{1, 1, 1},
				{2, 2, 2},
				{3, 3, 3},
			}},
			[]int{1, 2, 3},
		},
		{
			"4.1",
			args{74, []Op{
				{65, 65, 1},
				{66, 66, 2},
				{67, 67, 3},
			}},
			[]int{1, 2, 3},
		},
		{
			"4.2",
			args{140, []Op{
				{129, 129, 1},
				{130, 130, 2},
				{131, 131, 3},
			}},
			[]int{1, 2, 3},
		},
	}
	for i := range tests {
		f.Add(int64(i))
	}
	const max_n = 100
	f.Fuzz(func(t *testing.T, a int64) {
		var n int
		var ops []Op
		var want []int
		if 0 <= a && a < int64(len(tests)) {
			tt := tests[a]
			n = tt.args.n
			ops = tt.args.ops
			// want = tt.want
		} else {
			rnd := rand.New(rand.NewSource(a))
			n = rnd.Intn(max_n) + 1
			q := rnd.Intn(max_n) + 1
			ops = make([]Op, q)
			for i := range ops {
				l := rnd.Intn(n) + 1
				r := rnd.Intn(n) + 1
				x := rnd.Intn(n) + 1
				if l > r {
					l, r = r, l
				}
				ops[i] = Op{l, r, x}
			}
			// want = slowSolve(n, slices.Clone(ops))
		}
		want = slowSolve(n, slices.Clone(ops))
		got := solve(n, ops)
		if !reflect.DeepEqual(got, want) {
			t.Logf("ops: %v", ops)
			t.Errorf("got = \n%v, \nwant \n%v", got, want)
		}
	})
}
