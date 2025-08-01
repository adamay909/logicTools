package gentzen

import "strings"

func ltree(n *Node, simple bool) string {

	var lt func(n *Node) string

	disp := func(k *Node) string {
		if k.IsAtomic() {
			return k.StringF(O_Latex)
		}
		return k.MainConnective().Stringf(O_Latex)
	}

	lt = func(m *Node) (r string) {
		if simple {
			r = `[ \p{` + disp(m) + `} `
		} else {
			r = `[ \p{` + Parse(m.String(), !allowGreekUpper).StringF(O_Latex) + `} `
		}
		r = r + "\n"

		if !m.IsAtomic() {

			r = r + lt(m.subnode1)

			if m.IsBinary() {
				r = r + lt(m.subnode2)
			}

		}
		r = r + ` ] ` + "\n" // !{\qbalance} `

		return r
	}

	w := new(strings.Builder)
	w.WriteString(`\begin{forest}` + "\n")
	w.WriteString(`%syntax tree of:` + "\n% " + n.String() + "\n")
	w.WriteString(`%generated by gentzen` + "\n")
	w.WriteString(`DATA`)
	w.WriteString(`\end{forest}` + "\n\n")
	templ := w.String()

	//templ := `\begin{forest}{for tree={grow=south}}
	//%generated by gentzen
	//DATA\end{forest}

	//`
	return strings.ReplaceAll(templ, `DATA`, lt(n))

}

// SyntaxTree returns the latex code for printing the
// syntax tree of n. Each node will be printed as a full formula.
func (n *Node) SyntaxTree() string {

	return ltree(n, false)

}

// SyntaxTreeSimple is like SyntaxTree but the nodes are represented by a single connective or, in the case leadnodes, by atomic sentence letters.
func (n *Node) SyntaxTreeSimple() string {

	return ltree(n, true)

}
