package intesect

import (
	"reflect"
	"sort"
	"testing"
)

func TestSolve(t *testing.T) {
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
			got := Solve(tt.args.ops)
			for _, nums := range got {
				sort.Ints(nums)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Solve() = %v, want %v", got, tt.want)
			}
		})
	}
}
