import math

def solve(n, a, b, c):
	n, a, b, c = prepare(n, a, b, c)

	row = [False]*a
	row_count = 0

	prev_i = 0
	count = 0

	def put(val):
		nonlocal prev_i, row_count, count
		
		i = val // a
		j = val % a

		if row[j]:
			return False

		row[j] = True

		count += (i-prev_i)*row_count + 1
		row_count += 1
		prev_i = i

		return True

	b_queue = []

	put(0)
	c_next = c
	b_queue.append(b)

	while row_count < len(row) and (len(b_queue) > 0 and b_queue[0] <= n or c_next < n):
		if len(b_queue) > 0 and b_queue[0] <= c_next:
			v = b_queue.pop(0)
			if put(v):
				b_queue.append(v + b)
			continue

		v = c_next
		c_next += c

		if put(v):
			b_queue.append(v + b)

	count += (n//len(row) - prev_i) * row_count
	for j in range(n % len(row), len(row)):
		if row[j]:
			count -= 1

	return count


def prepare(n, a, b, c):
	d = math.gcd(a, math.gcd(b, c))
	a //= d
	b //= d
	c //= d
	n = (n + d - 1) // d

	if a > b:
		a, b = b, a
	if a > c:
		a, c = c, a
	if b > c:
		b, c = c, b

	return n, a, b, c


def main():
	n = int(input())
	a, b, c = map(int, input().split())
	res = solve(n, a, b, c)
	print(res)


if __name__ == '__main__':
    main()