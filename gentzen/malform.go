package gentzen

import (
	"math/rand"
	"strings"
)

// Malform returns a LaTeX string that is likely  to be a malformed formula.
func Malform(s string) string {

	var s2 string
	strategy := rand.Intn(4) + 1

	for len(s2) == 0 {

		switch strategy {

		case 1:
			s2 = addIllegal(s)

		case 2:
			s2 = removeToken(s)

		case 3:
			s2 = addToken(s)

		case 4:
			s2 = changeToken(s)
		}
	}

	return s2
}

func addIllegal(s string) string {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {
		panic(err)
		return s
	}

	tks := latexTokens(n)

	cand := `!*+-;:"'?/<>,.`

	r := string(cand[rand.Intn(len(cand))])

	pos := rand.Intn(len(tks) - 1)

	var ntks []string

	ntks = append(ntks, tks[:pos]...)

	ntks = append(ntks, r)

	ntks = append(ntks, tks[pos:]...)

	return strings.Join(ntks, "")

}

func removeToken(s string) string {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {
		return s
	}

	tks := latexTokens(n)

	i := rand.Intn(len(tks))
	if i == 0 {
		return strings.Join(tks[1:], "")
	}

	if i == len(tks)-1 {
		tks = tks[:i]
		return strings.Join(tks[:1], "")
	}

	ntks := tks[:i]
	ntks = append(ntks, tks[i+1:]...)
	return strings.Join(ntks, "")

}

func addToken(s string) string {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {
		return s
	}

	var cand = []string{
		`\land`,
		`\lor`,
		`\lnot`,
		`P`,
		`Q`,
		`R`,
	}

	if oCOND {
		cand = append(cand, `\limplies`)
	}

	c := n.classP()
	if c > len(brackets)-1 {
		c = len(brackets)
	}

	for i := range c {
		cand = append(cand, brackets[i][0])
		cand = append(cand, brackets[i][1])
	}

	tks := latexTokens(n)

	t := cand[rand.Intn(len(cand))]

	i := rand.Intn(len(tks) - 1)

	ntks := tks[:i]
	ntks = append(ntks, t)
	ntks = append(ntks, tks[i:]...)

	return strings.Join(ntks, "")
}

func changeToken(s string) string {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {
		return s
	}

	var cand = []string{
		`\land`,
		`\lor`,
		`\lnot`,
		`P`,
		`Q`,
		`R`,
	}

	if oCOND {
		cand = append(cand, `\limplies`)
	}

	c := n.classP()
	if c > len(brackets)-1 {
		c = len(brackets)
	}

	for i := range c {
		cand = append(cand, brackets[i][0])
		cand = append(cand, brackets[i][1])
	}

	tks := latexTokens(n)

	i := rand.Intn(len(tks))

	tks[i] = cand[rand.Intn(len(cand))]

	return strings.Join(tks, "")

}

func latexTokens(n *Node) []string {

	var resp []string

	w := new(strings.Builder)

	ingressFunc := func(e *Node) {
		w.Reset()
		latexIngressFunc(e, w)
		if w.String() != "" {
			resp = append(resp, w.String())
		}
	}

	pivotFunc := func(e *Node) {
		w.Reset()
		latexPivotFunc(e, w)
		if w.String() != "" {
			resp = append(resp, w.String())
		}
	}

	egressFunc := func(e *Node) {
		w.Reset()
		latexEgressFunc(e, w)
		if w.String() != "" {
			resp = append(resp, w.String())
		}
	}

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

	return resp

}

func infixTextTokens(n *Node, m PrintMode) []string {

	if m == O_Latex {
		return latexTokens(n)
	}

	var resp []string

	w := new(strings.Builder)

	ingressFunc := func(n *Node) {
		w.Reset()
		nomarkupIngressFunc(n, w, m)
		if w.String() != "" {
			resp = append(resp, w.String())
		}

	}

	pivotFunc := func(n *Node) {
		w.Reset()
		nomarkupPivotFunc(n, w, m)
		if w.String() != "" {
			resp = append(resp, w.String())
		}
	}

	egressFunc := func(n *Node) {
		w.Reset()
		nomarkupEgressFunc(n, w, m)
		if w.String() != "" {
			resp = append(resp, w.String())
		}
	}

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

	return resp

}
