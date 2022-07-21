package gentzen

import (
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

	ns1 := getSubnodes(n1)

	ns2 := getSubnodes(n2)

	ns2 = ns2[1:]

	if len(ns1) != len(ns2) {
		log.Println(ns1, " and ", ns2, "not same")
		return
	}

	for i := range ns1 {

		if ns1[i].MainConnective() != ns2[i].MainConnective() {
			log.Println(ns1, " and ", ns2, "not same2")

			return
		}

		if ns1[i].IsAtomic() {
			if len(ns1[i].Terms()) != len(ns2[i].Terms()) {
				log.Println(ns1, " and ", ns2, "not same3")
				return
			}
			if ns1[i].Predicate() != ns2[i].Predicate() {
				log.Println(ns1, " and ", ns2, "not same4")
				return
			}
			j := findPos(v, ns2[i].Terms())
			if j == -1 {
				log.Println(ns1, " and ", ns2, "not same5")
				return
			}
			r = ns1[i].Terms()[j]
			break

		}
	}
	n3 := replaceTerms(ns2[0], v, r)

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
