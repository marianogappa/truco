package newbot

import (
	"errors"

	"github.com/marianogappa/truco/truco"
)

var (
	rules    []rule
	_ruleset = map[string]struct{}{}
)

// Every rule must be registered in order to be considered by the bot.
func registerRule(r rule) {
	if _, ok := _ruleset[r.name]; ok {
		panic("rule already registered; is the name unique?")
	}
	rules = append(rules, r)
	_ruleset[r.name] = struct{}{}
}

type state map[string]any
type stateChange struct {
	fn          func(*state)
	description string
}

type ruleResult struct {
	action            truco.Action
	stateChanges      []stateChange
	resultDescription string
}

type rule struct {
	name         string
	description  string
	isApplicable func(state, truco.ClientGameState) bool
	dependsOn    []rule
	run          func(state, truco.ClientGameState) (ruleResult, error)
}

func topologicalSortKahn(rules []rule) ([]rule, error) {
	nameToRule := map[string]rule{}
	for _, r := range rules {
		nameToRule[r.name] = r
	}

	indegree := map[string]int{}
	graph := map[string][]rule{}
	for _, r := range rules {
		indegree[r.name] = 0
		graph[r.name] = []rule{}
	}
	for _, r := range rules {
		for _, d := range r.dependsOn {
			graph[d.name] = append(graph[d.name], r)
			indegree[r.name]++
		}
	}
	queue := []rule{}
	for r, d := range indegree {
		if d == 0 {
			queue = append(queue, nameToRule[r])
		}
	}
	sorted := []rule{}
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]
		sorted = append(sorted, r)
		for _, d := range graph[r.name] {
			indegree[d.name]--
			if indegree[d.name] == 0 {
				queue = append(queue, d)
			}
		}
	}
	if len(sorted) != len(rules) {
		return nil, errors.New("circular dependency")
	}
	return sorted, nil
}
