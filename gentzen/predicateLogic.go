package gentzen

import (
	"errors"
	"log"
	"strings"
)

//check if s1 is instance of s2; if yes, variable and term return which
//variable was replaced by which term
func isInstanceOf(s1, s2 string) (val bool, variable, term string) {

	var v, r string

	val = false

	n1 := Parse(s1)
	n2 := Parse(s2)

	if !n2.MainConnective().isQuantifier() {
		return
	}

	v = n2.BoundVariable()

	have := getSubnodes(n1)

	want := getSubnodes(n2.Child1Must())

	if len(have) != len(want) {
		log.Println(have, " and ", want, "not same")
		return
	}

	for i := range have {

		if have[i].MainConnective() != want[i].MainConnective() {
			log.Println(have, " and ", want, "not same2")

			return
		}

		if have[i].IsAtomic() {
			if len(have[i].Terms()) != len(want[i].Terms()) {
				log.Println(have, " and ", want, "not same3")
				return
			}
			if have[i].Predicate() != want[i].Predicate() {
				log.Println(have, " and ", want, "not same4")
				return
			}
			j := findPos(v, want[i].Terms())
			if j == -1 {
				continue
			}
			r = have[i].Terms()[j]
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

func replaceTerms(n *Node, old, subst string) *Node {

	s := n.String()

	s = strings.ReplaceAll(s, old, subst)

	return Parse(s)
}

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

func (n *Node) renewRaw() {

	n.raw = n.predicateLetter

	for _, t := range n.term {
		n.raw = n.raw + t
	}
}

func (n *Node) hasTerm(t string) bool {

	ns := getSubnodes(n)

	for _, i := range ns {

		if !i.IsAtomic() {
			continue
		}

		if slicesContains(i.Terms(), t) {
			return true
		}
	}
	return false
}

func (n *Node) hasIllegalBoundVariables() (err error) {

	ns := getSubnodes(n)

	for _, e := range ns {
		if e.IsQuantifier() {
			v := e.BoundVariable()

			for _, f := range getSubnodes(e)[1:] {
				if f.IsQuantifier() {
					if f.BoundVariable() == v {
						err = errors.New("nested quantifiers with same variable name")
						return
					}
				}
			}
		}
	}
	return
}
