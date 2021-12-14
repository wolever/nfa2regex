``nfa2regex``: convert NFAs (and DFAs) to regular expressions
================================================================

An implementation of the state removal technique for converting an NFA to a
regular expression.

No optimization is done, so results may be surprising.

Examples
========

A Simple NFA
------------

A simple example NFA:

.. code:: golang

    nfa := New()
    nfa.addEdge("1", "1", "a")
    nfa.addEdge("1", "2", "b")
    nfa.addEdge("2", "2", "c")
    nfa.addEdge("2", "3", "d")
    nfa.addEdge("3", "3", "e")
    nfa.addEdge("3", "1", "x")
    nfa.nodes["1"].isInitial = true
    nfa.nodes["3"].isTerminal = true

Which produces the state diagram:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/simple-nfa.svg

And the regular expression::

    a*bc*d(e|xa*bc*d)*


Binary Multiples of 3
---------------------

An NFA which matches multiples of 3:

.. code:: golang

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

Which produces the state diagram:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/multiples-of-3-nfa.svg

And the regular expression::

    (0(0|11|10(1|00)*01)*|11(0|11|10(1|00)*01)*|10(1|00)*01(0|11|10(1|00)*01)*)


State Removal Steps
-------------------

Given an NFA with multiple initial and terminal nodes:

.. code:: golang

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

These are the steps used to convert it to a regular expression:

1. Initial NFA:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-01-start.svg

2. Creation of initial and terminal states:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-02-create-initial-terminal.svg

3. Removal of node ``4``:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-03-remove-node-4.svg

4. Removal of node ``5``:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-04-remove-node-5.svg

5. Removal of node ``1``:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-05-remove-node-1.svg

6. Removal of node ``2``:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-06-remove-node-2.svg

7. And finally, the removal of node ``3``:

.. image:: https://github.com/wolever/nfa-to-regex/raw/main/examples/state-removal-steps/nfa-07-remove-node-3.svg

Which yields the regular expression::

    (xl*b|xl*y|al*b|al*y)
