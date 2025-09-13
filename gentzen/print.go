package gentzen

import (
	"fmt"
	"os"
)

// PrintMode specifies the style of output of Stringer and similiar methods.
type PrintMode int

func printNodeInfix(n *Node, m PrintMode) string {

	switch m {

	case O_Latex:
		return laTeXString(n)

	case O_PlainText, O_PlainASCII, O_English, O_ProofChecker:
		return plainString(n, m)

	default:
		fmt.Println(m, "UNSUPPORTED")
		os.Exit(1)
		return n.String()

	}
}

func (c logicalConstant) isQuantifier() bool {
	return c == uni || c == ex
}

func (c logicalConstant) isModalOperator() bool {
	return c == nec || c == pos
}

func (c logicalConstant) isNegation() bool {
	return c == neg
}

func connectiveDisplay(n *Node, m PrintMode) string {

	if m == O_ProofChecker {
		m = O_PlainText
	}

	var s string

	for _, c := range connectives {
		if codeOf(n.Val.connective) == c[0] {
			s = c[int(m)]
		}
	}
	if n.Val.connective.isQuantifier() {
		s = s + n.Val.variable
		if m == O_Latex {
			s = s + ` `
		}
	}

	return s
}
