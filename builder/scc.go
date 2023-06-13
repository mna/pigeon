package builder

import "fmt"

func min(a1 int, a2 int) int {
	if a1 <= a2 {
		return a1
	}
	return a2
}

// StronglyConnectedComponents compute strongly сonnected сomponents of a graph.
// Tarjan's strongly connected components algorithm
func StronglyConnectedComponents(
	vertices []string, edges map[string]map[string]bool,
) []map[string]bool {
	// Tarjan's strongly connected components algorithm
	var (
		identified = map[string]bool{}
		stack      = []string{}
		index      = map[string]int{}
		lowlink    = map[string]int{}
		dfs        func(v string) []map[string]bool
	)

	dfs = func(vertex string) []map[string]bool {
		index[vertex] = len(stack)
		stack = append(stack, vertex)
		lowlink[vertex] = index[vertex]

		sccs := []map[string]bool{}
		for w := range edges[vertex] {
			if _, ok := index[w]; !ok {
				sccs = append(sccs, dfs(w)...)
				lowlink[vertex] = min(lowlink[vertex], lowlink[w])
			} else if _, ok := identified[w]; !ok {
				lowlink[vertex] = min(lowlink[vertex], lowlink[w])
			}
		}

		if lowlink[vertex] == index[vertex] {
			scc := map[string]bool{}
			for _, v := range stack[index[vertex]:] {
				scc[v] = true
			}
			stack = stack[:index[vertex]]
			for v := range scc {
				identified[v] = true
			}
			sccs = append(sccs, scc)
		}
		return sccs
	}

	sccs := []map[string]bool{}
	for _, v := range vertices {
		if _, ok := index[v]; !ok {
			sccs = append(sccs, dfs(v)...)
		}
	}
	return sccs
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func reduceGraph(
	graph map[string]map[string]bool, scc map[string]bool,
) map[string]map[string]bool {
	reduceGraph := map[string]map[string]bool{}
	for src, dsts := range graph {
		if _, ok := scc[src]; !ok {
			continue
		}
		reduceGraph[src] = map[string]bool{}
		for dst := range dsts {
			if _, ok := scc[dst]; !ok {
				continue
			}
			reduceGraph[src][dst] = true
		}
	}
	return reduceGraph
}

// FindCyclesInSCC find cycles in SCC emanating from start.
// Yields lists of the form ['A', 'B', 'C', 'A'], which means there's
// a path from A -> B -> C -> A.  The first item is always the start
// argument, but the last item may be another element, e.g.  ['A',
// 'B', 'C', 'B'] means there's a path from A to B and there's a
// cycle from B to C and back.
func FindCyclesInSCC(
	graph map[string]map[string]bool, scc map[string]bool, start string,
) [][]string {
	// Basic input checks.
	if _, ok := scc[start]; !ok {
		panic(fmt.Sprintf("scc %v have not %v", scc, start))
	}
	extravertices := []string{}
	for k := range scc {
		if _, ok := graph[k]; !ok {
			extravertices = append(extravertices, k)
		}
	}
	if len(extravertices) != 0 {
		panic(fmt.Sprintf("graph have not scc. %v", extravertices))
	}

	// Reduce the graph to nodes in the SCC.
	graph = reduceGraph(graph, scc)
	if _, ok := graph[start]; !ok {
		panic(fmt.Sprintf("graph %v have not %v", graph, start))
	}

	// Recursive helper that yields cycles.
	var dfs func(node string, path []string) [][]string
	dfs = func(node string, path []string) [][]string {
		ret := [][]string{}
		if contains(path, node) {
			t := make([]string, 0, len(path)+1)
			t = append(t, path...)
			t = append(t, node)
			ret = append(ret, t)
			return ret
		}
		path = append(path, node) // TODO: Make this not quadratic.
		for child := range graph[node] {
			ret = append(ret, dfs(child, path)...)
		}
		return ret
	}

	return dfs(start, []string{})
}
