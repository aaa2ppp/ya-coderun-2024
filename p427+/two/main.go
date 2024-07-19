package main

import "log"

// 6 2 3
//    0 1
// 0+ *
// 2+ * x
// 4+ * *

var debugEnable bool

func solve(n, a, b int) int {
	// prepare
	d := gcd(a, b)
	a /= d
	b /= d
	n = (n + d - 1) / d

	if a > b {
		a, b = b, a
	}

	if n <= a {
		return 1
	}

	if debugEnable {
		log.Println("prepared:", n, a, b, d)
	}

	count := 0

	bits := make([]byte, a)
	bits[0] = 1
	bits_count := 1

	b_cur := 0
	b_next := b

	for {
		if debugEnable {
			log.Println(b_cur, b_next, count, bits_count)
		}
		i0 := b_cur / a

		if b_next >= n || bits_count == a {
			i1 := n / a
			count += bits_count * (i1 - i0)
			n %= a
			if bits_count == a {
				count += n
			} else {
				for i := 0; i < n; i++ {
					count += int(bits[i])
				}
			}
			break
		}

		// if bits_count == a {
		// 	i1 := n / a
		// 	count += bits_count * (i1-i0)
		// 	count += n % a
		// 	break
		// }

		i1 := b_next / a
		count += bits_count * (i1 - i0)
		bits[b_next%a] = 1
		bits_count++

		b_cur = b_next
		b_next = b_cur + b
	}

	return count
}

func gcd(a, b int) int {
	for a > 0 {
		a, b = b%a, a
	}
	return b
}
