def solve(t, s, v):
    count = 0

    v.sort()

    s0 = v[0] * t
    s0_i = 0
    s0_div_s = s0 // s
    s0_mod_s = s0 % s

    sum_div_s, si_count = 0, 0

    nodes = [Node() for _ in range(len(v))]
    tree = Node.insert(None, s0_mod_s, nodes)

    for i in range(1, len(v)):
        si = v[i] * t

        if si != s0:
            sum_div_s += s0_div_s * (i-s0_i)
            si_div_s = si // s
            si_count = (si_div_s-1)*i - sum_div_s
            s0 = si
            s0_i = i
            s0_div_s = si_div_s

        si_mod_s = si % s
        si_last = Node.find_idx(tree, si_mod_s)
        tree = Node.insert(tree, si_mod_s, nodes)
        s0_mod_s = si_mod_s

        count += si_count + si_last

    return count


class Node:
    def __init__(self):
        self._key = 0
        self._left = None
        self._right = None
        self._size = 1
        self._hght = 1

    @staticmethod
    def size(node):
        return 0 if node is None else node._size

    @staticmethod
    def height(node):
        return 0 if node is None else node._hght

    @staticmethod
    def find_idx(node, key):
        '''Возвращает индекс (начиная с 0) первого не меньшиго по ключу узла.'''
        '''Если такого ключа не существует, возвращает node.size()'''
        idx = 0
        while node is not None:
            if key <= node._key:
                node = node._left
            else:
                idx += Node.size(node._left) + 1
                node = node._right
        return idx

    @staticmethod
    def insert(node, key, nodes):
        if node is None:
            node = nodes.pop()
            node._key = key
            return node
        if key <= node._key:
            new_left = Node.insert(node._left, key, nodes)
            node._left = new_left
            return Node._repair(node)
        else:
            new_right = Node.insert(node._right, key, nodes)
            node._right = new_right
            return Node._repair(node)

    @staticmethod
    def _update(node):
        node._size = Node.size(node._left) + Node.size(node._right) + 1
        node._hght = max(Node.height(node._left), Node.height(node._right)) + 1

    @staticmethod
    def _repair(node):
        d = Node.height(node._left) - Node.height(node._right)
        if d < -1:
            return Node._left_rotate(node)
        elif d > 1:
            return Node._right_rotate(node)
        else:
            Node._update(node)
            return node

    @staticmethod
    def _left_rotate(node):
        al = node
        bt = al._right

        if Node.height(bt._right) - Node.height(bt._left) > 0:
            al._right = bt._left
            Node._update(al)
            bt._left = al
            Node._update(bt)
            return bt
        else:
            ga = bt._left
            al._right = ga._left
            Node._update(al)
            bt._left = ga._right
            Node._update(bt)
            ga._left = al
            ga._right = bt
            Node._update(ga)
            return ga

    @staticmethod
    def _right_rotate(node):
        al = node
        bt = al._left

        if Node.height(bt._left) - Node.height(bt._right) > 0:
            al._left = bt._right
            Node._update(al)
            bt._right = al
            Node._update(bt)
            return bt
        else:
            ga = bt._right
            al._left = ga._right
            Node._update(al)
            bt._right = ga._left
            Node._update(bt)
            ga._right = al
            ga._left = bt
            Node._update(ga)
            return ga


def main():
    n, t, s = map(int, input().split())
    v = list(map(int, input().split()))
    res = solve(t, s, v)
    print(res)


if __name__ == '__main__':
    main()
