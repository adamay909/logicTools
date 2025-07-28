package gentzen

//go:generate stringer -type=tokenID
import (
	"strings"
)

type tokenID int

type token struct {
	tokenType tokenID
	str       string
	variable  string   //for PL
	predicate string   //for PL
	term      []string //for PL
	index     int      //index within string

}

type tokenStr []token

const (
	tConj tokenID = 1 << iota
	tDisj
	tNeg
	tCond
	tAtomicSentence
	tUni
	tEx
	tNec
	tPos
	tPredicate
	tTerm
	tOpenb
	tCloseb
	tIdent
	tComma
	tNull tokenID = 0
)

const (
	tConnective tokenID = tConj | tDisj | tNeg | tCond | tUni | tEx | tNec | tPos

	tQuantifier tokenID = tUni | tEx

	tModalOperator tokenID = tNec | tPos

	tUnary tokenID = tNeg | tUni | tEx | tNec | tPos

	tBinary tokenID = tConj | tCond | tDisj
)

func (t token) isQuantifier() bool {
	return t.tokenType == tUni || t.tokenType == tEx
}

func (t token) isModalOperator() bool {
	return t.tokenType == tNec || t.tokenType == tPos
}

func (t token) isPredicate() bool {
	return t.tokenType == tPredicate
}

func (t token) isTerm() bool {
	return t.tokenType == tTerm
}

func (t token) isConnective() bool {
	if oPL {
		return (t.tokenType & tConnective) != 0
	}
	return t.tokenType != tAtomicSentence
}

func (t token) isAtomicSentence() bool {
	return t.tokenType == tAtomicSentence
}

func (t token) String() string {
	if isGreekLowerCase(t.str) {
		return greekCharOf(t.str)
	}
	if t.isQuantifier() {
		return t.str + t.variable
	}
	return t.str
}

func (t token) StringF(m PrintMode) string {

	if m != O_ProofChecker {
		return t.String()
	}

	if isGreekLowerCase(t.str) {
		return greekCharOf(t.str)
	}
	if t.isQuantifier() {
		for _, e := range connectivesPL {
			if e[2] == t.str {
				return e[O_PlainText] + t.variable
			}
		}
		//		return t.str + t.variable
	}
	return t.str
}
func (t token) isConj() bool {
	return t.tokenType == tConj
}

func (t token) isDisj() bool {
	return t.tokenType == tDisj
}

func (t token) isNeg() bool {
	return t.tokenType == tNeg
}

func (t token) isCond() bool {
	return t.tokenType == tCond
}

func (t token) isKCA() bool {
	return (t.tokenType & (tConj | tDisj | tCond)) != 0
}

func (t token) isUni() bool {
	return t.tokenType == tUni
}

func (t token) isEx() bool {
	return t.tokenType == tEx
}

func (t token) isNec() bool {
	return t.tokenType == tNec
}

func (t token) isPos() bool {
	return t.tokenType == tPos
}

func (t token) isUnary() bool {

	return t.isNeg() || t.isUni() || t.isEx() || t.isNec() || t.isPos()
}

func (t token) isBinary() bool {
	return t.isConj() || t.isDisj() || t.isCond()
}

func (t token) isOpenb() bool {
	return t.tokenType == tOpenb
}

func (t token) isCloseb() bool {
	return t.tokenType == tCloseb
}

func (t token) boundVariable() string {
	if !t.isQuantifier() {
		return ""
	}
	return t.variable
}

func (t token) terms() []string {
	return t.term
}

func (t tokenID) logicConstant() LogicalConstant {
	switch t {
	case tNeg:
		return Neg
	case tConj:
		return Conj
	case tDisj:
		return Disj
	case tCond:
		return Cond
	case tUni:
		return Uni
	case tEx:
		return Ex
	case tNec:
		return Nec
	case tPos:
		return Pos
	default:
		return None
	}
}

func cleanString(s string) string {

	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")

	for _, g := range greekLCBindings {
		s = strings.ReplaceAll(s, g[0], g[2])
	}

	return s
}

func isLowerCase(s string) bool {
	if len(s) > 1 {
		return false
	}

	if s == "=" {
		return false
	}

	if s == strings.ToUpper(s) {
		return false
	}

	return true
}

func isGreekLowerCase(s string) bool {

	for _, e := range greekLCBindings {
		if s == e[1] {
			return true
		}
	}

	return false
}

func isGreekFormulaVar(s string) bool {

	for _, e := range greekLCBindings {
		if s == e[2] {
			return true
		}
	}

	return false
}

func (t *token) isIdentity() bool {

	if t.tokenType != tPredicate {
		return false
	}

	return t.str == `=`

}

func (t *token) isNegIdentity() bool {

	if t.tokenType != tPredicate {
		return false
	}

	return t.str == `/=`

}
