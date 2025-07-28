package gentzen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	allowSubscr = true

	allowNumeral = true

	allowGreekUpper = true

	allowIdentity = true

	allowSpecial = true
)

// Parse should only be used when s is known to be well-formed. It will panic
// otherwise. permitGreekUpper controls whether to allow upper Greek letters as sentence or predicate letters.
func Parse[S ~string](s S, permitGreekUpper bool) *Node {

	n, err := ParseStrict(s, permitGreekUpper)

	if err != nil {
		panic("malformed formula: " + string(s) + " " + err.Error())
	}
	return n
}

// ParseStrict parses s and returns the top node of the resulting tree. err is non-nil iff. s is not a wff. permitGreekUpper controls whether to allow upper Greek letters as sentence or predicate letters.
func ParseStrict[str ~string](s str, permitGreekUpper bool) (n *Node, err error) {

	input := string(s)

	tk, err := tokenize(input, permitGreekUpper, !allowSpecial)

	if err != nil {
		return n, errors.New(showerrloc(err, input))
	}

	n, err = parse(tk)

	if err != nil {
		return n, errors.New(showerrloc(err, input))
	}

	return
}

func parse(tk tokenStr) (n *Node, err error) {

	defer func() {

		r := recover()

		if r != nil {

			err = errors.New("0\n" + fmt.Sprint(r))

		}
	}()

	prevNode := new(Node)

	currentNode := new(Node)

	nc := 0

	for c := 0; c < len(tk); c++ {

		t := tk[c]

		currentNode = new(Node)
		currentNode.index = t.index
		nc++

		switch {

		case t.isConnective() && !t.isQuantifier():
			currentNode.SetConnective(t.tokenType.logicConstant())

		case t.isQuantifier():
			currentNode.SetConnective(t.tokenType.logicConstant())
			c++
			if c == len(tk) {
				err = errors.New(strconv.Itoa(t.index))
				return
			}
			t = tk[c]
			if t.tokenType != tTerm {
				err = errors.New(strconv.Itoa(t.index))
				return
			}
			currentNode.variable = t.str

		case t.isPredicate():
			currentNode.SetAtomic()
			currentNode.predicateLetter = t.str
			for c++; c < len(tk); c++ {
				t = tk[c]
				if t.tokenType != tTerm {
					c--
					break
				}
				currentNode.term = append(currentNode.term, t.str)
			}

			if currentNode.predicateLetter == "=" && len(currentNode.term) != 2 {
				err = errors.New(strconv.Itoa(t.index))
				return
			}

			currentNode.setraw()

		case t.tokenType == tAtomicSentence:
			currentNode.SetFormula(t.str)
			currentNode.SetAtomic()

		case t.tokenType == tTerm:
			err = errors.New(strconv.Itoa(t.index))
			return

		default:
			err = errors.New(strconv.Itoa(t.index))
			return
		}

		if nc == 1 {
			prevNode = currentNode
			continue
		}

		if !prevNode.isSaturated() {
			err = prevNode.AddChild(currentNode)
		} else {
			err = prevNode.openAncestor().AddChild(currentNode)
		}

		if err != nil {
			err = errors.New(strconv.Itoa(t.index))
			return
		}

		prevNode = currentNode
	}

	n = prevNode.rootNode()

	err = n.validate()

	if err != nil {
		return
	}

	return
}

func tokenize(s string, permitGreekUpper bool, permitSpecial bool) (t tokenStr, err error) {

	t = make(tokenStr, 0, len(s))

	for i := 0; i < len(s); i++ {

		var e token

		switch {

		case strings.HasPrefix(s[i:], ldisj):
			e.tokenType = tDisj
			e.str = ldisj
		case strings.HasPrefix(s[i:], lconj):
			e.tokenType = tConj
			e.str = lconj
		case strings.HasPrefix(s[i:], lcond):
			e.tokenType = tCond
			e.str = lcond
		case strings.HasPrefix(s[i:], lneg):
			e.tokenType = tNeg
			e.str = lneg
		case !oPL:
			e.tokenType = tAtomicSentence
			ch, err := getFirstChar(s[i:], allowSubscr, !allowNumeral, permitGreekUpper, !allowIdentity, permitSpecial)
			if err != nil {
				err = errors.New(strconv.Itoa(i))
				return t, err
			}
			e.str = ch
			i = i + len(ch) - 1

		case oPL:
			switch {

			case strings.HasPrefix(s[i:], luni):
				e.tokenType = tUni
				e.str = luni
			case strings.HasPrefix(s[i:], lex):
				e.tokenType = tEx
				e.str = lex
			default:
				ch, err := getFirstChar(s[i:], allowSubscr, !allowNumeral, permitGreekUpper, allowIdentity, permitSpecial)
				if err != nil {
					err = errors.New(strconv.Itoa(i))
					return t, err
				}

				e.str = ch

				chb := strings.Split(ch, "_")[0]

				switch {
				case isRomanLower(rune(chb[0])):
					e.tokenType = tTerm

				case isRomanUpper(rune(chb[0])):
					e.tokenType = tPredicate

				case isGreekLower(chb):
					e.tokenType = tPredicate

				case isGreekUpper(chb):
					e.tokenType = tAtomicSentence

				case chb == "=":
					e.tokenType = tPredicate

				case chb == "@":
					e.tokenType = tPredicate

				case chb == "#":
					e.tokenType = tTerm

				case chb == "$":
					e.tokenType = tPredicate
				}
				i = i + len(ch) - 1
			}
		}

		e.index = i
		t = append(t, e)

	}
	return
}

func getFirstChar(s string, permitSubscr bool, permitNumeral bool, permitGreekUpper bool, permitIdentity bool, permitSpecial bool) (ch string, err error) {

	if len(s) < 1 {
		err = errors.New("malformed: too short")
		return
	}

	i := 0

	if isRoman(rune(s[0])) {

		ch = s[:1]

	}

	if len(s) > 1 && isGreekLower(s[:2]) {
		i++
		ch = s[:2]

	}

	if permitGreekUpper {

		if len(s) > 1 && isGreekUpper(s[:2]) {
			i++
			ch = s[:2]

		}
	}

	if permitNumeral {

		for j := 1; j < len(s)+1; j++ {
			_, err = strconv.Atoi(s[:j])
			if err != nil {
				break
			}
			ch = s[:j]
		}
		err = nil
	}

	if permitSpecial {

		if s[:1] == "@" { //for atomic sentences (SL) and predicates (PL)
			ch = s[:1]
		}

		if s[:1] == "#" { //for  terms
			ch = s[:1]

		}

		if s[:1] == "$" { //for function formulas
			ch = s[:1]
		}
	}

	if permitIdentity {

		if rune(s[0]) == '=' {
			ch = s[:1]
		}

		if strings.HasPrefix(s, "/=") {
			ch = s[:2]
		}
	}

	if !permitSubscr {
		return
	}

	if ch == "" {
		err = errors.New("malformed: no permissible letter string found")
		return
	}

	if len(s[i:]) > 2 {
		if s[i+1] == '_' {

			sub, err := getFirstChar(s[i+2:], !allowSubscr, allowNumeral, allowGreekUpper, !allowIdentity, !allowSpecial)

			if err != nil {
				return "", err
			}

			ch = ch + "_" + sub
		}
	}

	return

}

const (
	romanLower = "abcdefghijklmnopqrstuvwxyz"
	romanUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	greekLower = "abgdezhtiklmnxoprsyufcqw"
	greekUpper = "GDTLXPRSUFQW"
	numeral    = "0123456789"
)

func isNumeral(r rune) bool {

	return strings.ContainsRune(numeral, r)

}

func isRomanLower(r rune) bool {

	return strings.ContainsRune(romanLower, r)

}

func isRomanUpper(r rune) bool {

	return strings.ContainsRune(romanUpper, r)

}

func isRoman(r rune) bool {

	return isRomanLower(r) || isRomanUpper(r)

}

// s must be two ASCII code points
func isGreekUpper(s string) bool {

	for _, e := range greekUCBindings {

		if s == e[0] {
			return true
		}
	}

	return false
}

// s must be two ASCII code points
func isGreekLower(s string) bool {

	for _, e := range greekLCBindings {

		if s == e[0] {
			return true
		}
	}

	return false
}

func isGreek(s string) bool {

	return isGreekLower(s) || isGreekUpper(s)

}

func showerrloc(err error, s string) string {

	oerr := strings.Split(err.Error(), "\n")

	msg := "malformed:\n"
	if len(oerr) > 1 {
		msg = oerr[1] + "\n"
	}

	pos, _ := strconv.Atoi(oerr[0])

	if pos < 0 {
		return msg + s + "\n"
	}

	msg = msg + s + "\n"

	for c := 0; c < pos; c++ {
		msg = msg + " "
	}
	msg = msg + "\u25b2" + "\n"

	return msg
}
