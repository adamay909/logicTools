//This is for parsing logical sentences given in Polish notation

package gentzen

import (
	"errors"
	"log"
)

func parseTokens(t tokenStr) (*Node, error) {

	var n Node
	var err error

	if len(t) == 0 {
		err = errors.New("too short")

		return &n, err
	}

	if !t[0].isConnective() {

		if len(t) > 1 {
			err = errors.New("malformed: " + t.String())
			return &n, err
		}
		n.SetFormula(t.String())
		n.SetAtomic()
		n.predicateLetter = t[0].predicate
		n.term = t[0].term
		return &n, err
	}

	if len(t) < 2 {
		err = errors.New("malformed: " + t.String())
	}
	n.SetFormula(t.String())
	switch {
	case t[0].isUnary():
		n.SetConnective(t[0].tokenType.logicConstant())
		n.variable = t[0].variable
		s1 := findNextSentence(t[1:])
		n.subnode1, err = parseTokens(s1)
		if err != nil {
			return &n, err
		}
		n.subnode1.parent = &n

	case t[0].isBinary():
		n.SetConnective(t[0].tokenType.logicConstant())
		s1 := findNextSentence(t[1:])
		n.subnode1, err = parseTokens(s1)
		n.subnode1.parent = &n
		s2 := t[len(s1)+1:]
		n.subnode2, err = parseTokens(s2)
		if err != nil {
			return &n, err
		}
		n.subnode2.parent = &n

	default:
		return &n, err
	}

	return &n, err
}

func findNextSentence(s tokenStr) tokenStr {

	if len(s) == 0 {
		return s
	}
	i, m, k := 0, 0, 0

	for i = range s {
		if !s[i].isConnective() {
			m++
			if m == k+1 {
				break
			}
			continue
		}
		if s[i].isKCA() {
			k++
			continue
		}
	}

	if k == 0 {
		return s[:i+1]
	}

	if m != k+1 {
		log.Fatal("malformed", s)
	}

	return s[:i+1]

}
