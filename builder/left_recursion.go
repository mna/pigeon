package builder

import (
	"errors"
	"fmt"

	"github.com/mna/pigeon/ast"
)

var (
	// ErrNoLeader is no leader error.
	ErrNoLeader = errors.New(
		"SCC has no leadership candidate (no element is included in all cycles)")
	// ErrHaveLeftRecirsion is recursion error.
	ErrHaveLeftRecirsion = errors.New("have left recursion")
)

// PrepareGramma evaluates parameters associated with left recursion
func PrepareGramma(grammar *ast.Grammar) error {
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	ComputeNullables(mapRules)
	if err := ComputeLeftRecursives(mapRules); err != nil {
		return fmt.Errorf("error compute left recursive: %w", err)
	}
	rulesWithLeftRecursion := []string{}
	for _, rule := range grammar.Rules {
		if rule.LeftRecursive {
			rulesWithLeftRecursion = append(rulesWithLeftRecursion, rule.Name.Val)
		}
	}
	if len(rulesWithLeftRecursion) > 0 {
		return fmt.Errorf("%w: %v", ErrHaveLeftRecirsion, rulesWithLeftRecursion)
	}

	return nil
}

// ComputeNullables evaluates nullable nodes
func ComputeNullables(rules map[string]*ast.Rule) {
	// Compute which rules in a grammar are nullable
	for _, rule := range rules {
		rule.NullableVisit(rules)
	}
}

func findLeader(
	graph map[string]map[string]bool, scc map[string]bool,
) (string, error) {
	// Try to find a leader such that all cycles go through it.
	leaders := make(map[string]bool, len(scc))
	for k := range scc {
		leaders[k] = true
	}
	for start := range scc {
		for _, cycle := range FindCyclesInSCC(graph, scc, start) {
			mapCycle := map[string]bool{}
			for _, k := range cycle {
				mapCycle[k] = true
			}
			for k := range scc {
				if _, okCycle := mapCycle[k]; !okCycle {
					delete(leaders, k)
				}
			}
			if len(leaders) == 0 {
				return "", ErrNoLeader
			}
		}
	}
	// Pick an arbitrary leader from the candidates.
	var leader string
	for k := range leaders {
		leader = k // The only element.
		break
	}
	return leader, nil
}

// ComputeLeftRecursives evaluates left recursion
func ComputeLeftRecursives(rules map[string]*ast.Rule) error {
	graph := MakeFirstGraph(rules)
	vertices := make([]string, 0, len(graph))
	for k := range graph {
		vertices = append(vertices, k)
	}
	sccs := StronglyConnectedComponents(vertices, graph)
	for _, scc := range sccs {
		if len(scc) > 1 {
			for name := range scc {
				rules[name].LeftRecursive = true
			}
			leader, err := findLeader(graph, scc)
			if err != nil {
				return fmt.Errorf("error find leader %v: %w", scc, err)
			}
			rules[leader].Leader = true
		} else {
			var name string
			for k := range scc {
				name = k // The only element.
				break
			}
			if _, ok := graph[name][name]; ok {
				rules[name].LeftRecursive = true
				rules[name].Leader = true
			}
		}
	}
	return nil
}

// MakeFirstGraph compute the graph of left-invocations.
// There's an edge from A to B if A may invoke B at its initial position.
// Note that this requires the nullable flags to have been computed.
func MakeFirstGraph(rules map[string]*ast.Rule) map[string]map[string]bool {
	graph := make(map[string]map[string]bool)
	vertices := make(map[string]bool)
	for rulename, rule := range rules {
		names := rule.InitialNames()
		graph[rulename] = names
		for name := range names {
			vertices[name] = true
		}
	}
	for vertex := range vertices {
		if _, ok := graph[vertex]; !ok {
			graph[vertex] = make(map[string]bool)
		}
	}
	return graph
}
