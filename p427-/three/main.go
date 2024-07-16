package main

import "log"

// 10 3 4 5
//    0 1 2
// 0+ *
// 3+ * x +
// 6+ * * *
// 9+ *

// 10 3 5 7
//    0 1 2
// 0+ *
// 3+ *   x
// 6+ * + *
// 9+ *

var debugEnable bool

func solve(n, a, b, c int) int {
	// prepare

	// вытаскикавем наименьшее число на первую позицию
	if a > b {
		a, b = b, a
	}
	if a > c {
		a, c = c, a
	}

	// выкидываем кратные числа за пределы (>= n), чтобы не мешались
	// делаем их кратными a только для того, чтобы на следующем шаге можно
	// было просто сократить на gcd
	if c%a == 0 || c%b == 0 {
		c = (n + a - 1) / a * a
	} 
	if b%a == 0 {
		b = (n + a - 1) / a * a
		b, c = c, b
	}

	d := gcd(a, gcd(b, c))
	a /= d
	b /= d
	c /= d
	n = (n + d - 1) / d

	if n <= a {
		return 1
	}

	if debugEnable {
		log.Println("prepared:", n, a, b, c, d)
	}

	count := 0

	bits := make([]byte, a)
	bits_count := 0

	cur := 0

	b_next := 0
	q_c_next := newQueue(c + 1)
	q_c_next.push(0)

	for {
		i0 := cur / a
		if debugEnable {
			log.Println(bits, i0, cur, b_next, q_c_next.front(), count, bits_count)
		}

		next := b_next
		if q_c_next.front() <= next {
			next = q_c_next.pop()
		}

		if next >= n || bits_count == a {
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
		// 	count += bits_count * (i1 - i0)
		// 	count += n % a
		// 	break
		// }

		i1 := next / a
		count += bits_count * (i1 - i0)
		if i := next % a; bits[i] == 0 {
			bits[i] = 1
			bits_count++
		}

		cur = next

		if b_next == next {
			b_next += b
		}

		if q_c_next.len() == 0 || q_c_next.front() < n {
			q_c_next.push(next + c)
		}
	}

	return count
}

func gcd(a, b int) int {
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

type queue struct {
	first int
	size  int
	items []int
}

func newQueue(n int) *queue {
	items := make([]int, n)
	return &queue{
		items: items,
	}
}

func (q *queue) len() int {
	return q.size
}

func (q *queue) push(v int) {
	if q.size == len(q.items) {
		panic("queue is full")
	}
	i := q.first + q.size
	if i >= len(q.items) {
		i -= len(q.items)
	}
	q.items[i] = v
	q.size++
}

func (q *queue) front() int {
	if q.size == 0 {
		panic("queue is empty")
	}
	return q.items[q.first]
}

func (q *queue) pop() int {
	v := q.front()
	q.first++
	if q.first >= len(q.items) {
		q.first -= len(q.items)
	}
	q.size--
	return v
}
