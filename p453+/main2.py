import os

debug_enable = False


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
    color = [0]*len(graph)
    loop_starts = {}
    loop_count = 0

    def dfs(node):
        nonlocal loop_count
        stack = []
        color[node] = GREY
        stack.append([node, -1, 0, 0])
    
        while len(stack) > 0:
            node, prev, input_loop_count, neig_idx = stack[-1]

            if neig_idx < len(graph[node]):
                neig = graph[node][neig_idx]
                stack[-1][-1] = neig_idx + 1
                if color[neig] == WHITE:
                    color[neig] = GREY 
                    stack.append([neig, node, loop_count, 0])
                elif color[neig] == GREY and neig != prev:
                    loop_count += 1
                    loop_starts[neig] = loop_starts.get(neig, 0) + 1
                    edges.append(new_edge(node, neig))
                continue

            n = loop_starts.get(node, 0)
            if n != 0:
                loop_count -= n
                del loop_starts[node]

            if loop_count > input_loop_count:
                edges.append(new_edge(prev, node))

            color[node] = BLACK
            stack.pop()
                                 
    for node in range(len(graph)):
        if color[node] == WHITE:
            dfs(node)

    return edges


def devide_by_comps(graph):
    visited = [False]*len(graph)

    def dfs(node):
        size = 0
        stack = []
        visited[node] = True
        stack.append([node, 0])

        while len(stack) > 0:
            node, neig_idx = stack[-1]

            if neig_idx < len(graph[node]): 
                neig = graph[node][neig_idx]
                stack[-1][-1] = neig_idx + 1
                if not visited[neig]:
                    visited[neig] = True
                    stack.append([neig, 0])
                continue

            size += 1
            stack.pop()

        return size

    comps = []    
    for node in range(len(graph)):
        if not visited[node]:
            size = dfs(node)
            if size >= 3:
                comps.append((node, size))
    
    return comps


def count_edges_to_looping(graph, node, size):
    count = 0
    stack = []
    stack.append([node, -1, 0])
    while len(stack) > 0:
        node, prev, neig_idx = stack[-1]
        if neig_idx < len(graph[node]):
            neig = graph[node][neig_idx]
            stack[-1][-1] = neig_idx + 1
            if neig != prev:
                stack.append([neig, node, 0])
            continue
        count += size - len(graph[node]) - 1
        stack.pop()
    return count // 2


if __name__ == '__main__':
    debug_enable = "DEBUG" in os.environ
    main()