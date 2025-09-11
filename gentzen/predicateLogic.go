package gentzen

import (
	"errors"
	"slices"
	"strings"
)

// check if s1 is instance of s2; if yes, variable and term return which
// variable was replaced by which term

func findPos(v string, list []string) int {

	for i := 0; i < len(list); i++ {
		if list[i] == v {
			return i
		}
	}

	return -1
}

/*
// check if s1 is modal instance of s2 (i.e., s1 is s2 minus modal operator)
func isModalInstanceOf(s1, s2 string) bool {

		n1 := Parse(s1, !allowGreekUpper)
		n2 := Parse(s2, !allowGreekUpper)

		if !n2.Val.connective.isModalOperator() {
			return false
		}

		return n2.Child(0).String() == n1.String()
	}
*/
func replaceTerms(n *Node, old, subst string) *Node {

	s := n.String()

	s = strings.ReplaceAll(s, old, subst)

	return Parse(s, !allowGreekUpper)
}

/*
func (n *Node) replaceTerm(p int, v string) (old, subst string) {

		if p < 0 {

			return
		}

		if len(n.term) <= p {
			return
		}

		old = n.term[p]
		subst = v
		n.term[p] = v

		return
	}
*/
func renewRaw(n *Node) {

	n.Val.raw = n.Val.predicateLetter

	for _, t := range n.Val.term {
		n.Val.raw = n.Val.raw + t
	}
}

func hasTerm(n *Node, t string) bool {

	ns := n.Linearize()

	for _, i := range ns {

		if !i.Check(isAtomic) {
			continue
		}

		if slices.Contains(i.Val.term, t) {
			return true
		}
	}
	return false
}

func hasIllegalBoundVariables(n *Node) (err error) {

	ns := n.Linearize()

	for _, e := range ns {
		if e.Check(isQuantifier) {
			v := e.Val.variable

			for _, f := range e.Linearize()[1:] {
				if f.Check(isQuantifier) {
					if f.Val.variable == v {
						err = errors.New("nested quantifiers with same variable name")
						return
					}
				}
			}
		}
	}
	return
}
