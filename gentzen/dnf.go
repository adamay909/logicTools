package gentzen

import (
	"strings"
)

// replaceNK replaces NKpq with KNaNb
func replaceNK(n *Node) *Node {

	if n.MainConnective() != neg {
		return n
	}

	if n.Child1Must().MainConnective() != conj {
		return n
	}

	top := new(Node)

	top.SetConnective(disj)

	c1 := new(Node)
	c2 := new(Node)

	c1.SetConnective(neg)
	c2.SetConnective(neg)

	c1.subnode1 = n.Child1Must().Child1Must()
	c2.subnode1 = n.Child1Must().Child2Must()

	top.subnode1 = c1
	top.subnode2 = c2

	return top
}

// replaceNA replaces NApq with ANaNb
func replaceNA(n *Node) *Node {

	if n.MainConnective() != neg {
		return n
	}

	if n.Child1Must().MainConnective() != disj {
		return n
	}

	top := new(Node)

	top.SetConnective(conj)

	c1 := new(Node)
	c2 := new(Node)

	c1.SetConnective(neg)
	c2.SetConnective(neg)

	c1.subnode1 = n.Child1Must().Child1Must()
	c2.subnode1 = n.Child1Must().Child2Must()

	top.subnode1 = c1
	top.subnode2 = c2

	return top
}

func replaceC(n *Node) *Node {

	if n.MainConnective() != cond {
		return n
	}

	ant := n.Child1Must()
	cons := n.Child2Must()

	top := new(Node)

	top.SetConnective(disj)

	c1 := new(Node)

	c1.SetConnective(neg)

	c1.subnode1 = ant

	top.subnode1 = c1
	top.subnode2 = cons

	return top
}

func replaceNC(n *Node) *Node {

	if n.MainConnective() != neg {
		return n
	}

	if n.Child1Must().MainConnective() != cond {
		return n
	}

	ant := n.Child1Must().Child1Must()
	cons := n.Child1Must().Child2Must()

	top := new(Node)

	top.SetConnective(conj)

	c2 := new(Node)

	c2.SetConnective(neg)

	c2.subnode1 = cons

	top.subnode1 = ant
	top.subnode2 = c2

	return top
}

func replaceNN(n *Node) *Node {

	if n.MainConnective() != neg {
		return n
	}
	if n.Child1Must().MainConnective() != neg {
		return n
	}

	top := new(Node)

	top = n.Child1Must().Child1Must()

	return top
}

func dnfReplace(n *Node) *Node {

	if n.IsBasic() {
		return n
	}

	m := replaceNK(n)
	if !sameSentence(n, m) {
		return m
	}
	m = nil

	m = replaceNA(n)
	if !sameSentence(n, m) {
		return m
	}

	m = nil
	m = replaceC(n)
	if !sameSentence(n, m) {
		return m
	}

	m = nil
	m = replaceNC(n)
	if !sameSentence(n, m) {
		return m
	}

	m = nil
	m = replaceNN(n)
	if !sameSentence(n, m) {
		return m
	}

	return n
}

func reduceSentences(s []string) []string {

	var r []string
	for _, e := range s {
		r = append(r, dnfReplace(Parse(e)).String())
	}

	return r
}

func removeDuplicates(s []string) []string {

	var r []string

	for i := range s {
		found := false
		for j := range r {
			if s[i] == r[j] {
				found = true
				break
			}
		}
		if !found {
			r = append(r, s[i])
		}
	}
	return r
}

func splitConj(s []string) []string {

	var r []string

	for i := range s {

		if Parse(s[i]).MainConnective() != conj {
			r = append(r, s[i])
			continue
		}
		c1 := Parse(s[i]).Child1Must().String()
		c2 := Parse(s[i]).Child2Must().String()
		r = append(r, c1)
		r = append(r, c2)
	}

	return r
}

func splitDisjOnce(s []string) [][]string {

	var r1, r2 []string
	var r [][]string

	dp := func(s0 []string) int {
		k := -1
		for i := range s0 {
			if Parse(s0[i]).MainConnective() == disj {
				k = i
				break
			}
		}
		return k
	}

	n := dp(s)

	if n == -1 {
		r = append(r, s)
		return r
	}

	r1 = append(r1, s...)
	r2 = append(r2, s...)

	r1[n] = Parse(s[n]).Child1Must().String()
	r2[n] = Parse(s[n]).Child2Must().String()

	r = append(r, r1)
	r = append(r, r2)

	return r
}

func sameSentence(n, m *Node) bool {

	return n.String() == m.String()

}

func reduceDisj(s []string) []string {

	var r []string

	for _, e := range s {
		d := Parse(e)
		if d.MainConnective() != disj {
			r = append(r, e)
			continue
		}
		if d.Child1Must().String() != d.Child2Must().String() {
			r = append(r, e)
			continue
		}
		r = append(r, d.Child1Must().String())
	}

	return r
}

type reductionTree struct {
	set      []string
	children []*reductionTree
	parent   *reductionTree
}

func ReduceTree(s []string) *reductionTree {

	sameSet := func(s1, s2 []string) bool {
		if len(s1) != len(s2) {
			return false
		}

		for i := range s1 {
			if s1[i] != s2[i] {
				return false
			}
		}
		return true
	}

	n := new(reductionTree)

	n.set = append(n.set, s...)

	switch {
	case !sameSet(s, reduceSentences(s)):
		n.children = append(n.children, ReduceTree(reduceSentences(s)))

	case !sameSet(s, splitConj(s)):
		n.children = append(n.children, ReduceTree(splitConj(s)))

	case !sameSet(s, removeDuplicates(s)):
		n.children = append(n.children, ReduceTree(removeDuplicates(s)))

	case !sameSet(s, reduceDisj(s)):
		n.children = append(n.children, ReduceTree(reduceDisj(s)))

	case len(splitDisjOnce(s)) == 2:
		d := splitDisjOnce(s)
		n.children = append(n.children, ReduceTree(d[0]))
		n.children = append(n.children, ReduceTree(d[1]))

	default:
		return n

	}
	return n
}

func DNFtree(n *reductionTree) string {

	var lt func(n *reductionTree) string

	disp := func(k *reductionTree) string {
		var r string
		for i, e := range k.set {
			r = r + Parse(e).StringLatex()
			if i < len(k.set)-1 {
				r = r + ", "
			}

		}

		return r
	}

	lt = func(m *reductionTree) (r string) {

		r = `[ \p{` + disp(m) + ` } `
		r = r + "\n"

		if len(m.children) != 0 {

			for i := range m.children {

				r = r + lt(m.children[i])

			}

		}
		r = r + ` ] ` + "\n" // !{\qbalance} `

		return r
	}

	templ := `\begin{forest}{for tree={grow=south}}
%generated by gentzen
DATA\end{forest}

`
	return strings.ReplaceAll(templ, `DATA`, lt(n))

}

func DNF(n *reductionTree) *Node {

	t := flattenDNFtree(n)

	var dnf string

	var bc []string

	var bcs [][]string

	for _, s := range t {
		if len(s.children) != 0 {
			continue
		}
		bcs = append(bcs, s.set)
	}

	for _, cs := range bcs {
		var pc string

		for _, j := range cs {
			pc = lconj + pc + j
		}
		pc = pc[1:]
		bc = append(bc, pc)

	}

	for _, k := range bc {
		dnf = ldisj + dnf + k
	}

	dnf = dnf[1:]

	return Parse(dnf)

}

func flattenDNFtree(n *reductionTree) []*reductionTree {

	var gs func(n *reductionTree, list []*reductionTree) []*reductionTree

	gs = func(n *reductionTree, list []*reductionTree) []*reductionTree {

		list = append(list, n)

		if len(n.children) == 0 {
			return list
		}

		for _, c := range n.children {
			list = gs(c, list)
		}

		return list
	}

	var list []*reductionTree

	return gs(n, list)
}
