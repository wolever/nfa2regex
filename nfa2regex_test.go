package nfa2regex

import (
	"fmt"
	"regexp"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MakeNFAManyMany() *NFA {
	nfa := NewNFA()
	nfa.AddEdge("1", "2", "a")
	nfa.AddEdge("2", "3", "b")
	nfa.AddEdge("2", "2", "l")
	nfa.AddEdge("4", "2", "x")
	nfa.AddEdge("2", "5", "y")
	nfa.Nodes["1"].IsInitial = true
	nfa.Nodes["4"].IsInitial = true
	nfa.Nodes["3"].IsTerminal = true
	nfa.Nodes["5"].IsTerminal = true
	return nfa
}

func TestNFAMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "ab", expected: true},
		{input: "alllb", expected: true},
		{input: "aa", expected: false},
		{input: "xy", expected: true},
		{input: "fff", expected: false},
		{input: "ax", expected: false},
	}

	nfa := MakeNFAManyMany()

	for _, tc := range tests {
		testName := fmt.Sprintf("%s:%t", tc.input, tc.expected)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, nfa.Match(tc.input), tc.expected, "Match failed")
		})
	}
}

// nfaWalk walks every path through an NFA, calling `callback` at each node, up
// to a maximum depth of `depth`
func nfaWalk(nfa *NFA, maxDepth int, callback func(path string, node *NFANode)) {
	for _, node := range nfa.Nodes {
		if node.IsInitial {
			nfaWalkRecurse(nfa, "", node, maxDepth, callback)
		}
	}
}
func nfaWalkRecurse(
	nfa *NFA,
	path string,
	node *NFANode,
	depth int,
	callback func(path string, node *NFANode),
) {
	callback(path, node)
	if depth == 0 {
		return
	}

	for _, edgeOut := range nfa.EdgesOut(node.Name) {
		nfaWalkRecurse(nfa, path+edgeOut.Value, edgeOut.DstNode, depth-1, callback)
	}
}

func TestNFAWalk(t *testing.T) {
	nfa := MakeNFAManyMany()
	expected := []string{
		":false",
		"a:false",
		"ab:true",
		"al:false",
		"alb:true",
		"all:false",
		"aly:true",
		"ay:true",
		":false",
		"x:false",
		"xb:true",
		"xl:false",
		"xlb:true",
		"xll:false",
		"xly:true",
		"xy:true",
	}
	actual := make([]string, 0, 10)
	nfaWalk(nfa, 3, func(path string, node *NFANode) {
		actual = append(actual, fmt.Sprintf("%s:%t", path, node.IsTerminal))
	})
	sort.Strings(expected)
	sort.Strings(actual)
	assert.Equal(t, expected, actual)
}

func TestGeneratedRegex(t *testing.T) {
	nfa := MakeNFAManyMany()
	regexStr := ToRegex(nfa)
	regex := regexp.MustCompile(regexStr)

	nfaWalk(nfa, 6, func(path string, node *NFANode) {
		t.Run(fmt.Sprintf("%s:%t", path, node.IsTerminal), func(t *testing.T) {
			isMatch := regex.MatchString(path)
			assert.Equal(t, node.IsTerminal, isMatch)
		})
	})
}
