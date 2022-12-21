/*
Package gentzen provides some tools for checking proofs (both porpositional logic and first order predicate logic with identity),, printing truth tables, syntax trees (propositional logic) as well as a few other tools for creating exercises and the like.

The main entry points expect inputs as plain strings in the Polish notation. The default for the logical constants is non-standard:

- negation: -
- conjunction: ^
- disjunction: V
- conditional: >
- universal quantifier: U
- existential quantifier X

This frees up some letters of the alphabet for other uses .You can switch to a more standard Polish notation with

	SetStandardPolish(true)

which will switch the notation to:

- negation: N
- conjunction: K
- disjunction: A
- conditional: C
- universal quantifier: U
- existential quantifier X

There is a parser for an infix notation. That is designed to be used with the online proofchecker in github.com/adamay909/logicTools/proofChecker and requires a pseudo-tokenized input in the form of a slice of tokens.
*/
package gentzen

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	oPL   = false
	oTHM  = false
	oCOND = true
	oML   = true
	oDR   = true
)

func init() {

	rand.Seed(int64(time.Now().Nanosecond()))
	SetStandardPolish(true)
	//connectives = append(connectivesSL, connectivesPL...)
	logger = log.New(&checkLog, "", 0)

}

// ShowLog displays log. Currently, logging is only done by proof checker.
func ShowLog() string {

	return checkLog.String()

}

// ClearLog clears log.
func ClearLog() {
	checkLog.Reset()
	return
}

// WriteLog writes s to log.
func WriteLog(s string, p string) {
	logger.SetPrefix(p)
	logger.Print(s)
	return
}

// SetConditional determines whether we use the conditional.
func SetConditional(v bool) {
	oCOND = v
	if oCOND {
		connectivesSL = [][6]string{
			{string(neg), `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that "},
			{string(conj), `\land `, `\wedge `, "\u2227", "\u2227", " and "},
			{string(disj), `\lor `, `\vee `, "\u2228", "\u2228", " or "},
			{string(cond), `\limplies `, `\supset `, "\u2283", "\u2283", " if , then "},
		}
	} else {
		connectivesSL = [][6]string{
			{string(neg), `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that "},
			{string(conj), `\land `, `\wedge `, "\u2227", "\u2227", " and "},
			{string(disj), `\lor `, `\vee `, "\u2228", "\u2228", " or "},
		}
	}
	connectives = append(connectivesSL, connectivesPL...)
	return
}

// SetPL enables Predicate Logic Processing when v is true.
// By default, PL processing is disabled.
func SetPL(v bool) {

	oPL = v

}

// SetML specifies whether we allows modal logic
func SetML(v bool) {

	oML = v

}

// SetDR specifies whether we allow derived rules
func SetDR(v bool) {

	oDR = v

}

// SetAllowTheorems sets whether appeal to some standard theorems
// is allowed. Default is false.
func SetAllowTheorems(v bool) {
	oTHM = v
}

// SetStrict sets whether inferece rules should be checked strictly.
func SetStrict(v bool) {
	strictCheck = v
}

/*
CheckDeriv checks the derivation given by lines.
Each line represents a sequent followed by an annotation
consisting of line references and inference rule/theorem name.
The elements of each line must be semicolon separated. You may use
spaces and tabs for readability. The available rules are:

ne: negation elimination

ni: negation introduction

de: disjunction elimination

di: disjunction introduction

ke: conjunction elimination

ki: conjunction introduction

ce: conditional elimination

ci: conditional introduction

ue: universal quantifier elimination

ui: universal quantifier introduction

ei: existential quantifier introduction

ee: existential quantifier elimination

see the textbook on how the rules work.
*/
func CheckDeriv(lines []string, offset int) bool {

	return checkDerivation(lines, offset)

}

// GenTruthTable prints truth table for s. Obsolete. use
// MkTable and MkTextTable instead.
func GenTruthTable(s string) string {

	//return printTable(genTableSpec(getColumns(s)))
	return printTable(MkTable(s))
}

// GenEmptyTruthTable prints an empty truth table for s.
// Obsolete. Use functionality provided by TextTable type instead.
func GenEmptyTruthTable(s string) string {

	//	return printEmptyTable(genTableSpec(getColumns(s)))
	return printEmptyTable(MkTable(s))

}

// GenTruthTableValues generates the truth values for
// each cell of columns. Obsolete. Use functionality provided
// by Table type instead.
func GenTruthTableValues(s string) [][]bool {

	return MkTable(s).vals
	//genTable(genTableSpec(getColumns(s)))
}

// IsTautology returns whether s is a tautology.
func IsTautology(s string) bool {

	t := MkTable(s).vals
	//genTable(genTableSpec(getColumns(s)))

	lc := t[len(t)-1]

	for _, v := range lc {
		if !v {
			return false
		}
	}
	return true
}

// IsWff returns whether s is a wff.
func IsWff(s string) bool {

	_, err := ParseStrict(s)

	return err == nil

}

// PrintDeriv prints the derivation given by lines
// as latex formatted derivation.
func PrintDeriv(lines []string, offset int) (out string) {
	seq1, err1 := parseArgline(lines[0])
	seq2, err2 := parseArgline(lines[len(lines)-1])

	if err1 == nil && err2 == nil {
		if seq2.seq.d == "" {
			out = out + `%Prove \p{` + seq2.seq.StringLatex() + "}\n\n"
		} else {
			out = out + `%Derive from \p{` + seq1.seq.StringLatex() + `} to \p{` + seq2.seq.StringLatex() + "}\n\n"
		}
	}
	out = out + `\begin{argumentN}[` + strconv.Itoa(offset) + "]\n"

	out = out + `%generated by gentzen` + "\n\n"

	for i := range lines {

		out = out + printArgLine(lines[i], mLatex)
	}

	out = out + `\end{argumentN}` + "\n\n"

	return out
}

// PrintDeriv prints the derivation given by lines
// as latex formatted derivation.
func PrintDerivText(lines []string, offset int) (out string) {

	for i := range lines {

		out = out + strconv.Itoa(i+offset) + `. ` + printArgLine(lines[i], mPlainText) + "\n"
	}

	return out
}

// RandomSentence returns a randomly generated sentence of
// sentential logic with at most m atomic sentences (m is capped at 10).
// d specifies the maximum number of connectives. It uses the package math/rand
// and you need to seed the PRNG yourself.
func RandomSentence(m, d int, fixed bool) string {

	return genRand(m, d, fixed)

}

// Parse should only be used when s is known to be well-formed.
func Parse[S ~string](s S) *Node {

	n, err := ParseStrict(s)

	if err != nil {
		fmt.Println("malformed formula: ", s, " ", err.Error())
	}
	return n
}

// ParseStrict parses s.
func ParseStrict[str ~string](s str) (n *Node, err error) {
	tokens, err := tokenize(string(s))
	if err != nil {
		return n, err
	}
	n, err = parseTokens(tokens)
	if err != nil {
		return
	}
	err = n.hasIllegalBoundVariables()
	return
}
