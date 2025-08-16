package gentzen

import (
	"errors"
	"log"
	"strings"
)

// check if s1 is instance of s2; if yes, variable and term return which
// variable was replaced by which term
func isInstanceOf(s1, s2 string) (val bool, variable, term string) {

	var v, r string

	val = false

	n1 := Parse(s1, !allowGreekUpper)
	n2 := Parse(s2, !allowGreekUpper)

	if !n2.Val.connective.isQuantifier() {
		return
	}

	v = n2.Val.variable

	have := n1.Linearize()

	want := n2.Child(0).Linearize()

	if len(have) != len(want) {
		log.Println(have, " and ", want, "not same")
		return
	}

	for i := range have {

		if have[i].Val.connective != want[i].Val.connective {
			log.Println(have, " and ", want, "not same2")

			return
		}

		if have[i].Check(isAtomic) {
			if len(have[i].Val.term) != len(want[i].Val.term) {
				log.Println(have, " and ", want, "not same3")
				return
			}
			if have[i].Val.predicateLetter != want[i].Val.predicateLetter {
				log.Println(have, " and ", want, "not same4")
				return
			}
			j := findPos(v, want[i].Val.term)
			if j == -1 {
				continue
			}
			r = have[i].Val.term[j]
			break

		}
	}
	if r == "" {
		return
	}
	n3 := replaceTerms(want[0], v, r)

	return n1.String() == n3.String(), v, r
}

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

		if slicesContains(i.Val.term, t) {
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
