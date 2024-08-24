import os
import sys

debug_enable = False
sys.setrecursionlimit(101000)


def main():
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        a, b = map(int, input().split())
        edges.append(new_edge(a,b))
    if debug_enable:
        print("edges:", edges)

    graph = make_graph(n, edges)
    if debug_enable:
        print("graph:", graph)

    loop_edges = find_loop_edges(graph)
    if debug_enable:
        print("loop_edges:", loop_edges)

    exclude = set(loop_edges)
    if debug_enable:
        print("exclude:", exclude)

    edges2 = [edge for edge in edges if edge not in exclude]
    if debug_enable:
        print("edges2:", edges2)

    graph2 = make_graph(n, edges2)
    if debug_enable:
        print("graph2:", graph2)

    comps = devide_by_comps(graph2)
    if debug_enable:
        print("comps:", comps)

    total = 0
    for node, size in comps:
        total += count_edges_to_looping(graph2, node, size)

    print(total)


def new_edge(a, b):
    if a > b:
        a, b = b, a
    return a, b


def make_graph(n, edges):
    graph = [[] for _ in range(n+1)]
    for edge in edges:
        a, b = edge
        graph[a].append(b)
        graph[b].append(a)
    return graph
    

WHITE = 0
GREY  = 1
BLACK = 2 

def find_loop_edges(graph):
    edges = []
    visited = [0]*len(graph)
    loop_starts = {}
    loop_count = 0

    def dfs(node, prev):
        nonlocal loop_count

        if visited[node] == BLACK:
            return False
        
        if visited[node] == GREY:
            loop_count += 1
            loop_starts[node] = loop_starts.get(node, 0) + 1
            return True
        
        visited[node] = GREY
        input_loop_count = loop_count
        in_loop = False

        for neig in graph[node]:
            if neig == prev:
                continue

            if dfs(neig, node):
                in_loop = True
                edges.append(new_edge(node, neig))

        visited[node] = BLACK

        n = loop_starts.get(node, 0)
        if n != 0:
            loop_count -= n
            del loop_starts[node]

        return in_loop and loop_count > input_loop_count

    for node in range(len(graph)):
        if visited[node] == WHITE:
            dfs(node, -1)

    return edges


def devide_by_comps(graph):
    visited = [False]*len(graph)

    def dfs(node):
        if visited[node]:
            return 0
        visited[node] = True
        size = 1
        for neig in graph[node]:
            size += dfs(neig)
        return size

    comps = []    
    for node in range(len(graph)):
        if not visited[node]:
            size = dfs(node)
            if size >= 3:
                comps.append((node, size))
    
    return comps


def count_edges_to_looping(graph, node, size):
    def dfs(node, prev):
        count = size - len(graph[node]) - 1
        for neig in graph[node]:
            if neig != prev:
                count += dfs(neig, node)
        return count
    return dfs(node, -1) // 2



if __name__ == '__main__':
    debug_enable = "DEBUG" in os.environ
    main()