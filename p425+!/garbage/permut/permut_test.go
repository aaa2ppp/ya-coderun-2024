package permut

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestSolve(t *testing.T) {
	type args struct {
		n    int
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"1",
			args{5, []int{1, 3, 5}},
			[]int{1, 3, 4, 5},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve(tt.args.n, tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Solve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_solve_1000(t *testing.B) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1000
	q := n / 2
	nums := make([]int, q)

	for i := range nums {
		nums[i] = rd.Intn(n) + 1
	}

	for i := 0; i < t.N; i++ {
		solve(n, nums)
	}
}

func Benchmark_solve_10000(t *testing.B) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 10000
	q := n
	nums := make([]int, q)

	for i := range nums {
		nums[i] = rd.Intn(n) + 1
	}

	for i := 0; i < t.N; i++ {
		solve(n, nums)
	}
}

func Benchmark_solve_100000(t *testing.B) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 100000
	q := n / 2
	nums := make([]int, q)

	for i := range nums {
		nums[i] = rd.Intn(n) + 1
	}

	for i := 0; i < t.N; i++ {
		solve(n, nums)
	}
}
