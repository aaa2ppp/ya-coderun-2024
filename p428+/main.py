def main():
    t = input().strip()
    s = input().strip()
    res = solve(t, s)
    print(len(res))
    print(*res)


X = 269
MODULO = int(1e9) + 7

def solve(t, s):
    n = len(s)

    pow = [0]*(n+1)
    pow[0] = 1
    for i in range(1, len(pow)):
        pow[i] = pow[i-1] * X % MODULO

    pow_n = pow[n]
    s_pow = make_pattern(s, pow)

    t_hash = [0]*(len(t)+1)
    for i, c in enumerate(t):
        t_hash[i+1] = (t_hash[i] * X + ord(c)) % MODULO

    res = []
    for i in range(0, len(t)-len(s)+1):
        s_hash = 0
        for p in s_pow:
            s_hash = (s_hash + ord(t[i+p[0]]) * p[1]) % MODULO
        t_h1 = t_hash[i] * pow_n #% MODULO
        t_h2 = t_hash[i+n]
        if t_h2 == ((s_hash + t_h1) % MODULO) and check_unique(t, i, s_pow):
            res.append(i+1)
    return res


def check_unique(t, i, s_pow):
    l_set = [0, 0]
    for c in (ord(t[i+p[0]]) for p in s_pow): 
        if l_set[c>>6] & (1 << (c&63)):
            return False
        l_set[c>>6] |= 1 << (c&63)
    return True


def make_pattern(s, pow):
    n = len(s)
    s_pow_set = [0]*128
    s_pos_set = [-1]*128
    for i in range(len(s)):
        c = ord(s[i])
        if s_pos_set[c] == -1:
            s_pos_set[c] = i
        s_pow_set[c] = (s_pow_set[c] + pow[n-i-1]) % MODULO
    s_pow = []
    for i in range(128):
        if s_pos_set[i] != -1:
            s_pow.append((s_pos_set[i], s_pow_set[i]))
    return s_pow
	

if __name__ == '__main__':
    main()