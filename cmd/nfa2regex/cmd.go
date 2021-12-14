package main

import (
	"fmt"
	"net/url"
	"strconv"

	n "github.com/wolever/nfa2regex"
)

// MakeNFASimple generates a simple example NFA
func MakeNFASimple() *n.NFA {
	nfa := n.New()
	nfa.AddEdge("1", "1", "a")
	nfa.AddEdge("1", "2", "b")
	nfa.AddEdge("2", "2", "c")
	nfa.AddEdge("2", "3", "d")
	nfa.AddEdge("3", "3", "e")
	nfa.AddEdge("3", "1", "x")
	nfa.Nodes["1"].IsInitial = true
	nfa.Nodes["3"].IsTerminal = true
	return nfa
}

// MakeNFAManyMany genreates an NFA with multiple initial and terminal nodes
func MakeNFAManyMany() *n.NFA {
	nfa := n.New()
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

// MakeNFAMultiplesOfX generates an NFA which will match binary strings that
// are multiples of ``x``
func MakeNFAMultiplesOfX(x int) *n.NFA {
	str := func(i int) string { return strconv.Itoa((i)) }
	nfa := n.New()
	for i := 0; i < x; i += 1 {
		nfa.AddEdge(str(i), str((i*2)%x), "0")
		nfa.AddEdge(str(i), str((i*2+1)%x), "1")
	}
	nfa.AddEdge("start", "0", "0")
	nfa.AddEdge("start", "1", "1")
	nfa.Nodes["start"].IsInitial = true
	nfa.Nodes["0"].IsTerminal = true
	return nfa
}

func main() {
	nfa := MakeNFAManyMany()
	//nfa := MakeNFAMultiplesOfN(3)
	regex := n.NFA2Regex(nfa)
	fmt.Println("Graph:")
	fmt.Println("https://dreampuf.github.io/GraphvizOnline/#" + url.PathEscape(n.NFA2Dot(nfa)))
	fmt.Println()
	fmt.Println("Regex:")
	fmt.Println(regex)
}
