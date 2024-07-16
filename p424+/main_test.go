package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sort"
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
			"0",
			args{strings.NewReader(`1422`)},
			`1423`,
			false,
			true,
		},
		{
			"1",
			args{strings.NewReader(`999222`)},
			`999999`,
			false,
			true,
		},
		{
			"1.2",
			args{strings.NewReader(`899922`)},
			`899989`,
			false,
			true,
		},
		{
			"1.3",
			args{strings.NewReader(`899999`)},
			`900009`,
			false,
			true,
		},
		// {
		// 	"1.4",
		// 	args{strings.NewReader(`999999`)},
		// 	`000000`,
		// 	false,
		// 	true,
		// },
		// {
		// 	"1.4.2",
		// 	args{strings.NewReader(`9999`)},
		// 	`0000`,
		// 	false,
		// 	true,
		// },
		// {
		// 	"1.4.3",
		// 	args{strings.NewReader(`99`)},
		// 	`00`,
		// 	false,
		// 	true,
		// },
		{
			"1.4",
			args{strings.NewReader(`999999`)},
			`001001`,
			false,
			true,
		},
		{
			"1.4.2",
			args{strings.NewReader(`9999`)},
			`0101`,
			false,
			true,
		},
		{
			"1.4.3",
			args{strings.NewReader(`99`)},
			`11`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2214`)},
			`2222`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1014`)},
			`1102`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`0014`)},
			`0101`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`2294`)},
			`2305`,
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

// func Test_run2(t *testing.T) {
// 	n := 6
// 	f := fmt.Sprintf("%%0%dd", n)

// 	maximum := 1
// 	for i := 0; i < n; i++ {
// 		maximum *= 10
// 	}

// 	for x := 0; x < maximum; x++ {
// 		want := find(x, n)
// 		if want == maximum {
// 			want = find(0, n)
// 		}
// 		num := fmt.Sprintf(f, x)
// 		out := &strings.Builder{}
// 		run(strings.NewReader(num), out)
// 		out_s := strings.TrimSpace(out.String())
// 		got, err := strconv.Atoi(out_s)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if got != want {
// 			t.Fatalf("%s: got = %s("+f+"), want = "+f, num, out_s, got, want)
// 			// } else {
// 			// 	t.Logf(f+": "+f+" "+f, x, got, want)
// 		}
// 	}
// }

func Fuzz_run(t *testing.F) {
	n := 8
	f := fmt.Sprintf("%%0%dd", n)

	maximum := 1
	for i := 0; i < n; i++ {
		maximum *= 10
	}

	r := rand.New(rand.NewSource(1))
	for i := 0; i < 10; i++ {
		t.Add(uint32(r.Int31n(int32(maximum))))
	}

	wants := calcWants(n)

	t.Fuzz(func(t *testing.T, xx uint32) {
		x := int(xx) % maximum
		k := sort.Search(len(wants), func(k int) bool {
			return wants[k] > x
		})
		if k == len(wants) {
			k = 1
		}
		want := wants[k]

		num := fmt.Sprintf(f, x)
		out := &strings.Builder{}
		run(strings.NewReader(num), out)
		out_s := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(out_s)
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Fatalf("%s: got = %s("+f+"), want = "+f, num, out_s, got, want)
		}
	})
}

func find(x int, n int) int {
	x++
	for !check(x, n) {
		x++
	}
	return x
}

func check(x int, n int) bool {
	s1, s2 := 0, 0
	for i := n / 2; x > 0 && i > 0; i-- {
		s1 += x % 10
		x /= 10
	}
	for i := n / 2; x > 0 && i > 0; i-- {
		s2 += x % 10
		x /= 10
	}
	return s1 == s2
}

func calcWants(n int) []int {
	end := 1
	for i := 0; i < n; i++ {
		end *= 10
	}

	var wants []int
	for x := 0; x < end; x++ {
		if check(x, n) {
			wants = append(wants, x)
		}
	}

	return wants
}
