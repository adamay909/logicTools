package gentzen

import (
	"math/rand"
)

func genRand(m, d int, fixed bool) string {

	var (
		atomicE = []string{
			"P",
			"Q",
			"R",
			"S",
			"T",
			"W",
			"Y",
			"Z",
			"M",
			"L",
		}
	)

	if m > len(atomicE) {
		m = len(atomicE)
	}

	s := Parse(chooseAtomic(m, atomicE))

	for s.ConnectiveCount() < d {

		nodes := getSubnodes(s)
		changed := false
		for _, e := range nodes {
			if !e.IsAtomic() {
				continue
			}
			if rand.Intn(4) == 0 {
				continue
			}
			changed = true
			e.connective = chooseConnective(connectivesSL)
			c1 := Parse(chooseAtomic(m, atomicE))
			e.SetChild1(c1)

			if e.IsBinary() {
				c2 := Parse(chooseAtomic(m, atomicE))
				e.SetChild2(c2)
			}
			e.raw = e.String()
			if s.ConnectiveCount() >= d {
				break
			}
		}
		if !changed {
			break
		}
	}

	return s.String()
}

func chooseAtomic(m int, atomicE []string) string {

	if m > len(atomicE) {
		m = len(atomicE)
	}

	return atomicE[rand.Intn(m)]

}

func chooseConnective(candidates [][6]string) logicalConstant {

	return logicalConstant(candidates[rand.Intn(len(candidates))][0])

}

// m number of predicates, d number of connectives
func genRandPL(m, d int) string {

	var quantifiers [][6]string

	quantifiers = append(quantifiers, connectivesPL[:2]...) //need to exclude equality sign

	atomicE := genAtomicE(m)

	s := Parse(chooseAtomic(len(atomicE), atomicE))

	for s.ConnectiveCount() < d {

		changed := false

		for _, e := range getSubnodes(s) {
			if !e.IsAtomic() {
				continue
			}
			if rand.Intn(4) == 0 {
				continue
			}

			changed = true
			e.SetConnective(chooseConnective(connectivesSL))

			c1 := Parse(chooseAtomic(len(atomicE), atomicE))
			e.clear()
			e.SetChild1(c1)

			if e.IsBinary() {

				c2 := Parse(chooseAtomic(len(atomicE), atomicE))
				e.clear()
				e.SetChild2(c2)
			}
			e.raw = e.String()
		}
		if !changed {
			break
		}
	}
	for s.HasFreeVars() {

		es := getSubnodes(s)

		e := es[rand.Intn(len(es))]

		if e.HasFreeVarsOnBranch() {
			if rand.Intn(4) == 0 {
				e.closure(chooseConnective(quantifiers), chooseAtomic(100, e.FreeVarsOnBranch()))
				e.raw = e.String()
			}
		}
	}
	s.renameVars()
	return s.String()
}

func (n *Node) FreeVarsOnBranch() []string {

	fv := n.FreeVars()

	for _, e := range n.Ancestors() {
		if !e.IsQuantifier() {
			continue
		}

		if slicesContains(fv, e.BoundVariable()) {
			fv = slicesRemove(fv, e.BoundVariable())
		}
	}
	return fv
}

func (n *Node) HasFreeVarsOnBranch() bool {
	return len(n.FreeVarsOnBranch()) > 0
}

func (n *Node) closure(quantifier logicalConstant, v string) {

	c1 := Parse(n.String())
	if _, ok := n.Child1(); ok {
		c1.SetChild1(n.Child1Must())
	}
	if n.IsBinary() {
		c1.SetChild2(n.Child2Must())
	}

	n.SetConnective(quantifier)

	n.SetChild1(c1)
	n.variable = v
	n.term = nil
	return
}

// generates the candidate atomic sentences for PL
func genAtomicE(m int) (r []string) {
	var (
		predicates = []string{
			"F",
			"G",
			"H",
			"D",
			"B",
			"P",
			"Q",
		}

		variables = []string{
			"a",
			"b",
			"c",
			"d",
			"e",
			"f",
		}
	)

	if m > len(predicates) {
		m = len(predicates)
	}

	nv := rand.Intn(len(variables)) + 1
	for _, p := range predicates[:m] {
		for _, v := range variables[:nv] {
			r = append(r, p+v)
		}
	}

	return r
}

func (s *Node) cleanEmptyQuantifiers() {

	for _, n := range getSubnodes(s) {

		if !n.hasEmptyQuantifiers() {
			continue
		}

		c := n.Child1Must()

		if n.parent == nil {
			n.connective = c.connective
			n.SetChild1(c.Child1Must())
			if c.IsBinary() {
				n.SetChild2(c.Child2Must())
			}
			continue
		}

		if n.parent.IsBinary() {
			if n.parent.Child2Must() == c {
				n.parent.SetChild2(c)
			} else {
				n.parent.SetChild1(c)

			}
		} else {
			n.parent.SetChild1(c)
		}
	}
	return
}
