package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type NFANodeName = string
type NFAEdgeValue = string

type NFA struct {
	nodes map[NFANodeName](*NFANode)
	edges [](*NFAEdge)
}

type NFAEdge struct {
	srcNode *NFANode
	dstNode *NFANode
	value   string
}

type NFANode struct {
	name       NFANodeName
	isInitial  bool
	isTerminal bool
}

// Adds an edge between nodes ``srcName`` and ``dstName`` with ``value`` to the NFA.
func (nfa *NFA) addEdge(srcName NFANodeName, dstName NFANodeName, value NFAEdgeValue) {
	srcNode := nfa._getOrCreateNode(srcName)
	dstNode := nfa._getOrCreateNode(dstName)
	nfa.edges = append(nfa.edges, &NFAEdge{
		srcNode: srcNode,
		dstNode: dstNode,
		value:   value,
	})
}

func (nfa *NFA) replaceNode(name NFANodeName, newNode *NFANode) {
	oldNode := nfa.nodes[name]
	nfa.nodes[name] = newNode
	for _, edge := range nfa.edges {
		if edge.srcNode == oldNode {
			edge.srcNode = newNode
		}
		if edge.dstNode == oldNode {
			edge.dstNode = newNode
		}
	}
}

// Removes a node an all associated edges from the NFA.
func (nfa *NFA) removeNode(nodeName NFANodeName) {
	node := nfa.nodes[nodeName]
	delete(nfa.nodes, nodeName)
	newEdges := make([](*NFAEdge), 0, len(nfa.edges))
	for _, edge := range nfa.edges {
		if edge.srcNode == node || edge.dstNode == node {
			continue
		}
		newEdges = append(newEdges, edge)
	}
	nfa.edges = newEdges
}

// Returns a list of edges into ``nodeName`` (ie, where edge.dstNode == nodeName)
func (nfa *NFA) edgesIn(nodeName NFANodeName) [](*NFAEdge) {
	node := nfa.nodes[nodeName]
	res := [](*NFAEdge){}
	for _, edge := range nfa.edges {
		if edge.dstNode == node {
			res = append(res, edge)
		}
	}
	return res
}

// Returns a list of edges out from ``nodeName`` (ie, where edge.srcNode == nodeName)
func (nfa *NFA) edgesOut(nodeName NFANodeName) [](*NFAEdge) {
	node := nfa.nodes[nodeName]
	res := [](*NFAEdge){}
	for _, edge := range nfa.edges {
		if edge.srcNode == node {
			res = append(res, edge)
		}
	}
	return res
}

// Internal: gets node ``name``, or creates if it does not exist.
func (nfa *NFA) _getOrCreateNode(name NFANodeName) *NFANode {
	node := nfa.nodes[name]
	if node == nil {
		node = &NFANode{
			name:       name,
			isInitial:  false,
			isTerminal: false,
		}
		nfa.nodes[name] = node
	}
	return node
}

// Creates a shallow copy of the NFA
func (nfa *NFA) shallowCopy() *NFA {
	res := &NFA{
		nodes: map[NFANodeName]*NFANode{},
		edges: make([]*NFAEdge, len(nfa.edges)),
	}
	for key, val := range nfa.nodes {
		res.nodes[key] = val
	}
	copy(res.edges, nfa.edges)
	return res
}

// Creates a new NFA
func New() *NFA {
	return &NFA{
		nodes: map[NFANodeName]*NFANode{},
		edges: []*NFAEdge{},
	}
}

var _svgTempDir string
var _svgCounter int
var DEBUG_SHOW_STEPS = false

// Generates an SVG for ``nfa`` and saves it to a temp directory
func debugShowStep(nfa *NFA, description string) {
	if !DEBUG_SHOW_STEPS {
		return
	}

	if len(_svgTempDir) == 0 {
		tempDir, err := ioutil.TempDir("", "nfa-to-regex-svgs")
		if err != nil {
			panic(err)
		}
		_svgTempDir = tempDir
		fmt.Println("Saving debug steps to:", _svgTempDir)
	}

	dot := NFA2Dot(nfa)

	_svgCounter += 1
	fname := fmt.Sprintf("nfa-%02d-%s.svg", _svgCounter, description)
	dotProc := exec.Command("dot", "-Tsvg", "-o", path.Join(_svgTempDir, fname))
	dotProc.Stdin = strings.NewReader(dot)

	err := dotProc.Run()
	if err != nil {
		panic(err)
	}
}

// Converts a NFA to a regular expression using the state removal technique
func NFA2Regex(nfa *NFA) string {
	nfa = nfa.shallowCopy()

	debugShowStep(nfa, "start")

	// 1. Create single initial and terminal nodes with empty transitions to
	initialNode := nfa._getOrCreateNode("__initial__")
	terminalNode := nfa._getOrCreateNode("__terminal__")
	for _, node := range nfa.nodes {
		if node.isInitial {
			nfa.addEdge(initialNode.name, node.name, "")
			nfa.replaceNode(node.name, &NFANode{
				name:       node.name,
				isInitial:  false,
				isTerminal: false,
			})
		}
		if node.isTerminal {
			nfa.addEdge(node.name, terminalNode.name, "")
			nfa.replaceNode(node.name, &NFANode{
				name:       node.name,
				isInitial:  false,
				isTerminal: false,
			})
		}
	}
	initialNode.isInitial = true
	terminalNode.isTerminal = true

	debugShowStep(nfa, "create-initial-terminal")

	// 2. Iteritively remove nodes which aren't the initial or terminal node
	for len(nfa.nodes) > 2 {
		for _, node := range nfa.nodes {
			if node == initialNode || node == terminalNode {
				continue
			}

			// Collect any loops (ie, where the node references its self) so they
			// can be converted to kleen star in the middle of new edges
			kleenStarValues := []string{}
			inEdges := nfa.edgesIn(node.name)
			for _, inEdge := range inEdges {
				if inEdge.srcNode == inEdge.dstNode {
					kleenStarValues = append(kleenStarValues, inEdge.value)
				}
			}

			kleenStarMiddle := addKleenStar(orJoin(kleenStarValues), len(kleenStarValues) > 1)
			for _, inEdge := range inEdges {
				if inEdge.srcNode == inEdge.dstNode {
					continue
				}
				for _, outEdge := range nfa.edgesOut(node.name) {
					if outEdge.srcNode == outEdge.dstNode {
						continue
					}

					nfa.addEdge(
						inEdge.srcNode.name,
						outEdge.dstNode.name,
						inEdge.value+kleenStarMiddle+outEdge.value,
					)
				}
			}

			nfa.removeNode(node.name)
			debugShowStep(nfa, fmt.Sprintf("remove-node-%s", node.name))
		}
	}

	// 3. Produce the regular expression
	res := make([]string, 0, len(nfa.edges))
	for _, edge := range nfa.edges {
		res = append(res, edge.value)
	}
	return orJoin(res)
}

// Generates a graphviz dot file from a NFA
func NFA2Dot(nfa *NFA) string {
	res := make([]string, 0, len(nfa.edges)+5)

	res = append(res, "\trankdir = LR;")

	for _, edge := range nfa.edges {
		label := edge.value
		if len(label) == 0 {
			label = "''"
		}
		res = append(res, fmt.Sprintf(
			"\t%q -> %q [label=%q];",
			edge.srcNode.name,
			edge.dstNode.name,
			label,
		))
	}

	for _, node := range nfa.nodes {
		if node.isInitial {
			res = append(res, fmt.Sprintf("\t%q [shape=point];", node.name+"__initial"))
			res = append(res, fmt.Sprintf("\t%q -> %q;", node.name+"__initial", node.name))
		}
		if node.isTerminal {
			res = append(res, fmt.Sprintf("\t%q [peripheries=2];", node.name))
		}

	}

	return "digraph g {\n" + strings.Join(res, "\n") + "\n}\n"
}

// Adds a kleen star to ``s``:
//   addKleenStar("") -> ""
//   addKleenStar("a") -> "a*"
//   addKleenStar("abc") -> "(abc)*"
//   addKleenStar("(abc|123)", true) -> "(abc|123)*"
func addKleenStar(s string, noWrap ...bool) string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s + "*"
	default:
		if len(noWrap) > 0 && noWrap[0] {
			return s + "*"
		}
		return fmt.Sprintf("(%s)*", s)
	}
}

// Joins a series of strings together in an "or" statement, ignoring empty
// strings::
//   orJoin({"a"}) -> "a"
//   orJoin({"a", "b"}) -> "(a|b)"
//   orJoin({"", "a", "b"}) -> "(a|b)"
func orJoin(inputStrs []string) string {
	strs := make([]string, 0, len(inputStrs))
	for _, s := range inputStrs {
		if len(s) > 0 {
			strs = append(strs, s)
		}
	}

	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	default:
		return "(" + strings.Join(strs, "|") + ")"
	}
}

// Generates a NFA which will match binary strings which are multiples of ``n``
func MakeNFAMultiplesOfN(n int) *NFA {
	str := func(i int) string { return strconv.Itoa((i)) }
	nfa := New()
	for i := 0; i < n; i += 1 {
		nfa.addEdge(str(i), str((i*2)%n), "0")
		nfa.addEdge(str(i), str((i*2+1)%n), "1")
	}
	nfa.addEdge("start", "0", "0")
	nfa.addEdge("start", "1", "1")
	nfa.nodes["start"].isInitial = true
	nfa.nodes["0"].isTerminal = true
	return nfa
}

// Generates an example NFA
func MakeNFASimple() *NFA {
	nfa := New()
	nfa.addEdge("1", "1", "a")
	nfa.addEdge("1", "2", "b")
	nfa.addEdge("2", "2", "c")
	nfa.addEdge("2", "3", "d")
	nfa.addEdge("3", "3", "e")
	nfa.addEdge("3", "1", "x")
	nfa.nodes["1"].isInitial = true
	nfa.nodes["3"].isTerminal = true
	return nfa
}

// Generates an NFA with multiple initial and terminal nodes
func MakeNFAManyMany() *NFA {
	nfa := New()
	nfa.addEdge("1", "2", "a")
	nfa.addEdge("2", "3", "b")
	nfa.addEdge("2", "2", "l")
	nfa.addEdge("4", "2", "x")
	nfa.addEdge("2", "5", "y")
	nfa.nodes["1"].isInitial = true
	nfa.nodes["4"].isInitial = true
	nfa.nodes["3"].isTerminal = true
	nfa.nodes["5"].isTerminal = true
	return nfa
}

func main() {
	DEBUG_SHOW_STEPS = true
	nfa := MakeNFAManyMany()
	//nfa := MakeNFAMultiplesOfN(3)
	regex := NFA2Regex(nfa)
	fmt.Println("Graph:")
	fmt.Println("https://dreampuf.github.io/GraphvizOnline/#" + url.PathEscape(NFA2Dot(nfa)))
	fmt.Println()
	fmt.Println("Regex:")
	fmt.Println(regex)
}
