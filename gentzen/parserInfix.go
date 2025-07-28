package gentzen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseInfix is like [ParseStrict] but parses a formula given in infix form. The input must be what you would get through n.StringF(O_PlainASCII).
func ParseInfix(s string, permitGreekUpper bool) (n *Node, err error) {

	tk, err := tokenizeInfix(s, permitGreekUpper)

	if err != nil {
		err = errors.New(showerrloc(err, s))
		return n, err
	}

	n, err = parseInfix(tk)

	if err != nil {
		err = errors.New(showerrloc(err, s))
		return n, err
	}

	return

}

// ParseLatex parses a formula given in LaTeX. The LaTeX must be code
// you would get through n.StringF(O_Latex).
func ParseLatex(s string, permitGreekUpper bool) (n *Node, err error) {

	s = latex2plaintxt(s, permitGreekUpper)

	return ParseInfix(s, permitGreekUpper)

}

func latex2plaintxt(s string, permitGreekUpper bool) string {

	for _, b := range brackets[1:] {

		s = strings.ReplaceAll(s, b[0], "(")
		s = strings.ReplaceAll(s, b[1], ")")

	}

	for _, g := range greekLCBindings {

		s = strings.ReplaceAll(s, g[1], g[0])

	}

	if permitGreekUpper {

		for _, g := range greekUCBindings {

			s = strings.ReplaceAll(s, g[1], g[0])

		}

	}

	s = strings.ReplaceAll(s, `\mc{\land }`, `^`)

	s = strings.ReplaceAll(s, `\mc{\lor }`, `v`)

	s = strings.ReplaceAll(s, `\mc{\lnot }`, `-`)

	s = strings.ReplaceAll(s, `\mc{\limplies }`, `>`)

	s = strings.ReplaceAll(s, `\land`, `^`)

	s = strings.ReplaceAll(s, `\lor`, `v`)

	s = strings.ReplaceAll(s, `\lnot`, `-`)

	s = strings.ReplaceAll(s, `\limplies`, `>`)

	if oPL {

		s = strings.ReplaceAll(s, `\lthereis`, `X`)

		s = strings.ReplaceAll(s, `\lforall`, `U`)

		s = strings.ReplaceAll(s, `\mc{\lthereis }`, `X`)

		s = strings.ReplaceAll(s, `\mc{\lforall }`, `U`)
	}

	return s

}

func tokenizeInfix(s string, permitGreekUpper bool) (tk tokenStr, err error) {

	s = strings.ReplaceAll(s, `\`, `/`)

	for c := 0; c < len(s); c++ {

		var t token

		t.index = c

		switch s[c : c+1] {

		case " ":
			continue //we allow space within string for readability

		case "(":
			t.tokenType = tOpenb
			t.str = "("

		case ")":
			t.tokenType = tCloseb
			t.str = ")"

		case "^":
			t.tokenType = tConj
			t.str = s[c : c+1]

		case "v":
			t.tokenType = tDisj
			t.str = s[c : c+1]

		case ">":
			t.tokenType = tCond
			t.str = s[c : c+1]

		case "-":
			t.tokenType = tNeg
			t.str = s[c : c+1]

		case ",":
			t.tokenType = tComma
			t.str = ","

		}

		if t.tokenType != tNull {
			tk = append(tk, t)
			continue
		}

		switch {
		case !oPL:
			t.tokenType = tAtomicSentence
			ch, err := getFirstChar(s[c:], allowSubscr, !allowNumeral, permitGreekUpper, !allowIdentity, !allowSpecial)
			if err != nil {
				return tk, err
			}
			t.str = ch
			c = c + len(ch) - 1

		case oPL:
			switch {

			case strings.HasPrefix(s[c:], luni):
				t.tokenType = tUni

			case strings.HasPrefix(s[c:], lex):
				t.tokenType = tEx

			default:
				ch, err := getFirstChar(s[c:], allowSubscr, !allowNumeral, permitGreekUpper, allowIdentity, !allowSpecial)

				if err != nil {
					err = errors.New(strconv.Itoa(c))
					return tk, err
				}

				t.str = ch

				chb := strings.Split(ch, "_")[0]

				switch {
				case isRomanLower(rune(chb[0])):
					t.tokenType = tTerm

				case isRomanUpper(rune(chb[0])):
					t.tokenType = tPredicate

				case isGreekLower(chb):
					t.tokenType = tPredicate

				case isGreekUpper(chb):
					t.tokenType = tAtomicSentence

				case chb == "=":
					t.tokenType = tIdent

				case chb == "/=":
					t.tokenType = tIdent

				}

				c = c + len(ch) - 1
			}
		}

		tk = append(tk, t)
	}

	if !oPL {
		return
	}

	//Clean up identities

	for c, t := range tk {

		if t.tokenType == tIdent {

			if c == 0 {
				err = errors.New(strconv.Itoa(c))
				return
			}

			if tk[c-1].tokenType != tTerm {
				err = errors.New(strconv.Itoa(c - 1))
				return
			}

			if c+1 == len(tk) {
				err = errors.New(strconv.Itoa(c + 1))
				return
			}

			if tk[c+1].tokenType != tTerm {
				err = errors.New(strconv.Itoa(c + 1))
				return
			}

			term1 := tk[c-1].str

			tk[c-1].tokenType = tIdent
			tk[c-1].str = t.str

			tk[c].tokenType = tTerm
			tk[c].str = term1

		}
	}

	return
}

func parseInfix(tk tokenStr) (n *Node, err error) {

	defer func() {

		r := recover()

		if r != nil {

			fmt.Println(tk, ": parseInfix recovered:", r)

			err = errors.New("0\nmalformed:" + tk.String() + "\nrecovered: " + fmt.Sprint(r))

		}
	}()
	var normalNodeOpen = tOpenb | tUnary | tAtomicSentence

	if oPL {
		normalNodeOpen = normalNodeOpen | tPredicate | tIdent
	}

	var expect = normalNodeOpen

	var prevNode *Node

	openbrackets := 0
	brackets := 0

	var t token

	for c := 0; c < len(tk); c++ {

		t = tk[c]

		//fmt.Println(t.tokenType, t.str)

		if t.tokenType&expect == 0 {
			err = errors.New(strconv.Itoa(t.index) + "\nunexpected symbol")
			return
		}

		switch {

		case t.tokenType == tOpenb:

			openbrackets++
			brackets++

			if prevNode.isFunctionFormula() {
				expect = tTerm
				continue
			}

			n = new(Node)
			n.index = t.index

			if prevNode != nil {
				err = prevNode.AddChild(n)
				if err != nil {
					err = errors.New(strconv.Itoa(t.index))
					return
				}
			}

			expect = normalNodeOpen

			prevNode = n

		case t.tokenType == tAtomicSentence:

			n = new(Node)
			n.index = t.index

			n.raw = t.str

			if prevNode != nil {
				err = prevNode.AddChild(n)
				if err != nil {
					err = errors.New(strconv.Itoa(t.index))
					return
				}
			}

			expect = tBinary

			e := n

			for e = n.parent; e != nil && e.IsUnary(); e = e.parent {
			}

			if e != nil && e.isSaturated() {
				expect = tCloseb
			}

			prevNode = n

		case t.isBinary():

			e := prevNode.openAncestor()

			if e == nil {

				n = new(Node)
				n.index = t.index
				n.SetConnective(t.tokenType.logicConstant())

				err = n.AddChild(prevNode.rootNode())
				if err != nil {
					err = errors.New(strconv.Itoa(t.index))
					return
				}

			} else {
				e.SetConnective(t.tokenType.logicConstant())
				e.index = t.index
				n = e
			}

			expect = normalNodeOpen

			prevNode = n

		case t.tokenType == tNeg:

			n = new(Node)
			n.index = t.index

			n.SetConnective(t.tokenType.logicConstant())

			if prevNode != nil {

				err = prevNode.AddChild(n)
				if err != nil {
					err = errors.New(strconv.Itoa(t.index))
					return
				}

			}

			expect = normalNodeOpen

			prevNode = n

		case t.isQuantifier():

			n = new(Node)
			n.index = t.index

			n.SetConnective(t.tokenType.logicConstant())

			if prevNode != nil {

				err = prevNode.AddChild(n)
				if err != nil {
					err = errors.New(strconv.Itoa(t.index))
					return
				}

			}

			expect = tTerm

			prevNode = n

		case t.tokenType == tPredicate:

			n = new(Node)
			n.index = t.index

			n.SetAtomic()

			n.predicateLetter = t.str

			n.setraw()

			if prevNode != nil {

				err = prevNode.AddChild(n)
				if err != nil {
					err = errors.New(strconv.Itoa(t.index))
					return
				}

			}

			expect = tTerm | tCloseb | tBinary

			if oPL {
				if n.isFunctionFormula() {
					expect = tOpenb | tCloseb | tBinary
				}
			}

			prevNode = n

		case t.tokenType == tIdent:

			n = new(Node)
			n.index = t.index

			n.predicateLetter = "="

			n.setraw()

			if t.str == "/=" {
				n2 := new(Node)
				n.index = t.index
				n2.SetConnective(Neg)
				n2.AddChild(n)
				n = n2
			}

			if prevNode != nil {
				prevNode.AddChild(n)
			}

			expect = tTerm

			prevNode = n

			if n.IsNegation() {
				prevNode = n.children[0]
			}

		case t.tokenType == tTerm && prevNode.isFunctionFormula():

			prevNode.term = append(prevNode.term, t.str)

			prevNode.setraw()

			expect = tComma | tCloseb

		case t.tokenType == tComma:

			expect = tTerm

		case t.tokenType == tTerm && prevNode.IsAtomic():

			prevNode.term = append(prevNode.term, t.str)

			prevNode.setraw()

			if prevNode.predicateLetter == "=" {
				if len(prevNode.term) == 2 {
					expect = tCloseb | tBinary
				} else {
					expect = tTerm
				}

			} else {
				expect = tTerm | tCloseb | tBinary
			}

		case t.tokenType == tTerm && prevNode.IsQuantifier():

			prevNode.variable = t.str

			expect = normalNodeOpen

		case t.tokenType == tCloseb:

			openbrackets--

			if openbrackets < 0 {
				err = errors.New(strconv.Itoa(t.index) + "\ntoo many closing brackets")
				return

			}

			e := prevNode.binaryAncestor()

			if e == nil {

				prevNode = prevNode.rootNode()

			} else {

				prevNode = e

			}

			expect = tCloseb | tBinary

		default:

			err = errors.New(strconv.Itoa(t.index) + "\nunknown token type: " + t.tokenType.String())
			return
		}

	}

	if openbrackets > 0 {

		err = errors.New("-1\n" + "there are unclosed brackets")
		return
	}

	//	fmt.Println("binary count", prevNode.rootNode().binaryCount())
	//	fmt.Println("bracket pairs", brackets)

	bcount := prevNode.rootNode().binaryCount()

	if bcount != 0 {
		if brackets < bcount-1 {
			err = errors.New("-1\n" + "not enough bracket pairs")
			return
		}

		if brackets > bcount {
			err = errors.New("-1\n" + "too many bracket pairs")
			return
		}
	}

	return prevNode.rootNode(), prevNode.rootNode().validate()

}
