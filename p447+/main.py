import math


MODULO = 1_000_000_000


def solve(a, b):
	if len(a) > len(b):
		a, b = b, a

	d = 1
	truncated = False

	for i in range(len(a)):
		for j in range(len(b)):
			v = math.gcd(a[i], b[j])
			if v > 1:
				d *= v
				if d >= MODULO:
					d %= MODULO
					truncated = True
					if d == 0:
						return 0, True
				a[i] //= v
				b[j] //= v

	return d, truncated


def main():
	n = int(input())
	a = list(map(int, input().split()))
	m = int(input())
	b = list(map(int, input().split()))

	d, truncated = solve(a, b)
	if truncated:
		print(f"{d:09}")
	else:
		print(d)


if __name__ == '__main__':
    main()