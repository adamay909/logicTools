package gentzen

import (
	"strings"
)

/*
Serialize walks depth-first and left-branch-first the tree rooted at n.
The three functions specify what to do as the process enters a node,
as it switches to a sibling, and as it is done with the node.

Here is an example of the use of Serialize. You can serialize the nodes
as a one-dimensional slice like this:

	func order(n *Node) []*Node {

	 var response []*Node

	 infunc := func (n *Node) {
		response = append(response, n)
	   }

	 pivotfunc := func (n *Node) {
	   }

	 efunc := func (n *Node) {
	   }

	  return Serialize(n, infunc, pivotfunc, func)
	}

You can get a one-dimensional slice in a reverse Polish style like this:

	func order(n *Node) []*Node {

	 var response []*Node

	 infunc := func (n *Node) {
	   }

	 pivotfunc := func (n *Node) {
	   }

	 efunc := func (n *Node) {
		response = append(response, n)
	   }

	  return Serialize(n, infunc, pivotfunc, func)
	}
*/
func Serialize(n *Node, ingressFunc, pivotFunc, egressFunc func(*Node)) {

	var walk func(*Node)

	walk = func(e *Node) {

		if e == nil {
			return
		}

		ingressFunc(e)

		for i, c := range e.children {

			walk(c)

			if i < len(e.children)-1 {
				pivotFunc(e)
			}

		}

		egressFunc(e)

	}

	walk(n)

}

func linearize(n *Node) []*Node {

	var resp []*Node

	ingressFunc := func(e *Node) {

		if e == nil {
			return
		}

		resp = append(resp, e)

	}

	pivotFunc := func(e *Node) {
		return
	}

	egressFunc := func(e *Node) {
		return
	}

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

	return resp

}

// linearize postfix style
func linearizeBT(n *Node) []*Node {

	var resp []*Node

	egressFunc := func(e *Node) {

		if e == nil {
			return
		}

		resp = append(resp, e)

	}

	pivotFunc := func(e *Node) {
		return
	}

	ingressFunc := func(e *Node) {
		return
	}

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

	return resp

}

func polishIngressFunc(n *Node, w *strings.Builder) {

	if n.IsConnective() {
		w.WriteString(codeOf(n.connective))
		if n.IsQuantifier() {
			w.WriteString(n.variable)
		}
		return
	}

	w.WriteString(n.raw)

}

func polishPivotFunc(n *Node, w *strings.Builder) {
	return
}

func polishEgressFunc(n *Node, w *strings.Builder) {
	return
}

func latexIngressFunc(n *Node, w *strings.Builder) {

	if !n.IsConnective() {
		w.WriteString(n.predicateString(O_Latex))
		return
	}

	if n.IsBinary() {

		if n.parent == nil {
			return
		}

		d := n.nestingDepth()

		if d > len(brackets)-1 {
			d = len(brackets) - 1
		}

		w.WriteString(brackets[d][0])

		return
	}

	if n.IsNegation() && n.children[0].predicateLetter == "=" {

		w.WriteString(`\nident{` + encodeString(n.children[0].term[0], O_Latex) + `}` + `{` + encodeString(n.children[0].term[1], O_Latex) + `} `)
		return
	}

	if n.parent == nil {
		w.WriteString(`\mc{`)
	}

	w.WriteString(n.connectiveDisplay(O_Latex))

	if n.parent == nil {
		w.WriteString(`}`)
	}

	w.WriteString(` `)
}

func latexPivotFunc(n *Node, w *strings.Builder) {

	if n.IsBinary() {

		if n.parent == nil {
			w.WriteString(`\mc{`)
		}

		w.WriteString(n.connectiveDisplay(O_Latex))

		if n.parent == nil {
			w.WriteString(`}`)
		}
	}
}

func latexEgressFunc(n *Node, w *strings.Builder) {

	if !n.IsBinary() {
		return
	}

	if n.parent == nil {
		return
	}

	d := n.nestingDepth()

	if d > len(brackets)-1 {
		d = len(brackets) - 1
	}

	w.WriteString(brackets[d][1])
}

func (n *Node) laTeXString() string {

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

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

	return w.String()

}

func (n *Node) predicateString(mode PrintMode) string {

	if mode == O_ProofChecker {
		mode = O_PlainText
	}

	if n.IsConnective() {
		return ""
	}

	n.setraw()

	if !oPL {
		return encodeString(n.raw, mode)
	}

	if !n.IsPredicate() {
		return encodeString(n.raw, mode)
	}

	if isRomanUpper(rune(n.predicateLetter[0])) {
		return encodeString(n.raw, mode)
	}

	resp := ""

	if rune(n.predicateLetter[0]) == '=' {
		if n.parent != nil && n.parent.IsNegation() {
			return resp
		}

		switch mode {

		case O_Latex:
			resp = `\ident{` + encodeString(n.term[0], mode) + `}` + `{` + encodeString(n.term[1], mode) + `}`

		case O_PlainText:
			resp = encodeString(n.term[0], mode) + "=" + encodeString(n.term[1], mode)

		default:
			resp = n.raw
		}
		return resp
	}

	resp = encodeString(n.predicateLetter, mode)

	if len(n.term) > 0 {

		var terms []string

		for _, e := range n.term {
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

	if !n.IsConnective() {
		w.WriteString(n.predicateString(O_PlainText))
		return
	}

	if n.IsBinary() {

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

		if n.parent != nil {
			d := n.nestingDepth()

			if d > len(br)-1 {
				d = len(br) - 1
			}

			w.WriteString(br[d][0])
		}

		if mode == O_English {
			if n.IsConditional() {
				w.WriteString("if ")
			}
		}

		return
	}

	if n.IsNegation() && n.children[0].predicateLetter == "=" {

		if mode == O_PlainText {
			w.WriteString(toUnicodeString(n.children[0].term[0]) + "≠" + toUnicodeString(n.children[0].term[1]))
		} else {
			w.WriteString(toUnicodeString(n.children[0].term[0]) + "/=" + toUnicodeString(n.children[0].term[1]))
		}
		return
	}

	w.WriteString(n.connectiveDisplay(mode))

}

func nomarkupPivotFunc(n *Node, w *strings.Builder, mode PrintMode) {

	if !n.IsBinary() {
		return
	}

	if mode == O_ProofChecker {
		mode = O_PlainText
	}

	if mode == O_English {
		if n.IsConditional() {
			w.WriteString(`, then `)
			return
		}
	}

	w.WriteString(n.connectiveDisplay(mode))

}

func nomarkupEgressFunc(n *Node, w *strings.Builder, mode PrintMode) {

	if !n.IsBinary() {
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

	if n.parent != nil {
		d := n.nestingDepth()

		if d > len(br)-1 {
			d = len(br) - 1
		}

		w.WriteString(br[d][1])
	}
}

func (n *Node) plainString(m PrintMode) string {

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

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

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
