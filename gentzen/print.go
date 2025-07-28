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
		return n.laTeXString()

	case O_PlainText, O_PlainASCII, O_English, O_ProofChecker:
		return n.plainString(m)

	default:
		fmt.Println(m, "UNSUPPORTED")
		os.Exit(1)
		return n.String()

	}
}

func (c LogicalConstant) isQuantifier() bool {
	return c == Uni || c == Ex
}

func (c LogicalConstant) isModalOperator() bool {
	return c == Nec || c == Pos
}

func (c LogicalConstant) isNegation() bool {
	return c == Neg
}

func (n *Node) connectiveDisplay(m PrintMode) string {

	if m == O_ProofChecker {
		m = O_PlainText
	}

	var s string

	for _, c := range connectives {
		if codeOf(n.MainConnective()) == c[0] {
			s = c[int(m)]
		}
	}
	if n.MainConnective().isQuantifier() {
		s = s + n.BoundVariable()
		if m == O_Latex {
			s = s + ` `
		}
	}

	return s
}
