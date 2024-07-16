package main

import (
	"bytes"
	"io"
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
			args{strings.NewReader(`4 103 123 20 4567`)},
			`3`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5 22 43 55 42 4`)},
			`4`,
			false,
			true,
		},
		// [7887 1847] => 1
		// [9786 9692 2040 6104 3854 1730 2964] => 17
		// [3977 4896 1850 7777 2779 5916 7517 1144 747] => 25
		// [3156 1013 175 6937 1197 3924 6813 711 5296 1417] => 40
		// [3836 4489 8480 4387 2266 5086 149] => 15
		{
			"3",
			args{strings.NewReader(`2 7887 1847`)},
			`1`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`7 9786 9692 2040 6104 3854 1730 2964`)},
			`17`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`9 3977 4896 1850 7777 2779 5916 7517 1144 747`)},
			`25`,
			false,
			true,
		},
		{
			"6",
			args{strings.NewReader(`10 3156 1013 175 6937 1197 3924 6813 711 5296 1417`)},
			`40`,
			false,
			true,
		},
		{
			"7",
			args{strings.NewReader(`7 3836 4489 8480 4387 2266 5086 149`)},
			`15`,
			false,
			true,
		},
		{
			"8",
			args{strings.NewReader(`7 3836 3836 3836 3836 3836 3836 3836`)},
			`21`,
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

// func Fuzz_countPairs(t *testing.F) {
// 	for i := int64(1); i <= 5; i++ {
// 		t.Add(i)
// 	}

// 	t.Fuzz(func(t *testing.T, seed int64) {
// 		r := rand.New(rand.NewSource(seed))
// 		n := r.Intn(10) + 1
// 		nums := make([]int, n)
// 		for i := range nums {
// 			nums[i] = r.Intn(10000)
// 		}
// 		t.Log(nums)
// 		digs := prepareDigs(nums)
// 		want := bruteforceCountPairs(digs)
// 		if got := countPairs(digs); got != want {
// 			t.Errorf("countPairs() = %v, want %v", got, want)
// 		} else {
// 			t.Logf("countPairs() = %v", got)
// 		}
// 	})
// }

// func bruteforceCountPairs(digs []int16) int {
// 	n := len(digs)
// 	pairs := make(map[[2]int]struct{}, n*(n-1)/2)
// 	for i := 0; i < n; i++ {
// 		for j := i + 1; j < n; j++ {
// 			if digs[i]&digs[j] != 0 {
// 				pairs[[2]int{i, j}] = struct{}{}
// 			}
// 		}
// 	}
// 	return len(pairs)
// }

// func Test_bruteforceCountPairs(t *testing.T) {
// 	tests := []struct {
// 		nums []int
// 		want int
// 	}{
// 		{
// 			[]int{103, 123, 20, 4567},
// 			3,
// 		},
// 		{
// 			[]int{22, 43, 55, 42, 4},
// 			4,
// 		},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprint(tt.nums), func(t *testing.T) {
// 			if got := bruteforceCountPairs(prepareDigs(tt.nums)); got != tt.want {
// 				t.Errorf("bruteforceCountPairs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
