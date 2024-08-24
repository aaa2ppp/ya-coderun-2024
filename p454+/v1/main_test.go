package main

import (
	"bytes"
	"io"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
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
		debug   bool
	}{
		// {
		// 	"1",
		// 	args{strings.NewReader(`2 2 2
		// 1 2
		// 2 1`)},
		// 	`2
		// 1 2 `,
		// 	true,
		// },
		{
			"2",
			args{strings.NewReader(`3 3 6
1 -1
1 2
2 -2
2 3
3 -3
3 1`)},
			`0`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`3 3 5
1 1
1 2
1 3
2 3
3 -2`)},
			`2
1 2`,
			true,
		},
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
			run(tt.args.in, out)
			if gotOut := out.String(); strings.TrimRight(gotOut, "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func checkResult(n, m int, ab [][2]int, res []int) bool {
	type topping struct {
		enabled    bool
		interested []int
		// tolging_count int
	}

	type person struct {
		wants           []bool
		dissatisfaction int
	}

	people := make([]person, n+1)
	toppings := make([]topping, m+1)

	for i := range people {
		people[i].wants = make([]bool, m+1)
	}

	for i := range ab {
		a, b := ab[i][0], ab[i][1]

		var want bool
		if b > 0 {
			want = true
		} else {
			b = -b
		}

		people[a].wants[b] = want
		toppings[b].interested = append(toppings[b].interested, a)
	}

	for _, v := range res {
		toppings[v].enabled = true
	}

	for t := 1; t < len(toppings); t++ {
		for _, p := range toppings[t].interested {
			if toppings[t].enabled != people[p].wants[t] {
				if people[p].dissatisfaction != 0 {
					return false
				}
				people[p].dissatisfaction = t
			}
		}
	}

	return true
}

type testArgs struct {
	n, m int
	ab   [][2]int
}

func generateTest(n, m int, seed int64) testArgs {
	type topping struct {
		enabled bool
		// interested    []int
		// tolging_count int
	}

	type person struct {
		wants           []bool
		dissatisfaction int
	}

	rand := rand.New(rand.NewSource(seed))
	people := make([]person, n+1)
	toppings := make([]topping, m+1)

	for t := 1; t < len(toppings); t++ {
		toppings[t].enabled = rand.Intn(100) < 50
	}

	for p := 1; p < len(people); p++ {
		people[p].wants = make([]bool, m+1)
		d := rand.Intn(m + 1)
		people[p].dissatisfaction = d
		for t := 1; t < len(toppings); t++ {
			if t == d {
				people[p].wants[t] = !toppings[t].enabled
			} else {
				people[p].wants[t] = toppings[t].enabled
			}
		}
	}

	var ab [][2]int
	for t := 1; t < len(toppings); t++ {
		for p := 1; p < len(people); p++ {
			if t == people[p].dissatisfaction || rand.Intn(100) < 50 {
				if people[p].wants[t] {
					ab = append(ab, [2]int{p, t})
				} else {
					ab = append(ab, [2]int{p, -t})
				}
			}
		}
	}

	return testArgs{n, m, ab}
}

func Test_checkResult(t *testing.T) {
	debugEnable = false

	type args struct {
		n   int
		m   int
		ab  [][2]int
		res []int
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		debug bool
	}{
		{
			"1",
			args{2, 2, [][2]int{{1, 2}, {2, 1}}, []int{1, 2}},
			true,
			true,
		},
		{
			"1.1",
			args{2, 2, [][2]int{{1, 2}, {2, 1}}, []int{}},
			true,
			true,
		},
		{
			"2",
			args{
				3, 3,
				[][2]int{
					{1, -1},
					{1, 2},
					{2, -2},
					{2, 3},
					{3, -3},
					{3, 1},
				},
				[]int{},
			},
			true,
			true,
		},
		{
			"2.1",
			args{
				3, 3,
				[][2]int{
					{1, -1},
					{1, 2},
					{2, -2},
					{2, 3},
					{3, -3},
					{3, 1},
				},
				[]int{1, 2, 3},
			},
			true,
			true,
		},
		{
			"2.2",
			args{
				3, 3,
				[][2]int{
					{1, -1},
					{1, 2},
					{2, -2},
					{2, 3},
					{3, -3},
					{3, 1},
				},
				[]int{1, 3},
			},
			false,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			if got := checkResult(tt.args.n, tt.args.m, tt.args.ab, tt.args.res); got != tt.want {
				t.Errorf("checkResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_solve(t *testing.T) {
	debugEnable = false

	type args struct {
		n  int
		m  int
		ab [][2]int
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
		debug  bool
	}{
		{
			"1",
			args{2, 2, [][2]int{{1, 2}, {2, 1}}},
			true,
			false,
		},
		{
			"2",
			args{3, 3, [][2]int{{1, -1}, {1, 2}, {2, -2}, {2, 3}, {3, -3}, {3, 1}}},
			true,
			false,
		},
		{
			"3",
			args{3, 3, [][2]int{{1, 1}, {1, 2}, {1, 3}, {2, 3}, {3, -2}}},
			true,
			true,
		},
		{
			"4",
			args{
				3, 3, [][2]int{{1, -1}, {2, -1}, {3, 1}, {1, 2}, {2, 2}, {3, 2}, {1, -3}},
			},
			true,
			false,
		},
		{
			"5",
			args{
				3, 3, [][2]int{{1, -1}, {2, 1}, {3, -1}, {1, 2}, {2, 2}, {1, 3}, {2, -3}},
			},
			true,
			false,
		},
		{
			"6",
			args{
				3, 3, [][2]int{{1, -1}, {2, -1}, {3, 1}, {2, -2}, {3, -2}, {1, 3}, {2, -3}, {3, 3}},
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			got, gotOk := solve(tt.args.n, tt.args.m, tt.args.ab)
			if tt.wantOk && !checkResult(tt.args.n, tt.args.m, tt.args.ab, got) {
				t.Errorf("solve() got = %v: bat result", got)
			}
			if gotOk != tt.wantOk {
				t.Errorf("solve() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Fuzz_solve(f *testing.F) {
	type args = testArgs
	type test struct {
		name   string
		args   args
		wantOk bool
	}
	tests := []test{
		{
			"1",
			args{2, 2, [][2]int{{1, 2}, {2, 1}}},
			true,
		},
		{
			"2",
			args{3, 3, [][2]int{{1, -1}, {1, 2}, {2, -2}, {2, 3}, {3, -3}, {3, 1}}},
			true,
		},
		{
			"3",
			args{3, 3, [][2]int{{1, 1}, {1, 2}, {1, 3}, {2, 3}, {3, -2}}},
			true,
		},
	}

	for i := 0; i < len(tests); i++ {
		f.Add(i)
	}

	debugEnable = false
	f.Fuzz(func(t *testing.T, seed int) {
		var tt test
		if 0 <= seed && seed < len(tests) {
			tt = tests[seed]
		} else {
			args := generateTest(1000, 1000, int64(seed))
			sort.Slice(args.ab, func(i, j int) bool {
				return args.ab[i][1] < args.ab[i][1]
			})
			tt = test{
				strconv.Itoa(seed),
				args,
				true,
			}
		}

		got, gotOk := solve(tt.args.n, tt.args.m, tt.args.ab)
		if gotOk != tt.wantOk {
			t.Logf("name: %v", tt.name)
			t.Logf("args: %v", tt.args)
			t.Errorf("solve() ok = %v, want %v", gotOk, tt.wantOk)
		}
		if gotOk && tt.wantOk && !checkResult(tt.args.n, tt.args.m, tt.args.ab, got) {
			t.Logf("name: %v", tt.name)
			t.Logf("args: %v", tt.args)
			t.Errorf("solve() got = %v: bat result", got)
		}
	})
}
