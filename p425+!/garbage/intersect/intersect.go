package intesect

import (
	"sort"
)

func Intersect() {}

type Op struct {
	l, r, x int
}

func Solve(ops []Op) [][]int {
	var res [][]int

	type item struct {
		*Op
		used bool
	}

	list := make([]item, len(ops))
	for i := range ops {
		list[i] = item{Op: &ops[i]}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].l < list[j].l || list[i].l == list[j].l && list[i].r < list[j].r
	})

	l, r := 0, 1
	for l < len(list) {
		var nums []int

		for r < len(list) && list[r].l <= list[l].r {
			nums = append(nums, list[r].x)
			list[r].used = true
			r++
		}

		if !list[l].used || nums != nil {
			nums = append(nums, list[l].x)
			res = append(res, nums)
		}

		l++
	}

	return res
}
