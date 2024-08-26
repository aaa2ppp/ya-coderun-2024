package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/rand"
	"reflect"
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
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`3 3
1 7
10 5
8 9
3 0
3 1 1
6 2 1 2`)},
			`11
34
13`,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// TODO: Add test cases.
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

type args struct {
	w []warehouse
	o []order
}

func Fuzz_solve10(f *testing.F) {
	const N = 100

	tests := []args{
		{
			[]warehouse{{1, 7}, {10, 5}, {8, 9}},
			[]order{{3, nil}, {3, []int{1}}, {6, []int{1, 2}}},
		},
	}

	for i := range tests {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, a int) {
		var tt args
		if 0 <= a && a < len(tests) {
			tt = tests[a]
		} else {
			rand := rand.New(rand.NewSource(int64(a)))
			n := 10
			m := 10
			k := n

			w := make([]warehouse, n)
			for i := range w {
				x := rand.Intn(N)
				price := rand.Intn(N) + 1
				w[i] = warehouse{x, price}
			}
			sort.Slice(w, func(i, j int) bool {
				return w[i].x < w[j].x
			})

			banned := make([]int, k)
			for i := range banned {
				banned[i] = i + 1
			}
			rand.Shuffle(len(banned), func(i, j int) {
				banned[i], banned[j] = banned[j], banned[i]
			})

			o := make([]order, m)
			for i := range o {
				o[i].x = rand.Intn(N)
				for j, k := rand.Intn(len(banned)), rand.Intn(n/2); k > 0; j, k = (j+1)%len(banned), k-1 {
					o[i].banned = append(o[i].banned, banned[j])
				}
			}
			sort.Slice(o, func(i, j int) bool {
				return o[i].x < o[j].x
			})

			tt = args{w, o}
			log.Println("tt:", tt)
		}

		want := slowSolve(tt.w, tt.o)
		defer func(v bool) { debugEnable = v }(debugEnable)
		debugEnable = true
		got := solve(tt.w, tt.o)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

var tt args

func init() {
	rand := rand.New(rand.NewSource(1))
	n := 200000
	m := 200000

	w := make([]warehouse, n)
	for i := range w {
		x := rand.Intn(1e9)
		price := rand.Intn(1e9) + 1
		w[i] = warehouse{x, price}
	}

	banned := make([]int, n)
	for i := range banned {
		banned[i] = i + 1
	}
	rand.Shuffle(len(banned), func(i, j int) {
		banned[i], banned[j] = banned[j], banned[i]
	})

	o := make([]order, m)
	for i := range o {
		o[i].x = rand.Intn(1e9)
		for j, n := rand.Intn(len(banned)), rand.Intn(100); n > 0; j, n = (j+1)%len(banned), n-1 {
			o[i].banned = append(o[i].banned, banned[j])
		}
	}
}

var res []int

func Benchmark_solve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res = solve(tt.w, tt.o)
		bw := bufio.NewWriter(io.Discard)
		defer bw.Flush()
		writeInts(bw, res, writeOpts{})
	}
}
