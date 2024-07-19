// test
package main

import (
	"bytes"
	"io"
	"strconv"
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
			args{strings.NewReader(`6
2 3 4`)},
			`5`,
			false,
			true,
		},
		{
			"2.1",
			args{strings.NewReader(`1000000
		123 657 955`)},
			strconv.Itoa(calc(1000000, 123, 657, 955)),
			false,
			false,
		},
		{
			"2.2",
			args{strings.NewReader(`1000000000000000000
		123 657 955`)},
			"999999999999991606", // ???
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`11
2 4 6`)},
			`6`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`10 5 3 5`)},
			`6`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`16 33 93 6`)},
			`3`,
			false,
			true,
		},
		{
			"6",
			args{strings.NewReader(`18 16 49 6`)},
			`4`,
			false,
			true,
		},
		// {
		// 	"6.2",
		// 	args{strings.NewReader(`1000000000000000000
		// 91249 91493 91529`)},
		// 	`xxx`,
		// 	false,
		// 	false,
		// },
		// {
		// 	"6.3",
		// 	args{strings.NewReader(`1000000000000000000
		// 99707 99709	99713`)},
		// 	`xxx`,
		// 	false,
		// 	false,
		// },
		// {
		// 	"6.4",
		// 	args{strings.NewReader(`1000000000000000000
		// 99707 99708	99709`)},
		// 	`xxx`,
		// 	false,
		// 	false,
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
			if gotOut := out.String(); strings.TrimRight(gotOut, "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			} else {
				t.Logf("run() = %v", gotOut)
			}
		})
	}
}

func Benchmark_run(b *testing.B) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"1",
			args{strings.NewReader(`1000000000000000000
		99859 33469 65173`)},
		},
		{
			"2",
			args{strings.NewReader(`1000000000000000000
		91249 91493 91529`)},
		},
		{
			"3",
			args{strings.NewReader(`1000000000000000000
		99707 99709	99713`)},
		},
		{
			"4",
			args{strings.NewReader(`1000000000000000000
		99707 99708	99709`)},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = run(tt.args.in, io.Discard)
			}
		})
	}
}

func Test_gcd(t *testing.T) {
	oldDebugEnable := debugEnable
	defer func() {
		debugEnable = oldDebugEnable
	}()

	type args struct {
		a int
		b int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		debug bool
	}{
		{
			"4, 6",
			args{4, 6},
			2,
			true,
		},
		{
			"6, 4",
			args{6, 4},
			2,
			true,
		},
		{
			"96, 33",
			args{96, 33},
			3,
			false,
		},
		{
			"33, 96",
			args{33, 96},
			3,
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			if got := gcd(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("gcd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Fuzz_solve(t *testing.F) {
	type args struct {
		n int
		a int
		b int
		c int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"6 2 3 4",
			args{6, 2, 3, 4},
		},
		{
			"11 2 4 6",
			args{11, 2, 4, 6},
		},
		{
			"10 2 4 6",
			args{10, 2, 4, 6},
		},
		{
			"9 2 4 6",
			args{9, 2, 4, 6},
		},
		{
			"1000 33 93 6",
			args{1000, 33, 93, 6},
		},
		{
			"100000 75 100 125",
			args{100000, 75, 100, 125},
		},
		{
			"1000000 75 100 125",
			args{1000000, 75, 100, 125},
		},
		{
			"10000000 75 100 125",
			args{10000000, 75, 100, 125},
		},
		{
			"100000000 75 100 125",
			args{100000000, 75, 100, 125},
		},
		{
			"1000000000 75 100 125",
			args{1000000000, 75, 100, 125},
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Add(uint(tt.args.n), uint(tt.args.a), uint(tt.args.b), uint(tt.args.c))
	}
	t.Fuzz(func(t *testing.T, nn, aa, bb, cc uint) {
		nn %= 10 ^ 18
		if nn == 0 {
			nn++
		}

		aa %= 100001
		bb %= 100001
		cc %= 100001

		if aa == 0 {
			aa++
		}
		if bb == 0 {
			bb++
		}
		if cc == 0 {
			cc++
		}

		n, a, b, c := int(nn), int(aa), int(bb), int(cc)
		t.Logf("n, a, b, c = %d,%d,%d,%d", n, a, b, c)

		want := calc(n, a, b, c)
		if got := solve(n, a, b, c); got != want {
			t.Errorf("solve() = %v, want %v", got, want)
		} else {
			t.Logf("solve() = %v", got)
		}
	})
}

func calc(n int, a, b, c int) int {
	p := make([]bool, n+max(a, b, c))
	p[0] = true
	count := 0
	for i := 0; i < n; i++ {
		if !p[i] {
			continue
		}
		count++
		p[i+a] = true
		p[i+b] = true
		p[i+c] = true
	}
	return count
}
