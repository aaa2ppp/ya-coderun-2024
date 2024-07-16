package permut

import "sort"

// Находить все возможные суммы nums < len(set). Nums должен быть отсортирован.
// Отмечает найденные суммы в set. Возвращает неупорядоченный список
// найденных сумм.
func permut(nums []int, set []bool) []int {
	list := make([]int, 0, len(nums)*2)
	list = append(list, 0)
	for _, v := range nums {
		for i, n := 0, len(list); i < n; i++ {
			if v2 := list[i] + v; v2 < len(set) && !set[v2] {
				set[v2] = true
				list = append(list, v2)
			}
		}
	}
	return list[1:]
}

func solve(n int, nums []int) []int {
	set := make([]bool, n+1)
	res := permut(nums, set)
	sort.Ints(res)
	return res
}
