package gentzen

import (
	"strings"
)

type tablNode struct {
	formula  *Node
	child    []*tablNode
	ancestor []*tablNode
	handled  bool
	closed   bool
}

func PrintSemanticTableaux(n *Node) string {

	t := semanticTableau(n)

	t.pruneTableau()
	var lt func(n *tablNode) string

	lt = func(n *tablNode) (r string) {

		if len(n.ancestor) > 0 {
			if n.ancestor[0].nodeIsClosed() {
				return
			}
		}
		if n.nodeIsClosed() {

			r = `[ \p{\nullset} ` + "\n"

		} else {

			r = `[ \p{` + Parse(n.formula.String()).StringLatex() + `} ` + "\n"
			//	r = `[ \p{` + n.formula.String() + `} ` + "\n"
		}
		for k := 0; k < len(n.child); k++ {

			r = r + lt(n.child[k])
		}

		r = r + ` ] ` + "\n" // !{\qbalance} `

		return r
	}
	templ := `\begin{forest}
%generated by gentzen
DATA\end{forest}

`
	return strings.ReplaceAll(templ, `DATA`, lt(t))

}

func (n *tablNode) pruneTableau() {
	t := flattenTableaux(n)

	for _, e := range t {
		if e.isClosed() {
			e.child = nil
			x := e.mkchild()
			x.setClosed()
		}
	}

}

func (n *Node) isInconsistent() bool {

	t := semanticTableau(n)

	t.pruneTableau()

	f := flattenTableaux(t)

	for _, e := range f {
		if len(e.child) == 0 {
			if !e.closed {
				return false
			}
		}
	}
	return true
}

func (n *tablNode) nodeIsClosed() bool {
	if n.closed {
		return true
	}
	return false
}

func (n *tablNode) isRoot() bool {
	return len(n.ancestor) == 0
}

func (n *tablNode) isBasic() bool {
	if n.formula.IsAtomic() {
		return true
	}

	if !n.formula.IsNegation() {
		return false
	}

	return n.formula.Child1Must().IsAtomic()
}

func (n *tablNode) setHandled() {
	n.handled = true
}

func (n *tablNode) setUnhandled() {
	n.handled = false
}

func (n *tablNode) addChild(c *tablNode) {
	n.child = append(n.child, c)
}

func (n *tablNode) addParent(p *tablNode) {
	n.ancestor = append(n.ancestor, p)
}

func (n *tablNode) setFormula(f *Node) {
	n.formula = f
}

func semanticTableau(n *Node) *tablNode {
	r := new(tablNode)
	r.formula = n
	r.setUnhandled()
	growSemanticTab(r)
	return r
}

func findFirstUnhandledAncestor(n *tablNode) *tablNode {

	if n.isRoot() {
		return n
	}

	if n.ancestor[0].handled {
		return n
	}

	if !n.ancestor[0].handled {
		return findFirstUnhandledAncestor(n.ancestor[0])
	}

	return n
}

func (n *tablNode) lastChild() *tablNode {
	if len(n.child) == 0 {
		return n
	}
	return n.child[len(n.child)-1]
}

func branchEnd(n *tablNode) bool {

	if !n.isBasic() {
		return false
	}

	return findFirstUnhandledAncestor(n).isRoot()
}

func flattenTableaux(n *tablNode) []*tablNode {

	var gs func(n *tablNode, list []*tablNode) []*tablNode

	gs = func(n *tablNode, list []*tablNode) []*tablNode {

		list = append(list, n)

		/**	if n.isBasic() {
					return list
				}
		**/
		for i := 0; i < len(n.child); i++ {
			list = gs(n.child[i], list)
		}

		return list
	}

	var list []*tablNode

	return gs(n, list)
}

func (n *tablNode) parent() *tablNode {

	if len(n.ancestor) == 0 {
		return n
	}

	return n.ancestor[len(n.ancestor)-1]
}

func (n *tablNode) isAncestor(c *tablNode) bool {

	if n == c {
		return true
	}

	for _, a := range c.ancestor {
		if a == n {
			return true
		}
	}
	return false
}

func (n *tablNode) mkchild() *tablNode {

	c := new(tablNode)

	c.ancestor = append(c.ancestor, n.ancestor...)
	c.ancestor = append(c.ancestor, n)

	c.handled = false
	c.closed = false

	n.child = append(n.child, c)

	return c
}

func (n *tablNode) growthpoint() bool {

	if len(n.child) != 0 {
		return false
	}

	if n.closed {
		return true
	}

	return true

}

func (n *tablNode) allHandled() bool {

	for _, e := range flattenTableaux(n) {
		if !e.handled {
			return false
		}
	}
	return true
}

func (n *tablNode) setClosed() {
	n.closed = true
}

func (n *tablNode) isClosed() bool {

	for _, e := range n.ancestor {

		if n.formula.String() == lneg+e.formula.String() {
			return true
		}

		if lneg+n.formula.String() == e.formula.String() {
			return true
		}
	}

	return false
}

func growSemanticTab(n *tablNode) {

	for !n.allHandled() {
		tflat := flattenTableaux(n)
		for _, e := range tflat {

			if !e.handled {

				if e.isBasic() {
					e.setHandled()
					break
				}
				f := e.formula

				eflat := flattenTableaux(e)

				if f.IsNegation() {

					if f.Child1Must().IsNegation() {
						for _, e2 := range eflat {
							if !e2.growthpoint() {
								continue
							}
							c1 := e2.mkchild()
							c1.setFormula(f.Child1Must().Child1Must())
						}
						e.setHandled()
						break
					}

					if f.Child1Must().IsConjunction() {
						for _, e2 := range eflat {
							if !e2.growthpoint() {
								continue
							}
							c1 := e2.mkchild()

							n1 := f.Child1Must().Child1Must()
							n2 := f.Child1Must().Child2Must()
							n3 := Disjoin(Negate(n1), Negate(n2))
							c1.setFormula(n3)
						}
						e.setHandled()
						break
					}
					if f.Child1Must().IsDisjunction() {
						for _, e2 := range eflat {
							if !e2.growthpoint() {
								continue
							}
							c1 := e2.mkchild()
							n1 := f.Child1Must().Child1Must()
							n2 := f.Child1Must().Child2Must()
							n3 := Conjoin(Negate(n1), Negate(n2))
							c1.setFormula(n3)
						}
						e.setHandled()
						break
					}
					if f.Child1Must().IsConditional() {
						for _, e2 := range eflat {
							if !e2.growthpoint() {
								continue
							}
							c1 := e2.mkchild()
							n1 := f.Child1Must().Child1Must()
							n2 := f.Child1Must().Child2Must()
							n3 := Conjoin(n1, Negate(n2))
							c1.setFormula(n3)
						}
						e.setHandled()
						break
					}

				}

				if f.IsConditional() {
					for _, e2 := range eflat {
						if !e2.growthpoint() {
							continue
						}
						c1 := e2.mkchild()
						c2 := e2.mkchild()
						c1.setFormula(Negate(f.Child1Must()))
						c2.setFormula(f.Child2Must())
					}
					e.setHandled()
					break
				}

				if f.IsConjunction() {
					for _, e2 := range eflat {
						if !e2.growthpoint() {
							continue
						}
						c1 := e2.mkchild()
						c2 := c1.mkchild()
						c1.setFormula(f.Child1Must())
						c2.setFormula(f.Child2Must())

					}
					e.setHandled()
					break
				}
				if f.IsDisjunction() {
					for _, e2 := range eflat {
						if !e2.growthpoint() {
							continue
						}
						c1 := e2.mkchild()
						c2 := e2.mkchild()
						c1.setFormula(Parse(f.Child1Must().String()))
						c2.setFormula(Parse(f.Child2Must().String()))
					}
					e.setHandled()
					break
				}
			}
		}
	}
}

/**
start from top
find unhandled child
append new children
repeat
**/

/**

go back ancestors until find first unhandled
append appropriate sentences

**/
