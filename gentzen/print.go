package gentzen

import "strings"

type printMode int

func printNodePolish(n *Node) (s string) {

	switch {

	case n.IsUnary():

		s = n.connectiveString() + printNodePolish(n.subnode1)

	case n.IsBinary():

		s = n.connectiveString() + printNodePolish(n.subnode1) + printNodePolish(n.subnode2)

	default:
		s = n.raw

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
	return s
}

func printNodeInfix(n *Node, m printMode) (s string) {

	var br [][2]string

	br = brackets
	if m == mPlainText {
		br = append(brackets[:3], [2]string{`{`, `}`})
	}
	if m == mSimple {
		br = brackets[:2]
	}

	switch {

	case n.IsUnary():

		s = n.connectiveDisplay(m) + printNodeInfix(n.subnode1, m)
		//special case negated identity
		if n.MainConnective() == neg && n.subnode1.IsAtomic() && n.subnode1.predicateLetter == "=" {
			switch m {
			case mLatex:
				s = `\nident{` + n.subnode1.term[0] + `}{` + n.subnode1.term[1] + `}`
				return s
			case mPlainLatex:
				s = n.subnode1.term[0] + `\mathbin{\neq}` + n.subnode1.term[1] + ``
				return s
			case mPlainText:
				s = n.subnode1.term[0] + "\u2260" + n.subnode1.term[1]
				return s
			default:
			}
		}

	case n.IsBinary():

		s = printNodeInfix(n.subnode1, m) + n.connectiveDisplay(m) + printNodeInfix(n.subnode2, m)

	case n.predicateLetter == "=":
		switch m {
		case mLatex:
			s = `\ident{` + n.term[0] + `}{` + n.term[1] + `}`
		case mPlainLatex:
			s = n.term[0] + `\mathbin{=}` + n.term[1]
		case mPlainText:
			s = n.term[0] + "=" + n.term[1]
		default:
		}

	default:
		s = n.raw

	}

	if n.parent == nil {
		return s
	}

	if n.IsQuantifier() {
		return s
	}

	if n.IsNegation() {
		if !n.parent.IsQuantifier() {
			return s
		}
	}

	var ob1, ob2 string
	var blevel int

	blevel = n.BracketClass()

	if blevel+1 >= len(br) {
		ob1 = br[len(br)-1][0]
		ob2 = br[len(br)-1][1]
	} else {
		ob1 = br[blevel][0]
		ob2 = br[blevel][1]
	}

	return ob1 + s + ob2
}

func (c logicalConstant) StringLatex() string {
	switch c {
	case neg:
		return `\lnot `
	case disj:
		return ` \lor `
	case conj:
		return ` \land `
	case cond:
		return ` \limplies `
	case uni:
		return `\lforall `
	case ex:
		return `\lthereis `
	default:
		return ""
	}
}

func (c logicalConstant) StringMathJax() string {
	switch c {
	case neg:
		return `\neg `
	case disj:
		return ` \vee `
	case conj:
		return ` \wedge `
	case cond:
		return ` \supset `
	case uni:
		return `\forall `
	case ex:
		return `\exists `
	default:
		return ""
	}
}

func (c logicalConstant) StringPlain() string {
	switch c {
	case neg:
		return "\u00ac"
	case disj:
		return "\u2228"
	case conj:
		return "\u2227"
	case cond:
		return "\u2283"
	case uni:
		return "\u2200"
	case ex:
		return "\u2203"
	default:
		return ""
	}
}

func (n *Node) connectiveString() string {

	s := n.MainConnective().String()

	if n.MainConnective().isQuantifier() {
		s = s + n.BoundVariable()
	}

	return s
}

func (c logicalConstant) isQuantifier() bool {
	return c == uni || c == ex
}

func (seq sequent) StringLatex() string {
	var dat, suc string

	if len(seq.datum) != 0 {

		for _, d := range strings.Split(seq.datum, ",") {
			if d[:1] == `\` {
				dat = dat + d + `, `
			} else {
				dat = dat + printNodeInfix(Parse(d), mLatex) + `, `
			}
		}
		dat = strings.TrimRight(dat, ", ")
	}

	suc = printNodeInfix(Parse(seq.succedent), mLatex)

	return `\seq{` + dat + `}{` + suc + `}`
}

func (n *Node) connectiveDisplay(m printMode) string {
	var s string
	if m == mSimple {
		m = mPlainText
	}
	for _, c := range connectives {
		if string(n.MainConnective()) == c[0] {
			s = c[int(m)]
		}
	}
	if n.MainConnective().isQuantifier() {
		s = s + n.BoundVariable() + ` `
	}
	return s
}
