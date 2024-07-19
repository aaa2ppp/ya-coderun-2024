package main

import (
	"fmt"
	"testing"
)

func bruteforceSolve(n, a, b int) int {
	bits := make([]byte, n+max(a, b))
	bits[0] = 1
	count := 0
	for i := 0; i < n; i++ {
		if bits[i] == 1 {
			count++
			bits[i+a] = 1
			bits[i+b] = 1
		}
	}
	return count
}

func Test_solve(t *testing.T) {
	debugEnable = true
	type args struct {
		n int
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"1",
			args{6, 2, 3},
			5,
		},
		{
			"2",
			args{6, 58, 3},
			2,
		},
		{
			"3",
			args{6, 2, 18},
			3,
		},
		{
			"4",
			args{45, 40, 41},
			3,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.name = fmt.Sprint(tt.args)
		t.Run(tt.name, func(t *testing.T) {
			if got := solve(tt.args.n, tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			} else {
				t.Logf("solve() = %v", got)
			}
		})
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Fuzz_solve(f *testing.F) {
	type args struct {
		n int
		a int
		b int
	}
	tests := []args{
		{6, 2, 3},
	}
	for _, tt := range tests {
		f.Add(tt.n, tt.a, tt.b)
	}

	fix := func(n, max_val int) int {
		if !(0 < n && n <= max_val) {
			return abs(n)%max_val + 1
		}
		return n
	}

	f.Fuzz(func(t *testing.T, n, a, b int) {
		n = fix(n, 1e6)
		a = fix(a, 1e3)
		b = fix(b, 1e3)
		t.Logf("args = %d, %d, %d", n, a, b)
		want := bruteforceSolve(n, a, b)
		if got := solve(n, a, b); got != want {
			t.Errorf("solve() = %v, want %v", got, want)
		} else {
			t.Logf("solve() = %v", got)
		}
	})
}

func Benchmark_solve(b *testing.B) {
	debugEnable = false
	type args struct {
		n int
		a int
		b int
	}
	tests := []args{
		{1000000, 657, 955},
		{1e16, 657, 955},
		{1e18, 657, 955},
		{1e18, 33469, 65173},
		{1e18, 91493, 91529},
		{1e18, 99709, 99713},
		{1e18, 99708, 99709},
	}

	for _, tt := range tests {
		tt_name := fmt.Sprint(tt)
		b.Run(tt_name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				solve(tt.n, tt.a, tt.b)
			}
		})
	}
}
