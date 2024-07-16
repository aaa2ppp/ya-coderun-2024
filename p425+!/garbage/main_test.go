package main

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"slices"
	"sort"
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
		// TODO: Add test cases.
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

func Test_intersect(t *testing.T) {
	debugEnable = false
	type args struct {
		ops []Op
	}
	tests := []struct {
		name  string
		args  args
		want  [][]int
		debug bool
	}{
		// 4 3
		// 1 3 1
		// 2 4 2
		// 3 4 4
		{
			"1",
			args{[]Op{{1, 3, 1}, {2, 4, 2}, {3, 4, 4}}},
			[][]int{{1, 2, 4}},
			false,
		},
		{
			"1.2",
			args{[]Op{{1, 3, 1}, {2, 4, 2}, {3, 4, 4}, {4, 5, 5}}},
			[][]int{{1, 2, 4}, {2, 4, 5}},
			false,
		},
		// 7 2
		// 1 5 1
		// 3 7 2
		{
			"2",
			args{[]Op{{1, 5, 1}, {3, 7, 2}}},
			[][]int{{1, 2}},
			false,
		},
		// 10 3
		// 1 1 2
		// 1 1 3
		// 1 1 6
		{
			"3",
			args{[]Op{{1, 1, 2}, {1, 1, 3}, {1, 1, 6}}},
			[][]int{{2, 3, 6}},
			false,
		},
		{
			"4",
			args{[]Op{{6, 9, 1}, {1, 4, 5}, {5, 7, 4}, {2, 4, 10}, {3, 7, 10}, {2, 6, 8}, {1, 4, 3}, {9, 9, 9}, {1, 10, 2}, {1, 7, 4}}},
			[][]int{{1, 2, 4, 4, 8, 10}, {1, 2, 9}, {2, 3, 4, 5, 8, 10, 10}},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			got := intersect(tt.args.ops)
			sortNumSets(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Fuzz_intersect(f *testing.F) {
	debugEnable = false
	type args struct {
		ops []Op
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		// 4 3
		// 1 3 1
		// 2 4 2
		// 3 4 4
		{
			"1",
			args{[]Op{{1, 3, 1}, {2, 4, 2}, {3, 4, 4}}},
			[][]int{{1, 2, 4}},
		},
		{
			"1.2",
			args{[]Op{{1, 3, 1}, {2, 4, 2}, {3, 4, 4}, {4, 5, 5}}},
			[][]int{{1, 2, 4}, {2, 4, 5}},
		},
		// 7 2
		// 1 5 1
		// 3 7 2
		{
			"2",
			args{[]Op{{1, 5, 1}, {3, 7, 2}}},
			[][]int{{1, 2}},
		},
		// 10 3
		// 1 1 2
		// 1 1 3
		// 1 1 6
		{
			"3",
			args{[]Op{{1, 1, 2}, {1, 1, 3}, {1, 1, 6}}},
			[][]int{{2, 3, 6}},
		},
		// TODO: Add test cases.
	}

	const max_n = 10

	for i := 0; i < len(tests); i++ {
		f.Add(int64(i))
	}
	f.Add(int64(-1))

	f.Fuzz(func(t *testing.T, seed int64) {
		var ops []Op
		var want [][]int

		if 0 <= seed && seed < int64(len(tests)) {
			ops = tests[seed].args.ops
			want = tests[seed].want
		} else {
			rd := rand.New(rand.NewSource(seed))
			n := rd.Intn(max_n) + 1
			ops = make([]Op, rd.Intn(n)+1)
			for i := range ops {
				l := rd.Intn(n) + 1
				r := rd.Intn(n) + 1
				x := rd.Intn(n) + 1
				if l > r {
					l, r = r, l
				}
				ops[i] = Op{l, r, x}
			}
			want = bruteforceIntersect(slices.Clone(ops))
		}

		got := intersect(slices.Clone(ops))

		sortNumSets(got)
		sortNumSets(want)
		if !reflect.DeepEqual(got, want) {
			t.Logf("ops: %v", ops)
			t.Errorf("intersect() = \n%v, want \n%v", got, want)
		}
	})
}

func bruteforceIntersect(ops []Op) [][]int {
	segment := make([]uint64, 65)
	for i, op := range ops {
		for j := op.l; j <= op.r; j++ {
			segment[j] |= 1 << i
		}
	}
	sets := map[uint64]struct{}{}
	for _, p := range segment {
		sets[p] = struct{}{}
	}
	num_sets := make([][]int, 0, len(sets))

main_loop:
	for set := range sets {
		for set2 := range sets {
			if set != set2 && set&set2 == set {
				continue main_loop
			}
		}
		var nums []int
		for i := 0; i < len(ops); i++ {
			if set&(1<<i) != 0 {
				nums = append(nums, ops[i].x)
			}
		}
		sort.Ints(nums)
		num_sets = append(num_sets, nums)
	}
	return num_sets
}

func sortNumSets(num_sets [][]int) {
	sort.Slice(num_sets, func(i, j int) bool {
		return compareNumSet(num_sets[i], num_sets[j]) < 0
	})
}

func compareNumSet(set1, set2 []int) int {
	sign := func(a int) int {
		if a < 0 {
			return -1
		} else if a > 0 {
			return 1
		}
		return 0
	}
	for i := 0; i < len(set1) && i < len(set2); i++ {
		v := sign(set1[i] - set2[i])
		if v != 0 {
			return v
		}
	}
	return sign(len(set1) - len(set2))
}

func Test_bruteforceIntersect(t *testing.T) {
	debugEnable = false
	type args struct {
		ops []Op
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		// 4 3
		// 1 3 1
		// 2 4 2
		// 3 4 4
		{
			"1",
			args{[]Op{{1, 3, 1}, {2, 4, 2}, {3, 4, 4}}},
			[][]int{{1, 2, 4}},
		},
		{
			"1.2",
			args{[]Op{{1, 3, 1}, {2, 4, 2}, {3, 4, 4}, {4, 5, 5}}},
			[][]int{{1, 2, 4}, {2, 4, 5}},
		},
		// 7 2
		// 1 5 1
		// 3 7 2
		{
			"2",
			args{[]Op{{1, 5, 1}, {3, 7, 2}}},
			[][]int{{1, 2}},
		},
		// 10 3
		// 1 1 2
		// 1 1 3
		// 1 1 6
		{
			"3",
			args{[]Op{{1, 1, 2}, {1, 1, 3}, {1, 1, 6}}},
			[][]int{{2, 3, 6}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bruteforceIntersect(tt.args.ops); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_intersect(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		intersect(ops)
	}
}

func Test_permut(t *testing.T) {
	debugEnable = false
	type args struct {
		set_n int
		nums  []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantSet []uint16
		debug   bool
	}{
		{
			"1",
			args{5, []int{1, 3, 5}},
			// 1, 3, 4, 5
			4,
			[]uint16{0, 1, 0, 1, 1, 1},
			false,
		},
		{
			"1.2",
			args{6, []int{1, 3, 5}},
			// 1, 3, 4, 5, 6
			5,
			[]uint16{0, 1, 0, 1, 1, 1, 1},
			false,
		},
		{
			"1.3",
			args{6, []int{1, 3, 6}},
			// 1, 3, 4, 6
			4,
			[]uint16{0, 1, 0, 1, 1, 0, 1},
			false,
		},
		{
			"1.4",
			args{9, []int{1, 3, 6}},
			// 1, 3, 4, 6, 7, 9
			6,
			[]uint16{0, 1, 0, 1, 1, 0, 1, 1, 0, 1},
			false,
		},
		{
			"1.5",
			args{11, []int{1, 3, 6}},
			// 1, 3, 4, 6, 7, 9, 10
			7,
			[]uint16{0, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0},
			false,
		},
		{
			"1.5",
			args{11, []int{1, 1, 3, 6}},
			// 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
			11,
			[]uint16{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			set := make([]uint16, tt.args.set_n+1)
			buf := make([]int, tt.args.set_n+1)

			if got := permut(tt.args.nums, buf, set, 1); len(got)-1 != tt.want || !reflect.DeepEqual(set, tt.wantSet) {
				t.Errorf("permut() = %d %v, want %d %v", got, set, tt.want, tt.wantSet)
			}
		})
	}
}
