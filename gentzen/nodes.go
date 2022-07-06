package gentzen

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

//SetFlag sets a flag for n. Use for storing extra information.
func (n *Node) SetFlag(f string) {
	n.flags = append(n.flags, f)
}

//Child1 returns first chile of n if it exists. Returns ok=false
//if n is atomic.
func (n *Node) Child1() (m *Node, ok bool) {

	if n.IsAtomic() {
		ok = false
		return
	}

	ok = true

	return n.subnode1, ok
}

//Child2 returns first chile of n if it exists. Returns ok=false
//if n is atomic or has no second chile.
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

//Parent returns parent of n. If n has no parent, ok is false.
func (n *Node) Parent() (m *Node, ok bool) {

	if n.parent == nil {
		ok = false
		return
	}

	ok = true

	return n.parent, ok
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
		return false
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

func (n *Node) AddTerm(t string) {
	n.term = append(n.term, t)
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

func (n *Node) BracketClass() int {

	t := ""
	s := n.String()

	for i := 0; i < len(s); i++ {
		if s[i:i+1] == string(neg) {
			continue
		}
		if s[i:i+1] == luni {
			i++
			continue
		}
		if s[i:i+1] == lex {
			i++
			continue
		}
		t = t + s[i:i+1]
	}

	return Parse(t).Class() - 1
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

//String implements Stringer interface for node.
//Return string is formatted in Polish notation.
func (n *Node) String() string {

	return printNodePolish(n)
}

//StringLatex return Latex string for n.
func (n *Node) StringLatex() string {

	return printNodeInfix(n, mLatex)

	//	return printNodeLatex(n)
}

//StringLatex return Latex string for n.
func (n *Node) StringMathJax() string {

	return printNodeInfix(n, mPlainLatex)
	//return printNodeMathJax(n)
}

//StringPlain returns plain Unicode string
func (n *Node) StringPlain() string {

	return printNodeInfix(n, mPlainText)
	//return printNodePlain(n)
}

//ConnectiveCount returns the number of connectives in n.
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

//String returns string for c.
func (c logicalConstant) String() string {
	return c.Stringf(mPolish)
}

//Sring returns formatted string for c.
func (c logicalConstant) Stringf(m printMode) string {

	for _, e := range connectives {
		if string(c) == e[0] {
			return e[int(m)]
		}
	}
	return ""
}
