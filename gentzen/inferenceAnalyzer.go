package gentzen

import (
	"strconv"
	"strings"
)

type derivNode struct {
	line            argLine
	supportingLines []*derivNode
}

func getDerivation(lines []string, offset int) (al []argLine, ok bool) {

	ok = true
	for k, i := range lines {
		logger.SetPrefix("line " + strconv.Itoa(k+1) + ": ")
		pl, err := parseArgline(i)
		if err != nil {
			logger.Print(err.Error())
			ok = false
		}
		al = append(al, pl)
	}
	return al, ok
}

func lineRefsOK(derivation []argLine, offset int) bool {

	ok := true

	for n := range derivation {
		for _, l := range derivation[n].lines {
			if l < offset || l > n+offset {
				logger.Print("illegal reference to line ", l)
				ok = false
			}
		}
	}
	return ok
}

// getDerivTree retrieves the alleged derivation structure leading to
// sequent on line. Does not  check for actual validity. Illegal line references
// will cause a panic.
func getDerivTree(derivation []argLine, line int) (dt *derivNode) {

	dt = new(derivNode)
	dt.line = derivation[line]
	Debug(dt.line.lines, dt.line.inf)
	if len(dt.line.lines) > 0 {
		for _, n := range dt.line.lines {
			dt.supportingLines = append(dt.supportingLines, getDerivTree(derivation, n-1))
		}
	}

	return dt
}

// return derivation tree as a one dimensional slice
func flattenDerivTree(n *derivNode) []*derivNode {

	var list []*derivNode
	var gs func(n *derivNode, list []*derivNode) []*derivNode

	gs = func(n *derivNode, list []*derivNode) []*derivNode {

		list = append(list, n)

		for _, l := range n.supportingLines {
			list = gs(l, list)

		}
		return list
	}

	return gs(n, list)
}

func (d *derivNode) isTheorem() bool {

	if theorem(d) {
		return true
	}

	plines := flattenDerivTree(d)

	for _, l := range plines {
		if len(l.supportingLines) > 0 {
			continue
		}

		if l.line.inf != a {
			return false
		}
	}

	return true
}

func checkStep(d *derivNode) bool {

	return checkFunc(d.line.inf)(d)

}

func checkFunc(infRule string) func(*derivNode) bool {
	switch infRule {
	case a: //assumption
		return assumption

	case ki: //conjunction intro
		return conjI

	case ke: //conjunction elim
		return conjE

	case di: //disjunction intro
		return disjI

	case de: //disjunction elim
		return disjE

	case ci: //conditional intro
		return condI

	case ce: //conditional elim
		return condE

	case ni: //negation intro
		return negI

	case ne: //negation elim
		return negE

	case ue: //universal elim
		return uniE

	case ui: //universal intro
		return uniI

	case ei: //existential intro
		return exI

	case ee: //existential elimo
		return exE

	case ii: //identity introduction
		return idI

	case ie: //identity introduction
		return idE

	case li: //necessity introduction
		return necI

	case tli: //necessity introduction
		return necI_T

	case pli: //necessity introduction
		return necI_S4

	case mli: //necessity introduction
		return necI_S5

	case le: //necessity elim
		return necE

	case mi: //possibility intro
		return posI

	case me: //possibility elim
		return posE

	case mme: //possibility elim
		return posE_S5

	case sc: //possibility elim
		return scopeReplacement

	case sl: //possibility elim
		return sententialLogic

	case "premise": //premise
		return premise

	case "rewrite": //sequent rewrite
		return seqRewrite

	default: //check if we are dealing with a theorem
		return theoremDeriv
	}
}

func theoremDeriv(d *derivNode) bool {

	var tf []string

	thms := theoremsInUse()

	inf := d.line.inf

	Debug("Theorems:============")
	for _, e := range thms {
		Debug(Parse(e[2]).display())
	}
	Debug("=================")

	inf = strings.TrimSpace(inf)
	for _, thm := range thms {
		if inf == thm[1] {
			tf = append(tf, thm[2])
		}
	}

	if len(tf) > 0 {
		return (theorem(d))
	}

	//check name of derived rule
	t := strings.TrimSuffix(strings.TrimSpace(inf), "R")

	for i := range thms {
		if t == thms[i][1] {
			tf = append(tf, thms[i][2])
		}
	}
	if len(tf) > 0 {
		return (derivR(d))
	}

	logger.Print(inf, "unknown inference rule or theorem")

	return false
}

//check if it claims to be a derived rule
//if yes deriv(d)
//check if it claims to be a theorem

//if yes theorem(d)
