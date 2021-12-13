``nfa-to-regex``: convert NFAs (and DFAs) to regular expressions
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

    *bc*d(e|xa*bc*d)*


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