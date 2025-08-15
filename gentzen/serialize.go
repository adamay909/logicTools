package gentzen

import (
	"strings"
)

func polishIngressFunc(n *Node, w *strings.Builder) {

	if n.Val.connective != None {
		w.WriteString(codeOf(n.Val.connective))
		if n.Val.connective.isQuantifier() {
			w.WriteString(n.Val.variable)
		}
		return
	}

	w.WriteString(n.Val.raw)

}

func polishPivotFunc(n *Node, w *strings.Builder) {
	return
}

func polishEgressFunc(n *Node, w *strings.Builder) {
	return
}

func latexIngressFunc(n *Node, w *strings.Builder) {

	if n.Val.connective == None {
		w.WriteString(predicateString(n, O_Latex))
		return
	}

	if isBinary(n) {

		if n.Parent() == nil {
			return
		}

		d := binaryHeight(n)

		if d > len(brackets)-1 {
			d = len(brackets) - 1
		}

		w.WriteString(brackets[d][0])

		return
	}

	if n.Check(isNegation) && n.Children()[0].Val.predicateLetter == "=" {

		w.WriteString(`\nident{` + encodeString(n.Children()[0].Val.term[0], O_Latex) + `}` + `{` + encodeString(n.Children()[0].Val.term[1], O_Latex) + `} `)
		return
	}

	if n.Parent() == nil {
		w.WriteString(`\mc{`)
	}

	w.WriteString(connectiveDisplay(n, O_Latex))

	if n.Parent() == nil {
		w.WriteString(`}`)
	}

	w.WriteString(` `)
}

func latexPivotFunc(n *Node, w *strings.Builder) {

	if n.Check(isBinary) {

		if n.Parent() == nil {
			w.WriteString(`\mc{`)
		}

		w.WriteString(connectiveDisplay(n, O_Latex))

		if n.Parent() == nil {
			w.WriteString(`}`)
		}
	}
}

func latexEgressFunc(n *Node, w *strings.Builder) {

	if !n.Check(isBinary) {
		return
	}

	if n.Parent() == nil {
		return
	}

	d := binaryHeight(n)

	if d > len(brackets)-1 {
		d = len(brackets) - 1
	}

	w.WriteString(brackets[d][1])
}

func laTeXString(n *Node) string {

	w := new(strings.Builder)

	ingressFunc := func(e *Node) {
		latexIngressFunc(e, w)
	}

	pivotFunc := func(e *Node) {
		latexPivotFunc(e, w)
	}

	egressFunc := func(e *Node) {
		latexEgressFunc(e, w)
	}

	n.Walk(ingressFunc, pivotFunc, egressFunc)

	return w.String()

}

func predicateString(n *Node, mode PrintMode) string {

	if mode == O_ProofChecker {
		mode = O_PlainText
	}

	if n.Val.connective != None {
		return ""
	}

	setraw(n)

	if !oPL {
		return encodeString(n.Val.raw, mode)
	}

	if !n.Check(isPredicate) {
		return encodeString(n.Val.raw, mode)
	}

	if isRomanUpper(rune(n.Val.predicateLetter[0])) {
		return encodeString(n.Val.raw, mode)
	}

	resp := ""

	if rune(n.Val.predicateLetter[0]) == '=' {
		if n.Parent() != nil && n.Parent().Check(isNegation) {
			return resp
		}

		switch mode {

		case O_Latex:
			resp = `\ident{` + encodeString(n.Val.term[0], mode) + `}` + `{` + encodeString(n.Val.term[1], mode) + `}`

		case O_PlainText:
			resp = encodeString(n.Val.term[0], mode) + "=" + encodeString(n.Val.term[1], mode)

		default:
			resp = n.Val.raw
		}
		return resp
	}

	resp = encodeString(n.Val.predicateLetter, mode)

	if len(n.Val.term) > 0 {

		var terms []string

		for _, e := range n.Val.term {
			terms = append(terms, encodeString(e, mode))
		}
		resp = resp + `(`

		resp = resp + strings.Join(terms, `,`)

		resp = resp + `)`

	}

	return resp
}

func encodeString(s string, m PrintMode) string {

	mode := 0
	if m == O_Latex {
		mode = 1
	}

	if m == O_PlainText {
		mode = 2
	}

	resp := ""
	found := false

	for c := 0; c < len(s); c++ {

		found = false

		for _, e := range greekUCBindings {
			if strings.HasPrefix(s[c:], e[0]) {
				resp = resp + e[mode]
				c = c + len(e[0]) - 1
				found = true
				break
			}
		}

		if found {
			continue
		}

		for _, e := range greekLCBindings {
			if strings.HasPrefix(s[c:], e[0]) {
				resp = resp + e[mode]
				c = c + len(e[0]) - 1
				found = true
				break
			}
		}

		if found {
			continue
		}

		resp = resp + s[c:c+1]
	}

	return resp
}

func nomarkupIngressFunc(n *Node, w *strings.Builder, mode PrintMode) {

	if n.Val.connective == None {
		w.WriteString(predicateString(n, O_PlainText))
		return
	}

	if n.Check(isBinary) {

		var br [][2]string

		switch mode {

		case O_ProofChecker:
			br = proofCheckerBrackets

		case O_PlainText:
			br = simpleBrackets

		case O_English:
			br = textBrackets

		case O_PlainASCII:
			br = plainBrackets

		}

		if n.Parent() != nil {
			d := binaryHeight(n)

			if d > len(br)-1 {
				d = len(br) - 1
			}

			w.WriteString(br[d][0])
		}

		if mode == O_English {
			if n.Check(isConditional) {
				w.WriteString("if ")
			}
		}

		return
	}

	if n.Check(isNegation) && n.Children()[0].Val.predicateLetter == "=" {

		if mode == O_PlainText {
			w.WriteString(toUnicodeString(n.Children()[0].Val.term[0]) + "≠" + toUnicodeString(n.Children()[0].Val.term[1]))
		} else {
			w.WriteString(toUnicodeString(n.Children()[0].Val.term[0]) + "/=" + toUnicodeString(n.Children()[0].Val.term[1]))
		}
		return
	}

	w.WriteString(connectiveDisplay(n, mode))

}

func nomarkupPivotFunc(n *Node, w *strings.Builder, mode PrintMode) {

	if !isBinary(n) {
		return
	}

	if mode == O_ProofChecker {
		mode = O_PlainText
	}

	if mode == O_English {
		if isConditional(n) {
			w.WriteString(`, then `)
			return
		}
	}

	w.WriteString(connectiveDisplay(n, mode))

}

func nomarkupEgressFunc(n *Node, w *strings.Builder, mode PrintMode) {

	if !isBinary(n) {
		return
	}

	var br [][2]string

	switch mode {

	case O_ProofChecker:
		br = proofCheckerBrackets

	case O_PlainText:
		br = simpleBrackets

	case O_English:
		br = textBrackets

	case O_PlainASCII:
		br = plainBrackets

	}

	if n.Parent() != nil {
		d := binaryHeight(n)

		if d > len(br)-1 {
			d = len(br) - 1
		}

		w.WriteString(br[d][1])
	}
}

func plainString(n *Node, m PrintMode) string {

	w := new(strings.Builder)

	ingressFunc := func(n *Node) {
		nomarkupIngressFunc(n, w, m)
	}

	pivotFunc := func(n *Node) {
		nomarkupPivotFunc(n, w, m)
	}

	egressFunc := func(n *Node) {
		nomarkupEgressFunc(n, w, m)
	}

	n.Walk(ingressFunc, pivotFunc, egressFunc)

	return w.String()

}

func toUnicodeString(s string) string {

	resp := ""
	found := false

	for c := 0; c < len(s); c++ {

		found = false

		for _, e := range greekUCBindings {
			if strings.HasPrefix(s[c:], e[0]) {
				resp = resp + e[2]
				c = c + len(e[0]) - 1
				found = true
				break
			}
		}

		if found {
			continue
		}

		for _, e := range greekLCBindings {
			if strings.HasPrefix(s[c:], e[0]) {
				resp = resp + e[2]
				c = c + len(e[0]) - 1
				found = true
				break
			}
		}

		if found {
			continue
		}

		resp = resp + s[c:c+1]
	}

	return resp
}
