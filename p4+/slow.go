package main

func solveOne(t, s int, v []int) int {
	var count int

	s0 := v[0] * t
	s0_div_s := s0 / s
	s0_mod_s := s0 % s

	for i := 1; i < len(v); i++ {
		if v[0] <= v[i] {
			continue
		}

		si := v[i] * t

		count += s0_div_s - si/s - 1
		if s0_mod_s > si%s {
			count++
		}
	}

	return count
}

func slowSolve(t, s int, v []int) int {
	count := 0
	for k := range v {
		v[0], v[k] = v[k], v[0]
		count += solveOne(t, s, v)
		v[0], v[k] = v[k], v[0]
	} 
	return count
}
