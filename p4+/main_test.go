package main

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
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
		wantErr bool
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`5 10 10
10 2 3 4 5`)},
			`26`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`10 10 10
1 2 3 4 5 6 7 8 9 10`)},
			`120`,
			false,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
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

func Fuzz_solve(f *testing.F) {
	const max_n = 1000

	type test struct {
		t int
		s int
		v []int
	}
	tests := []test{
		{
			10, 10,
			[]int{10, 2, 3, 4, 5},
		},
		{
			10, 10,
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			1, 6,
			[]int{1, 2, 3, 4, 5},
		},
		{
			1, 8,
			[]int{1, 1, 4, 5},
		},
		{
			7, 2,
			[]int{1, 2, 8, 8, 8, 9, 9, 10, 10},
		},
	}

	for i := range tests {
		f.Add(int64(i))
	}

	f.Fuzz(func(t *testing.T, seed int64) {
		var tt test
		if 0 <= seed && seed < int64(len(tests)) {
			tt = tests[seed]
		} else {
			rand := rand.New(rand.NewSource(seed))
			n := rand.Intn(max_n) + 1
			tt.t = rand.Intn(max_n) + 1
			tt.s = rand.Intn(max_n) + 1
			tt.v = make([]int, n)
			for i := range tt.v {
				tt.v[i] = rand.Intn(max_n) + 1
			}
		}

		if len(tt.v) <= 10 {
			t.Log(tt.t, tt.s, tt.v)
		}

		want := int64(slowSolve(tt.t, tt.s, tt.v))
		gotBig := solve(tt.t, tt.s, tt.v)

		if !gotBig.IsInt64() {
			panic(fmt.Sprintf("can't represent %d as int64", gotBig))
		}
		got := gotBig.Int64()

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

var benchRes *big.Int
var benchData []int

func init() {
	n := 1000000
	v := make([]int, n)
	for i := range v {
		v[i] = rand.Intn(1000000)
	}
	benchData = v
}

func Benchmark_solve(b *testing.B) {
	rand := rand.New(rand.NewSource(1))
	t := rand.Intn(500000) + 500000
	s := rand.Intn(500000) + 500000
	v := benchData

	b.Run("100K", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchRes = solve(t, s, v[:100000])
		}
	})

	b.Run("250K", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchRes = solve(t, s, v[:250000])
		}
	})

	b.Run("500K", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchRes = solve(t, s, v[:500000])
		}
	})

	b.Run("1M", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchRes = solve(t, s, v)
		}
	})
}
