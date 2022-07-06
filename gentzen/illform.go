package gentzen

import (
	"math/rand"
)

//Illform returns a probably ill-formed string based on s.
//Make sure to seed properly (uses math/rand).
func Illform(s string) string {

	tks := latexTokens(Parse(s))

	strategy := rand.Intn(3) + 1

	var output string

	switch strategy {

	case 1:
		tks = removeToken(tks)
	case 2:
		tks = addToken(tks)
	case 3:
		tks = changeToken(tks)
	default:
	}

	for _, t := range tks {
		output = output + t.str
	}

	return output
}

func removeToken(tks tokenStr) tokenStr {

	i := rand.Intn(len(tks))
	if i == 0 {
		tks = tks[1:]
		return tks
	}

	if i == len(tks)-1 {
		tks = tks[:i]
		return tks
	}

	ntks := tks[:i]
	ntks = append(ntks, tks[i+1:]...)
	return ntks

}

var candidates = []string{
	"(",
	")",
	"[",
	"]",
	`\{`,
	`\}`,
	`\big(`,
	`\big)`,
	`\big[`,
	`\big]`,
	`\big\{`,
	`\big\}`,
	`\Big(`,
	`\Big)`,
	`\Big[`,
	`\Big]`,
	`\Big\{`,
	`\Big\}`,
	`\bigg(`,
	`\bigg)`,
	`\bigg[`,
	`\bigg]`,
	`\bigg\{`,
	`\bigg\}`,
	`\land `,
	`\lor `,
	`\lnot `,
	`P`,
	`Q`,
	`R`,
	//`\limplies `,
}

func addToken(tks tokenStr) tokenStr {

	i := rand.Intn(len(tks) + 1)

	var t token
	t.str = candidates[rand.Intn(len(candidates))]

	if i == len(tks) {
		tks = append(tks, t)
		return tks
	}

	if i == 0 {
		ntks := tokenStr{t}
		tks = append(ntks, tks...)
		return tks
	}

	ntks := tks[:i]
	ntks = append(ntks, t)
	ntks = append(ntks, tks[i:]...)

	return ntks
}

func changeToken(tks tokenStr) tokenStr {

	i := rand.Intn(len(tks))

	if tks[i].tokenType == tAtomicSentence {
		tks[i].str = candidates[rand.Intn(len(candidates)-3)]
	} else {
		tks[i].str = candidates[rand.Intn(len(candidates))]
	}
	return tks
}

func latexTokens(n *Node) (s tokenStr) {

	var t1, t2, t3 token
	switch {
	case n.IsUnary():
		t1.tokenType = tUnary
		t1.str = n.connectiveDisplay(mLatex)
		s = append(s, t1)
		s = append(s, latexTokens(n.subnode1)...)

	case n.IsBinary():
		t1.tokenType = tBinary
		t1.str = n.connectiveDisplay(mLatex)
		s = append(s, latexTokens(n.subnode1)...)
		s = append(s, t1)
		s = append(s, latexTokens(n.subnode2)...)

	default:
		t1.tokenType = tAtomicSentence
		t1.str = n.raw
		s = append(s, t1)

	}

	if n.parent == nil {
		return s
	}

	if n.MainConnective() == neg {
		return s
	}

	if n.IsAtomic() {
		return s
	}

	var ob1, ob2 string
	blevel := n.BracketClass()

	if blevel+2 > len(brackets) {
		ob1 = brackets[len(brackets)-1][0]
		ob2 = brackets[len(brackets)-1][1]
	} else {
		ob1 = brackets[blevel][0]
		ob2 = brackets[blevel][1]
	}
	t2.tokenType = tOpenb
	t2.str = ob1

	t3.tokenType = tCloseb
	t3.str = ob2

	tmp := tokenStr{t2}
	s = append(tmp, s...)
	s = append(s, t3)
	return s
}
