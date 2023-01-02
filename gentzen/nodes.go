package gentzen

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	raw                string
	connective         logicalConstant
	subnode1, subnode2 *Node
	parent             *Node
	variable           string
	predicateLetter    string
	term               []string
	flags              []string
}

// Child1 returns first chile of n if it exists. Returns ok=false
// if n is atomic.
func (n *Node) Child1() (m *Node, ok bool) {

	if n.IsAtomic() {
		ok = false
		return
	}

	ok = true

	return n.subnode1, ok
}

// Child2 returns first chile of n if it exists. Returns ok=false
// if n is atomic or has no second chile.
func (n *Node) Child2() (m *Node, ok bool) {

	if n.IsAtomic() {
		ok = false
		return
	}

	if n.IsUnary() {
		ok = false
		return
	}

	ok = true

	return n.subnode2, ok
}

func (n *Node) Child1Must() (m *Node) {
	var k *Node
	if _, ok := n.Child1(); !ok {
		return k
	}

	return n.subnode1
}

func (n *Node) Child2Must() (m *Node) {
	var k *Node
	if _, ok := n.Child2(); !ok {
		return k
	}
	return n.subnode2
}

// Parent returns parent of n. If n has no parent, ok is false.
func (n *Node) Parent() (m *Node, ok bool) {

	if n.parent == nil {
		ok = false
		return
	}

	ok = true

	return n.parent, ok
}

func (n *Node) ParentMust() (m *Node) {
	var k *Node
	if _, ok := n.Parent(); !ok {
		return k
	}
	return n.parent
}

// SetFlag sets a flag for n. Use for storing extra information.
func (n *Node) SetFlag(f string) {
	n.flags = append(n.flags, f)
}

func (n *Node) RmFlag(f string) {
	var newflags []string

	for _, s := range n.flags {
		if s != f {
			newflags = append(newflags, s)
		}
	}

	n.flags = newflags
	return
}

func (n *Node) HasFlag(f string) bool {
	for _, s := range n.flags {
		if s == f {
			return true
		}
	}
	return false
}

func (n *Node) IsAtomic() bool {
	return n.connective == ""
}

func (n *Node) SetFormula(f string) {
	n.raw = f
}

func (n *Node) Formula() string {
	return printNodePolish(n)
}

func (n *Node) IsBinary() bool {
	switch n.MainConnective() {
	case conj:
		return true
	case cond:
		return true
	case disj:
		return true
	default:
		return false
	}
}

func (n *Node) IsUnary() bool {

	return !n.IsBinary() && !n.IsAtomic()

}

func (n *Node) IsQuantifier() bool {
	return n.MainConnective() == uni || n.MainConnective() == ex
}

func (n *Node) IsNegation() bool {
	return n.MainConnective() == neg
}

func (n *Node) IsConditional() bool {
	return n.MainConnective() == cond
}

func (n *Node) IsConjunction() bool {
	return n.MainConnective() == conj
}

func (n *Node) IsDisjunction() bool {
	return n.MainConnective() == disj
}

func (n *Node) IsModal() bool {
	return n.MainConnective() == nec || n.MainConnective() == pos
}

func (n *Node) BracketClass() int {

	c := 0

	if n.IsAtomic() {
		if oPL && n.predicateLetter == "=" && n.parent != nil {
			if n.parent.IsQuantifier() {
				c++
			}
		}
		if isGreekFormulaVar(n.Predicate()) && len(n.Terms()) > 0 && n.parent != nil {
			if n.parent.IsNegation() {
				c = 2
			} else {
				c++
			}
		}
		return c
	}

	if n.IsQuantifier() {
		c = n.subnode1.BracketClass()
		return c
	}

	if n.IsModal() {
		c = n.subnode1.BracketClass()
		return c
	}

	if n.IsNegation() {
		c = n.subnode1.BracketClass()
		if oPL && n.parent != nil {
			if n.parent.IsQuantifier() {
				c++
			}
			if n.subnode1.Predicate() == "=" {
				if n.parent.IsNegation() {
					c++
				}
			}
		}
		return c
	}

	c = n.subnode1.BracketClass() + 1
	if n.subnode2.BracketClass() > n.subnode1.BracketClass() {
		c = n.subnode2.BracketClass() + 1
	}
	return c
}

func (n *Node) Class() int {

	c := 1

	if n.IsAtomic() {
		return c
	}

	c = n.subnode1.Class() + 1
	if n.IsBinary() {
		if n.subnode2.Class() > n.subnode1.Class() {
			c = n.subnode2.Class() + 1
		}
	}
	return c
}

func (n *Node) MainConnective() logicalConstant {
	return n.connective
}

func (n *Node) SetConnective(s logicalConstant) {
	n.connective = s
	return
}

func (n *Node) Generation() int {

	g := 0
	for ; n.parent != nil; n = n.parent {
		g++
	}

	return g
}

func (n *Node) SetBoundVar(v string) {
	n.variable = v
	return
}

func (n *Node) BoundVariable() string {
	return n.variable
}

func (n *Node) AddTerm(t ...string) {

	n.term = append(n.term, t...)
	return
}

func (n *Node) Terms() []string {
	return n.term
}

func (n *Node) Predicate() string {
	return n.predicateLetter
}

func (n *Node) SetAtomic() {
	n.subnode1 = nil
	n.subnode2 = nil
	n.connective = ""
}

func (n *Node) SetChild1(c *Node) {
	n.subnode1 = c
	return
}

func (n *Node) SetChild2(c *Node) {
	n.subnode2 = c
	return
}

func (n *Node) IsIdentity() bool {

	if !n.IsAtomic() {
		return false
	}

	if n.predicateLetter != "=" {
		return false
	}

	return true

}

func (n *Node) HasFreeVars() bool {

	return len(n.FreeVars()) == 0
}

func (n *Node) FreeVars() []string {

	var fv []string

	nodes := getSubnodes(n)

	for _, e := range nodes {
		if !e.IsAtomic() {
			continue
		}

		for _, t := range e.Terms() {

			f := e

			for ; f.parent != nil; f = f.parent {
				if f.BoundVariable() == t {
					break
				}
			}
			if f.parent == nil && f.BoundVariable() != t {
				fv = append(fv, t)
			}
		}
	}

	return fv
}

// String implements Stringer interface for node.
// Return string is formatted in Polish notation.
func (n *Node) String() string {

	return printNodePolish(n)
}

// StringEnglish return English string for n.
func (n *Node) StringEnglish() string {

	return printNodeInfix(n, mEnglish)
	//	return fixBrackets(printNodeInfix(n, mLatex), mLatex)

	//		return printNodeLatex(n)
}

// StringLatex return Latex string for n.
func (n *Node) StringLatex() string {

	return printNodeInfix(n, mLatex)
	//	return fixBrackets(printNodeInfix(n, mLatex), mLatex)

	//		return printNodeLatex(n)
}

// StringLatex return Latex string for n.
func (n *Node) StringMathJax() string {

	return printNodeInfix(n, mPlainLatex)
	//return printNodeMathJax(n)
}

// StringPlain returns plain Unicode string
func (n *Node) StringPlain() string {

	return printNodeInfix(n, mPlainText)
	//return printNodePlain(n)
}

// ConnectiveCount returns the number of connectives in n.
func (n *Node) ConnectiveCount() int {
	var count int
	s := n.String()
	for _, c := range s {
		for _, k := range connectives {
			if string(c) == k[0] {
				count++
			}
		}
	}

	return count
}

// String returns string for c.
func (c logicalConstant) String() string {
	return c.Stringf(mPolish)
}

// Sring returns formatted string for c.
func (c logicalConstant) Stringf(m printMode) string {

	for _, e := range connectives {
		if string(c) == e[0] {
			return e[int(m)]
		}
	}
	return ""
}

func getSubnodes(n *Node) []*Node {

	var gs func(n *Node, list []*Node) []*Node

	gs = func(n *Node, list []*Node) []*Node {

		list = append(list, n)

		if n.IsAtomic() {
			return list
		}

		if n.subnode1 != nil {
			list = gs(n.subnode1, list)
		}

		if n.subnode2 != nil {
			list = gs(n.subnode2, list)
		}

		return list
	}

	var list []*Node

	return gs(n, list)
}

// order nodes by depth
func _reorderNodes(nodes []*Node) (out []*Node) {

	d := findMaxDepth(nodes)

	for i := 0; i <= d; i++ {

		for _, j := range nodes {
			if j.Generation() == i {
				out = append(out, j)
			}
		}
	}
	return out
}

// get subnodes orderd by class
func orderedNodes(n *Node) (out []*Node) {

	d := n.Class() - 1
	nodes := getSubnodes(n)

	for i := d; i >= 0; i-- {

		for _, j := range nodes {
			if j.Class()-1 == i {
				out = append(out, j)
			}
		}
	}
	return out
}

func findMaxDepth(nodes []*Node) int {

	var ds []int

	for _, n := range nodes {
		ds = append(ds, n.Generation())
	}

	sort.Ints(ds)

	return ds[len(ds)-1]
}

// check if s1 is instance of s0
func sameStructure(s0, s1 string) bool {

	sn := normalize(s0, s1)
	n0 := getSubnodes(Parse(sn[0]))
	n1 := getSubnodes(Parse(sn[1]))

	if len(n0) > len(n1) {
		return false
	}

	atomic := n0[0].AtomicSentences()

	for _, a := range atomic {
		for i := range n0 {
			if n0[i].Formula() == a {
				repl := n1[i].Formula()
				n1[i].SetAtomic()
				for j := range n0 {
					if n0[j].HasFlag("c") {
						continue
					}
					if n0[j].Formula() == a {
						n0[j].SetFormula(repl)
						n0[j].SetFlag("c")
					}
				}
				n1 = getSubnodes(n1[0])
			}
		}
	}
	return n0[0].Formula() == n1[0].Formula()

}

func normalize(s ...string) []string {

	var out []string

	var allAtomic []string

	// set up function for returning series of sentence/predicate letters
	var nextatomic func() string

	for _, e := range s {
		allAtomic = append(allAtomic, Parse(e).AtomicSentences()...)
	}
	availLetters := []string{"P", "Q", "R", "S", "T", "F", "G", "H"}
	var normal []string

	for _, l := range availLetters {
		for n := 0; n < 10; n++ {
			normal = append(normal, l+"_"+strconv.Itoa(n))
		}
	}

	count := -1
	nextatomic = func() (ret string) {

		count++
		if count == len(normal) {
			fmt.Println("Too many atomic sentences/predicates")
			return "K"
		}
		ret = normal[count]

		if slicesContains(allAtomic, ret) {
			return nextatomic()
		}
		return ret
	}
	// done setting things up

	for _, e := range s {

		atomic := Parse(e).AtomicSentences()

		for _, a := range atomic {
			if !oPL {
				e = Parse(e).ReplaceAtomic(a, nextatomic()).Formula()
			} else {
				terms := strings.TrimPrefix(a, Parse(a).predicateLetter)
				e = Parse(e).ReplaceAtomic(a, nextatomic()+terms).Formula()
			}

		}
		out = append(out, e)
	}
	return out

}

func (n *Node) ReplaceAtomic(old, repl string) *Node {

	n1 := getSubnodes(n)

	for i := range n1 {
		if !n1[i].IsAtomic() {
			continue
		}
		if n1[i].Formula() == old {
			n1[i].SetFormula(repl)
		}
	}
	return n1[0]

}

// check if s2 is instance of s1
func _sameStructure(s1, s2 string) bool {

	ns1 := orderedNodes(Parse(s1))
	ns2 := orderedNodes(Parse(s2))

	//	fmt.Println("target: ", ns1[0].Formula(), "check: ", ns2[0].Formula())

	if len(ns2) < len(ns1) {
		return false
	}

	for i := range ns1 {
		//		fmt.Println("compare ", ns2[i].Formula(), " against ", ns1[i].Formula())
		if ns1[i].IsAtomic() {
			old := ns1[i].Formula()
			repl := ns2[i].Formula()
			for j := range ns1 {
				if ns1[j].Formula() == old {
					if ns2[j].Formula() != repl {
						return false
					}
				}
			}
		}
		continue

		if ns1[i].MainConnective() != ns2[i].MainConnective() {
			//			fmt.Println("mc mis 1: ", ns1[i].MainConnective(), "2 :", ns2[i].MainConnective())
			//			fmt.Println("Fail  ", ns1[i].Formula(), "  ", ns2[i].Formula())
			return false
		}
	}
	return true

}

func (n *Node) AtomicSentences() []string {

	var as []string

	ns := getSubnodes(n)

	for _, e := range ns {
		if !e.IsAtomic() {
			continue
		}

		if slicesContains(as, e.String()) {
			continue
		}

		as = append(as, e.String())
	}

	return as
}

func (n *Node) AtomicCount() int {

	return len(n.AtomicSentences())

}

func (n *Node) IsPureSL() bool {

	ns := getSubnodes(n)

	for _, n := range ns {

		if n.IsQuantifier() {
			return false
		}
		if n.IsModal() {
			return false
		}

	}

	return true
}
