package gentzen

import (
	"errors"
	"strconv"
	"strings"
)

// Node holds information about a node in the syntax tree of a formula
type Node struct {
	raw                string
	connective         LogicalConstant
	subnode1, subnode2 *Node
	parent             *Node
	children           []*Node
	variable           string
	predicateLetter    string
	term               []string
	flags              []string
	index              int    //vaule of .index of token on which Node is based
	tvassigned         []bool //whether truthvalue for given row has been assigned
	tvalue             []bool //map interpretation row number to truth value
}

// Child1 returns first child of n if it exists. Returns ok=false
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

// Child1Must returns first child of n. Returns nil if there is no
// first child.
func (n *Node) Child1Must() (m *Node) {
	var k *Node
	if _, ok := n.Child1(); !ok {
		return k
	}

	return n.subnode1
}

// Child2Must returns second child of n. Returns nil if there is no
// second child.
func (n *Node) Child2Must() (m *Node) {
	var k *Node
	if _, ok := n.Child2(); !ok {
		return k
	}
	return n.subnode2
}

func (n *Node) mkchild() *Node {

	c := new(Node)
	c.parent = n
	c.SetAtomic()
	n.children = append(n.children, c)
	n.subnode1 = n.children[0]
	if len(n.children) > 1 {
		n.subnode2 = n.children[1]
	}
	return c

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

// ParentMust returns the parent of n. nil if n has no parent.
func (n *Node) ParentMust() (m *Node) {
	var k *Node
	if _, ok := n.Parent(); !ok {
		return k
	}
	return n.parent
}

// SetParent sets the parent of n to p.
func (n *Node) SetParent(p *Node) {

	n.parent = p

}

// SetFlag sets a flag for n. Use for storing extra information.
func (n *Node) SetFlag(f string) {
	n.flags = append(n.flags, f)
}

// RmFlag removes the flag f.
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

// HasFlag returns whether flag fhas been set.
func (n *Node) HasFlag(f string) bool {
	for _, s := range n.flags {
		if s == f {
			return true
		}
	}
	return false
}

// IsAtomic returns true if n is an atomic formula.
func (n *Node) IsAtomic() bool {

	if n == nil {
		return false
	}

	return n.connective == None
}

// IsPredicate returns true if n is a predicate.
func (n *Node) IsPredicate() bool {

	if n == nil {
		return false
	}

	return n.predicateLetter != ""
}

// IsConnective returns true if n has a main connective.
func (n *Node) IsConnective() bool {

	if n == nil {
		return false
	}

	return n.connective != None
}

// SetFormula sets the formula of n to f.
func (n *Node) SetFormula(f string) {
	n.raw = f
}

// Formula returns the formula of n.
func (n *Node) Formula() string {
	return n.String()
}

// IsBinary returns true if n is a binary node.
func (n *Node) IsBinary() bool {
	switch n.MainConnective() {
	case Conj:
		return true
	case Cond:
		return true
	case Disj:
		return true
	default:
		return false
	}
}

// IsBasic returns true if n is atomic or negation of an
// atomic formula.
func (n *Node) IsBasic() bool {

	if n.IsAtomic() {
		return true
	}
	if n.IsNegation() {
		if n.Child1Must().IsAtomic() {
			return true
		}
	}
	return false
}

// IsUnary returns true if n is a unary connective node.
func (n *Node) IsUnary() bool {

	return !n.IsBinary() && !n.IsAtomic()

}

// IsQuantifier returns true if n is a quantifier node.
func (n *Node) IsQuantifier() bool {

	if n == nil {
		return false
	}

	return n.MainConnective() == Uni || n.MainConnective() == Ex
}

// IsNegation returns true if n is a negation node.
func (n *Node) IsNegation() bool {

	if n == nil {
		return false
	}

	return n.MainConnective() == Neg
}

// IsConditional returns true if n is a conditional node.
func (n *Node) IsConditional() bool {

	if n == nil {
		return false
	}

	return n.MainConnective() == Cond

}

// IsConjunction returns true if n is a conjunction node.
func (n *Node) IsConjunction() bool {

	if n == nil {
		return false
	}

	return n.MainConnective() == Conj
}

// IsDisjunction returhs true if n is a disjunction node.
func (n *Node) IsDisjunction() bool {

	if n == nil {
		return false
	}

	return n.MainConnective() == Disj
}

// IsModal returns true if n is a modal operator node.
func (n *Node) IsModal() bool {

	if n == nil {
		return false
	}

	return n.MainConnective() == Nec || n.MainConnective() == Pos
}

func (n *Node) classP() int {

	class, sclass := 0, 0

	ln := linearize(n)

	for _, e := range ln {
		if !e.IsAtomic() {
			continue
		}
		sclass = 0
		for f := e.parent; f != nil; f = f.parent {
			sclass++
		}
		if sclass > class {
			class = sclass
		}
	}

	return class
}

// MainConnective returns the main connective of n.
func (n *Node) MainConnective() LogicalConstant {
	return n.connective
}

// SetConnective sets the main connective of n to s.
func (n *Node) SetConnective(s LogicalConstant) {
	n.connective = s
	return
}

// BoundVariable returns the variable bound by n.
// Return value is the empty string unless n is a quantifier node.
func (n *Node) BoundVariable() string {
	return n.variable
}

// Terms returns the terms of n. Returns an empty slice unless
// n is a predicate node.
func (n *Node) Terms() []string {
	return n.term
}

// Predicate returns the predicat letter.
func (n *Node) Predicate() string {
	return n.predicateLetter
}

// SetAtomic sets n to be an atomic formula.
func (n *Node) SetAtomic() {
	n.subnode1 = nil
	n.subnode2 = nil
	n.connective = None
}

// IsIdentity returns true if n is an identity node.
func (n *Node) IsIdentity() bool {

	if !n.IsAtomic() {
		return false
	}

	if n.predicateLetter != "=" {
		return false
	}

	return true

}

// Children returns the child nodes of n in a slice.
func (n *Node) Children() []*Node {

	var resp []*Node

	resp = append(resp, n.children...)

	return resp
}

// ClearChildren removes all child nodes of n.
func (n *Node) ClearChildren() {

	n.children = nil
}

// HasFreeVars returns true if there are free variables
// in the formula represented by n.
func (n *Node) HasFreeVars() bool {

	return len(n.FreeVars()) == 0
}

// FreeVars returns all the free variables in the formula
// represented by n.
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

	w := new(strings.Builder)

	ingressFunc := func(e *Node) {
		polishIngressFunc(e, w)
	}

	pivotFunc := func(e *Node) {
		polishPivotFunc(e, w)
	}

	egressFunc := func(e *Node) {
		polishEgressFunc(e, w)
	}

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

	return w.String()
}

// StringF returns the formula in the format specified by mode.
func (n *Node) StringF(mode PrintMode) string {

	switch mode {

	case O_Latex, O_English, O_ProofChecker, O_PlainText, O_PlainASCII:

		return printNodeInfix(n, mode)

	default:

		return n.String()

	}
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

// PolishString returns Polish string for c.
func (c LogicalConstant) PolishString() string {
	return c.Stringf(O_Polish)
}

// Stringf returns formatted string for c.
func (c LogicalConstant) Stringf(m PrintMode) string {

	for _, e := range connectives {
		if codeOf(c) == e[0] {
			return e[m]
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

// get subnodes orderd by class
func orderedNodes(n *Node) (out []*Node) {

	d := n.classP() - 1
	nodes := getSubnodes(n)

	for i := d; i >= 0; i-- {

		for _, j := range nodes {
			if j.classP()-1 == i {
				out = append(out, j)
			}
		}
	}
	return out
}

// check if s1 is instance of s0
func sameStructure(s0, s1 string) bool {

	sn := normalize(s0, s1)
	n0 := getSubnodes(Parse(sn[0], !allowGreekUpper))
	n1 := getSubnodes(Parse(sn[1], !allowGreekUpper))

	if len(n0) > len(n1) {
		return false
	}

	atomic := n0[0].AtomicSentences()

	Debug("<--sameStructure*************************")
	for k, a := range atomic {

		Debug("Round", k, ": compare ", n0[0].display(), " against ", n1[0].display())

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
			}
			n0 = getSubnodes(n0[0])
			n1 = getSubnodes(n1[0])
		}
	}
	Debug("Result: ", n0[0].display(), " against ", n1[0].display())
	Debug("--done structure check-->")

	return n0[0].Formula() == n1[0].Formula()

}

//display is for displaying the text of a node that might have
//non-standard raw text

func (n *Node) display() string {

	return Parse(n.String(), !allowGreekUpper).StringF(O_PlainText)

}

func normalize(s ...string) []string {

	var out []string

	var allAtomic []string

	// set up function for returning series of sentence/predicate letters
	var nextatomic func() string

	for _, e := range s {
		allAtomic = append(allAtomic, Parse(e, !allowGreekUpper).AtomicSentences()...)
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
			Debug("Too many atomic sentences/predicates")
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

		atomic := Parse(e, !allowGreekUpper).AtomicSentences()

		for _, a := range atomic {
			if !oPL {
				e = Parse(e, !allowGreekUpper).replaceAtomic(a, nextatomic()).Formula()
			} else {
				terms := strings.TrimPrefix(a, Parse(a, !allowGreekUpper).predicateLetter)
				e = Parse(e, !allowGreekUpper).replaceAtomic(a, nextatomic()+terms).Formula()
			}

		}
		out = append(out, e)
	}
	return out

}

func (n *Node) replaceAtomic(old, repl string) *Node {

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

// AtomicSentences returns a slice of the atomic sentences in the formula
// represented by n.
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

// AtomicCount returns the number of atomic sentences in n.
func (n *Node) AtomicCount() int {

	return len(n.AtomicSentences())

}

// IsPureSL returns true if the only logical constants are
// those of sentential logic (plus identity).
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

// Conjoin produces a Node that results by conjoining n1 and n1.
func Conjoin(n1, n2 *Node) *Node {

	s1 := n1.String()
	s2 := n2.String()
	s3 := lconj + s1 + s2

	return Parse(s3, !allowGreekUpper)
}

// Disjoin produces a node that results by disjoining n1 and n2.
func Disjoin(n1, n2 *Node) *Node {

	s1 := n1.String()
	s2 := n2.String()
	s3 := ldisj + s1 + s2

	return Parse(s3, !allowGreekUpper)
}

// Negate produces a node that results by negating n.
func Negate(n *Node) *Node {
	s1 := n.String()
	s2 := lneg + s1

	return Parse(s2, !allowGreekUpper)
}

// Conditionalize returns a conditional node that takes
// n1 as the antecedent and n2 as cosequent.
func Conditionalize(n1, n2 *Node) *Node {

	s1 := n1.String()
	s2 := n2.String()
	s3 := lcond + s1 + s2

	return Parse(s3, !allowGreekUpper)
}

func (n *Node) addFirstChild(n2 *Node) (err error) {

	if len(n.children) != 0 {
		err = errors.New("malformed: cannot add more than one child to node")
		return
	}

	n.subnode1 = n2
	n.children = append(n.children, n2)
	n2.parent = n

	return
}

func (n *Node) addSecondChild(n2 *Node) (err error) {

	if n == nil {
		err = errors.New("malformed: no appropriate parent node found")
		return
	}

	if len(n.children) > 1 {
		err = errors.New("malformed: cannot add more than two children to node")
		return
	}

	n.subnode2 = n2
	n.children = append(n.children, n2)
	n2.parent = n

	return
}

func (n *Node) rootNode() *Node {

	if n == nil {
		return nil
	}

	e := new(Node)

	for e = n; e.parent != nil; e = e.parent {
	}

	return e

}

func (n *Node) validate() (err error) {

	var walk func(*Node) error

	walk = func(e *Node) error {

		if e.IsUnary() && len(e.children) != 1 {
			err = errors.New("malformed: unary connective must have exactly one child node")

			err = errors.Join(errors.New(strconv.Itoa(e.index)), err)
			return err
		}

		if e.IsBinary() && len(e.children) != 2 {
			err = errors.New("malformed: binary connective must have exactly two child nodes")
			err = errors.Join(errors.New(strconv.Itoa(e.index)), err)
			return err
		}

		if e.IsAtomic() && len(e.children) != 0 {
			err = errors.New("malformed: non-connective cannot have a child")
			err = errors.Join(errors.New(strconv.Itoa(e.index)), err)
			return err
		}

		if e.IsPredicate() && len(e.term) == 0 {
			var ch string
			ch, err = getFirstChar(e.predicateLetter, !allowSubscr, !allowNumeral, !allowGreekUpper, !allowIdentity, allowSpecial)
			if err != nil {
				err = errors.Join(errors.New(strconv.Itoa(e.index)), err)
				return err
			}
			if !isGreekLower(ch) {
				err = errors.New("malformed: predicate letter must be followed by at least one term (else use lower case Greek letter)")
				err = errors.Join(errors.New(strconv.Itoa(e.index)), err)
				return err
			}
		}

		if e.IsQuantifier() && e.parent != nil {

			for f := e.parent; f != nil; f = f.parent {
				if f.IsQuantifier() && f.variable == e.variable {
					err = errors.New("illegal nested quantifier variables")
					err = errors.Join(errors.New(strconv.Itoa(e.index)), err)
					return err
				}
			}
		}

		for _, c := range e.children {

			err = walk(c)

			if err != nil {
				break
			}
		}

		return err
	}

	return walk(n)
}

func (n *Node) nestingDepth() int {

	d := 0

	var walk func(*Node)

	walk = func(e *Node) {

		if l := e.binaryGenerationNumber(); l > d {
			d = l
		}

		for _, c := range e.children {
			walk(c)
		}
	}

	m := Parse(n.String(), !allowGreekUpper)

	walk(m)

	if !prettifyBrackets && d > 1 {
		d = 1
	}

	return d
}

func (n *Node) binaryGenerationNumber() int {

	d := 0

	if n.parent == nil {
		return d
	}

	e := new(Node)

	for e = n.parent; e != nil; e = e.parent {
		if e.IsBinary() {
			d++
		}
	}

	return d
}

func (n *Node) setraw() {

	if n.predicateLetter == "" {
		return
	}

	n.raw = n.predicateLetter

	for _, t := range n.term {

		n.raw = n.raw + t

	}

	return
}

// AddChild adds n2 as a child of n. The order of children usually matters
// so be sure to add them in the right order.
func (n *Node) AddChild(n2 *Node) (err error) {

	if n == nil {
		return errors.New("no parent to add child")
	}

	if n.isSaturated() {
		return errors.New("cannot add child")
	}

	n.children = append(n.children, n2)

	n2.parent = n

	n.subnode1 = n.children[0]

	if len(n.children) > 1 {

		n.subnode2 = n.children[1]

	}

	return

}

func (n *Node) isSaturated() bool {

	if n.IsBinary() {
		return len(n.children) > 1
	}

	if n.IsUnary() {
		return len(n.children) == 1
	}

	if n.raw != "" {
		return true
	}

	return false
}

func (n *Node) removeSecondChild() {

	if !n.IsBinary() {
		return
	}

	if len(n.children) < 2 {
		return
	}

	n.subnode2 = nil

	n.children = nil

	n.children = append(n.children, n.subnode1)

}

func (n *Node) openAncestor() (e *Node) {

	for e = n.parent; e != nil; e = e.parent {

		if !e.isSaturated() {
			break
		}

	}

	return e

}

func (n *Node) binaryAncestor() (e *Node) {

	for e = n.parent; e != nil; e = e.parent {
		if e.IsBinary() {
			break
		}
	}

	return e
}

func (n *Node) binaryCount() int {

	count := 0

	for _, e := range linearize(n) {

		if e.IsBinary() {
			count++
		}

	}

	return count
}

func (n *Node) isFunctionFormula() bool {

	if n == nil {
		return false
	}

	if n.predicateLetter == "" {
		return false
	}

	ch := strings.Split(n.predicateLetter, "_")[0]

	if isGreekLower(ch) {
		return true
	}

	return false

}
