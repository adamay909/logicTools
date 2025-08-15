package gentzen

import (
	"errors"
	"strconv"
	"strings"

	"github.com/adamay909/logicTools/ju"
)

// Node holds information about a node in the syntax tree of a formula
type syntaxNode struct {
	raw             string
	connective      LogicalConstant
	variable        string
	predicateLetter string
	term            []string
	index           int //vaule of .index of token on which Node is based
	//tvassigned      []bool //whether truthvalue for given row has been assigned
	//tvalue          []bool //map interpretation row number to truth value
}

type Node = ju.Node[syntaxNode]

// IsAtomic returns true if n is an atomic formula.
func isAtomic(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == None
}

// IsPredicate returns true if n is a predicate.
func isPredicate(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.predicateLetter != ""
}

// IsConnective returns true if n has a main connective.
func isConnective(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective != None
}

// IsBinary returns true if n is a binary node.
func isBinary(n *Node) bool {
	switch n.Val.connective {
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
func isBasic(n *Node) bool {

	if isAtomic(n) {
		return true
	}
	if isNegation(n) {
		if n.Child(0).Check(isAtomic) {
			return true
		}
	}
	return false
}

// IsUnary returns true if n is a unary connective node.
func isUnary(n *Node) bool {

	return !isBinary(n) && !isAtomic(n)

}

// IsQuantifier returns true if n is a quantifier node.
func isQuantifier(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == Uni || n.Val.connective == Ex
}

// IsNegation returns true if n is a negation node.
func isNegation(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == Neg
}

// IsConditional returns true if n is a conditional node.
func isConditional(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == Cond

}

// IsConjunction returns true if n is a conjunction node.
func isConjunction(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == Conj
}

// IsDisjunction returhs true if n is a disjunction node.
func isDisjunction(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == Disj
}

// IsModal returns true if n is a modal operator node.
func isModal(n *Node) bool {

	if n == nil {
		return false
	}

	return n.Val.connective == Nec || n.Val.connective == Pos
}

// IsIdentity returns true if n is an identity node.
func isIdentity(n *Node) bool {

	if !isAtomic(n) {
		return false
	}

	if n.Val.predicateLetter != "=" {
		return false
	}

	return true

}

func isOpenFormula(n *Node) bool {

	return len(freeVars(n)) == 0
}

// FreeVars returns all the free variables in the formula
// represented by n.
func freeVars(n *Node) []string {

	var fv []string

	nodes := n.Linearize()

	for _, e := range nodes {
		if !e.Check(isAtomic) {
			continue
		}

		for _, t := range e.Val.term {

			f := e

			for ; f.Parent() != nil; f = f.Parent() {
				if f.Val.variable == t {
					break
				}
			}
			if f.Parent() == nil && f.Val.variable != t {
				fv = append(fv, t)
			}
		}
	}

	return fv
}

// String implements Stringer interface for node.
// Return string is formatted in Polish notation.
func stringFunc(n *Node) string {

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

	n.Walk(ingressFunc, pivotFunc, egressFunc)

	return w.String()
}

func init() {

	ju.SetStringFunc[syntaxNode](stringFunc)

}

// StringF returns the formula in the format specified by mode.
func StringF(n *Node, mode PrintMode) string {

	switch mode {

	case O_Latex, O_English, O_ProofChecker, O_PlainText, O_PlainASCII:

		return printNodeInfix(n, mode)

	default:

		return n.String()

	}
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

/*
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
*/
/*
//display is for displaying the text of a node that might have
//non-standard raw text

func (n *Node) display() string {

	return Parse(n.String(), !allowGreekUpper).StringF(O_PlainText)

}
*/
/*
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
*/
/*
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
*/
/*
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
*/
// IsPureSL returns true if the only logical constants are
// those of sentential logic (plus identity).
func isPureSL(n *Node) bool {

	return len(n.FindMatch(isQuantifier)) == 0
}

// Conjoin produces a Node that results by conjoining n1 and n1.
func Conjoin(n1, n2 *Node) *Node {

	e := new(Node)

	e.Val.connective = Conj
	e.AddChild(n1)
	e.AddChild(n2)

	return e
}

// Disjoin produces a node that results by disjoining n1 and n2.
func Disjoin(n1, n2 *Node) *Node {

	e := new(Node)

	e.Val.connective = Disj
	e.AddChild(n1)
	e.AddChild(n2)

	return e
}

// Negate produces a node that results by negating n.
func Negate(n *Node) *Node {

	e := new(Node)

	e.Val.connective = Neg
	e.AddChild(n)

	return e
}

// Conditionalize returns a conditional node that takes
// n1 as the antecedent and n2 as cosequent.
func Conditionalize(n1, n2 *Node) *Node {

	e := new(Node)

	e.Val.connective = Cond
	e.AddChild(n1)
	e.AddChild(n2)

	return e
}

func validate(n *Node) (err error) {

	var walk func(*Node) error

	walk = func(e *Node) error {

		if isUnary(e) && len(e.Children()) != 1 {
			err = errors.New("malformed: unary connective must have exactly one child node")

			err = errors.Join(errors.New(strconv.Itoa(e.Val.index)), err)
			return err
		}

		if isBinary(e) && len(e.Children()) != 2 {
			err = errors.New("malformed: binary connective must have exactly two child nodes")
			err = errors.Join(errors.New(strconv.Itoa(e.Val.index)), err)
			return err
		}

		if isAtomic(e) && len(e.Children()) != 0 {
			err = errors.New("malformed: non-connective cannot have a child")
			err = errors.Join(errors.New(strconv.Itoa(e.Val.index)), err)
			return err
		}

		if e.Val.predicateLetter != "" && len(e.Val.term) == 0 {
			var ch string
			ch, err = getFirstChar(e.Val.predicateLetter, !allowSubscr, !allowNumeral, !allowGreekUpper, !allowIdentity, allowSpecial)
			if err != nil {
				err = errors.Join(errors.New(strconv.Itoa(e.Val.index)), err)

				return err
			}
			if !isGreekLower(ch) {
				err = errors.New("malformed: predicate letter must be followed by at least one term (else use lower case Greek letter)")
				err = errors.Join(errors.New(strconv.Itoa(e.Val.index)), err)
				return err
			}
		}

		if isQuantifier(e) && e.Parent() != nil {

			for f := e.Parent(); f != nil; f = f.Parent() {
				if isQuantifier(f) && f.Val.variable == e.Val.variable {
					err = errors.New("illegal nested quantifier variables")
					err = errors.Join(errors.New(strconv.Itoa(e.Val.index)), err)
					return err
				}
			}
		}

		for _, c := range e.Children() {

			err = walk(c)

			if err != nil {
				break
			}
		}

		return err
	}

	return walk(n)
}

func binaryHeight(n *Node) int {

	h := 0
	lh := 0

	ingressFunc := func(n *Node) {
		if isBinary(n) {
			lh++
			if lh > h {
				h = lh
			}
		}
	}

	egressFunc := func(n *Node) {
		if isBinary(n) {
			lh--
		}
	}

	n.Walk(ingressFunc, func(n *Node) {}, egressFunc)

	return h
}

func setraw(n *Node) {

	if n.Val.predicateLetter == "" {
		return
	}

	n.Val.raw = n.Val.predicateLetter

	for _, t := range n.Val.term {

		n.Val.raw = n.Val.raw + t

	}

	return
}

func isSaturated(n *Node) bool {

	if isBinary(n) {
		return len(n.Children()) > 1
	}

	if isUnary(n) {
		return len(n.Children()) == 1
	}

	if n.Val.raw != "" {
		return true
	}

	return false
}

// returns the first open ancestor. Result takes the form of a slice
// so it can be used with ju's FindFunc method.
func openAncestor(n *Node) (r []*Node) {

	for e := n.Parent(); e != nil; e = e.Parent() {

		if !isSaturated(e) {
			r = append(r, e)
			break
		}

	}

	if len(r) == 0 {
		r = append(r, nil)
	}

	return r

}

// returns the first binary ancestor
func binaryAncestor(n *Node) (r []*Node) {

	for e := n.Parent(); e != nil; e = e.Parent() {
		if isBinary(e) {
			r = append(r, e)
			break
		}
	}

	return r
}

func binaryCount(n *Node) int {

	count := 0

	for _, e := range n.Linearize() {

		if isBinary(e) {
			count++
		}

	}

	return count
}

func isFunctionFormula(n *Node) bool {

	if n == nil {
		return false
	}

	if n.Val.predicateLetter == "" {
		return false
	}

	ch := strings.Split(n.Val.predicateLetter, "_")[0]

	if isGreekLower(ch) {
		return true
	}

	return false

}
