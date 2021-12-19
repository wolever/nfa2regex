# `nfa2regex`: convert NFAs (and DFAs) to regular expressions

An implementation of the state removal technique for converting an NFA to a
regular expression.

No optimization is done, so results may be surprising.

For documentation and more examples, see: https://pkg.go.dev/github.com/wolever/nfa2regex

# Examples

## A Simple NFA

A simple example NFA:

```golang
nfa := nfa2regex.NewNFA()
nfa.AddEdge("1", "1", "a")
nfa.AddEdge("1", "2", "b")
nfa.AddEdge("2", "2", "c")
nfa.AddEdge("2", "3", "d")
nfa.AddEdge("3", "3", "e")
nfa.AddEdge("3", "1", "x")
nfa.Nodes["1"].IsInitial = true
nfa.Nodes["3"].IsTerminal = true
```

Which produces the state diagram:

![State diagram for a simple NFA](https://github.com/wolever/nfa-to-regex/raw/main/examples/simple-nfa.svg)

And the regular expression:

```
a*bc*d(e|xa*bc*d)*
```


## Binary Multiples of 3

An NFA which matches multiples of 3:

```golang
func MakeNFAMultiplesOfX(x int) *NFA {
    str := func(i int) string { return strconv.Itoa((i)) }
    nfa := NewNFA()
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
```

Which produces the state diagram:

![State diagram for NFA matching multiples of 3](https://github.com/wolever/nfa-to-regex/raw/main/examples/multiples-of-3-nfa.svg)

And the regular expression::

```
(0(0|11|10(1|00)*01)*|11(0|11|10(1|00)*01)*|10(1|00)*01(0|11|10(1|00)*01)*)
```


## State Removal Steps

Given an NFA with multiple initial and terminal nodes:

```golang
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
```

These are the steps used to convert it to a regular expression:

1. Initial NFA:

    ![Initial NFA](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-01-start.svg)

2. Creation of initial and terminal states:

    ![Creation of initial and terminal nodes](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-02-create-initial-terminal.svg)

3. Removal of node ``4``:

    ![Removal of node 4](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-03-remove-node-4.svg)

4. Removal of node ``5``:

    ![Removal of node 5](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-04-remove-node-5.svg)

5. Removal of node ``1``:

    ![Removal of node 1](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-05-remove-node-1.svg)

6. Removal of node ``2``:

    ![Removal of node 2](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-06-remove-node-2.svg)

7. And finally, the removal of node ``3``:

    ![Removal of node 3](https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-07-remove-node-3.svg)

Which yields the regular expression:

```
(xl*b|xl*y|al*b|al*y)
```