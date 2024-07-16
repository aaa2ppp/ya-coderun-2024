package main

import (
	"fmt"
	"testing"
)

func bruteforceSolve(n, a, b, c int) int {
	bits := make([]byte, n+max(a, b, c))
	bits[0] = 1
	count := 0
	for i := 0; i < n; i++ {
		if bits[i] == 1 {
			count++
			bits[i+a] = 1
			bits[i+b] = 1
			bits[i+c] = 1
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
		c int
	}
	tests := []struct {
		args args
		want int
	}{
		{
			args{6, 2, 3, 4},
			5,
		},
		{
			args{10, 3, 4, 5},
			8,
		},
		{
			args{10, 3, 5, 7},
			7,
		},
		{
			args{77, 50, 3, 56},
			35,
		},
		{
			args{98, 21, 4, 22},
			78,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt_name := fmt.Sprint(tt.args)
		t.Run(tt_name, func(t *testing.T) {
			if got := solve(tt.args.n, tt.args.a, tt.args.b, tt.args.c); got != tt.want {
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
	debugEnable = false
	type args struct {
		n int
		a int
		b int
		c int
	}
	tests := []args{
		{6, 2, 3, 4},
		{10, 3, 4, 5},
		{10, 3, 5, 7},
	}
	for _, tt := range tests {
		f.Add(tt.n, tt.a, tt.b, tt.c)
	}

	fix := func(n, max_val int) int {
		if !(0 < n && n <= max_val) {
			return abs(n)%max_val + 1
		}
		return n
	}

	f.Fuzz(func(t *testing.T, n, a, b, c int) {
		n = fix(n, 1e6)
		a = fix(a, 1e3)
		b = fix(b, 1e3)
		c = fix(c, 1e3)
		t.Logf("args = %d, %d, %d, %d", n, a, b, c)
		want := bruteforceSolve(n, a, b, c)
		if got := solve(n, a, b, c); got != want {
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
		c int
	}
	tests := []args{
		{1000000, 123, 657, 955},
		{1e16, 123, 657, 955},
		{1e18, 123, 657, 955},
		{1e18, 99859, 33469, 65173},
		{1e18, 91249, 91493, 91529},
		{1e18, 99707, 99709, 99713},
		{1e18, 99707, 99708, 99709},
	}

	for _, tt := range tests {
		tt_name := fmt.Sprint(tt)
		b.Run(tt_name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				solve(tt.n, tt.a, tt.b, tt.c)
			}
		})
	}
}
