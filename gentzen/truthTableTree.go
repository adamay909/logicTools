package gentzen

import (
	"strings"
)

type interpretation map[string]bool

// val must be a map from atomic sentences to truth values. rownum is the
// row number given the standard way of ordering the rows of a truth table.
// It assumes that interpretations are generated in order.
func (n *Node) genInterpretation(val map[string]bool, rownum int, rowTotal int) {

	ingressFunc := func(e *Node) {
		return
	}

	pivotFunc := func(e *Node) {
		return
	}

	egressFunc := func(e *Node) {

		switch {
		case e.IsAtomic():
			e.tvalue = append(e.tvalue, val[e.String()])

		case e.IsNegation():
			e.tvalue = append(e.tvalue, !e.children[0].tvalue[rownum])

		case e.IsConjunction():
			e.tvalue = append(e.tvalue, e.children[0].tvalue[rownum] && e.children[1].tvalue[rownum])

		case e.IsDisjunction():
			e.tvalue = append(e.tvalue, e.children[0].tvalue[rownum] || e.children[1].tvalue[rownum])

		case e.IsConditional():
			e.tvalue = append(e.tvalue, !e.children[0].tvalue[rownum] || e.children[1].tvalue[rownum])
		}

	}

	Serialize(n, ingressFunc, pivotFunc, egressFunc)

}

func (n *Node) assignTruthValues() {

	ac := getAtomicColumns(getColumnTitles(n.String()))

	atomicCount := len(ac)

	value := make(map[string]bool, atomicCount)

	rowTotal := 1 << atomicCount

	for rownum := 0; rownum < rowTotal; rownum++ {

		atomicValues, _ := valuation(rownum, atomicCount)

		for i, a := range ac {

			value[a] = atomicValues[i]
		}

		n.genInterpretation(value, rownum, rowTotal)

	}
}

func constructTruthTableWide(s string) (tt TruthTable, err error) {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {

		return tt, err

		return

	}

	tt.Formula = s

	var egressFunc, donothing func(*Node)

	donothing = func(*Node) {
		return
	}

	n.assignTruthValues()

	tt.ColumnTitles = getColumnTitles(s)

	tt.NumAtomic = len(getAtomicColumns(tt.ColumnTitles))

	tt.Rows = make([][]bool, 0, 1<<tt.NumAtomic)

	var row []bool

	v := make([]bool, tt.NumAtomic)

	appeared := make(map[string]bool)

	for rownum := 0; rownum < 1<<tt.NumAtomic; rownum++ {

		row = make([]bool, 0, len(tt.ColumnTitles))

		v, _ = valuation(rownum, tt.NumAtomic)

		row = append(row, v...)

		clear(appeared)

		egressFunc = func(e *Node) {

			if e.IsAtomic() {
				return
			}

			if !appeared[e.String()] {
				row = append(row, e.tvalue[rownum])
				appeared[e.String()] = true
			}
		}

		Serialize(n, donothing, donothing, egressFunc)

		tt.Rows = append(tt.Rows, row)
	}

	return
}

/*
GenerateTruthTableNarrow return the truth table for s in narrow format.
For wide format, use [GenerateTruthTable]. Here is a 'narrow' truth table
presented in O_PlainText mode:

		p  q  |  [(p  ⊃  q)  ∧  ¬  p]  |⊃|  ¬  q
	   -------+----------------------------------
		T  T  |   T   T  T   F  F  T   |T|  F  T
		F  T  |   F   T  T   T  T  F   |F|  F  T
		T  F  |   T   F  F   F  F  T   |T|  T  F
		F  F  |   F   T  F   T  T  F   |T|  T  F

Note that the narrow table can only be presented using some version of the infix notation (O_PlainText, O_PlainASCII, O_Latex).
*/
func GenerateTruthTableNarrow(s string) (tt TruthTable, err error) {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {

		return tt, err

		return

	}

	tt.Formula = s

	n.assignTruthValues()

	tt.ColumnTitles = getColumnTitles(s)

	tt.NumAtomic = len(getAtomicColumns(getColumnTitles(s)))

	for i := tt.NumAtomic; i < len(tt.ColumnTitles); i++ {

		tt.ColumnTitles[i] = ""

	}

	tt.Rows = make([][]bool, 0, 1<<tt.NumAtomic)

	tt.Narrow = true

	var ingressFunc, pivotFunc, egressFunc func(*Node)
	var row []bool

	v := make([]bool, tt.NumAtomic)

	for rownum := 0; rownum < 1<<tt.NumAtomic; rownum++ {

		row = make([]bool, 0, len(tt.ColumnTitles))

		v, _ = valuation(rownum, tt.NumAtomic)

		row = append(row, v...)

		ingressFunc = func(e *Node) {
			if !e.IsBinary() {
				row = append(row, e.tvalue[rownum])
				if e.parent == nil {
					tt.MainConnective = len(row) - 1
				}
			}
		}

		pivotFunc = func(e *Node) {
			if e.IsBinary() {
				row = append(row, e.tvalue[rownum])
				if e.parent == nil {
					tt.MainConnective = len(row) - 1
				}
			}
		}

		egressFunc = func(e *Node) {
		}

		Serialize(n, ingressFunc, pivotFunc, egressFunc)

		tt.Rows = append(tt.Rows, row)
	}

	tks, _ := tokenizeInfix(n.StringF(O_PlainASCII), !allowGreekUpper)

	var t token

	for i := 0; i < len(tks); i++ {

		t = tks[i]

		switch t.tokenType {
		case tNeg:
			tt.Boundary = append(tt.Boundary, i+1)

		case tConj, tDisj, tCond:
			tt.Boundary = append(tt.Boundary, i, i+1)
		}
	}

	tt.Boundary = append(tt.Boundary, len(tks))

	return
}

// PrintTruthTableNarrow prints the truth table tt in narrow format. You need
// to generate tt in the narrow format (by using GenerateTruthTableNarrow). Set
// rowsep to true if you want a separator line between each row.
func printTruthTableNarrow(tt *TruthTable, mode PrintMode, rowsep bool) string {

	if !tt.Narrow {
		tt.PrintTruthTable(mode, rowsep) //this really should not happen
	}

	if mode == O_Latex {
		return printTruthTableNarrowLatex(tt, rowsep)
	}

	tt.SetColumnTitles(mode)

	w := new(strings.Builder)

	numCols := len(tt.ColumnTitles)
	colWidths := make([]int, numCols)

	for i, f := range tt.ColumnTitles {

		if i == tt.MainConnective {
			f = "|" + f + "|"
		}

		colWidths[i] = len([]rune(f))

		w.WriteString(center(f, colWidths[i]+2))

		if i == tt.NumAtomic-1 {
			w.WriteString(` | `)
		}

	}

	w.WriteString("\n")

	//separator
	for i, wdth := range colWidths {
		w.WriteString(strings.Repeat("-", wdth+2))

		if i == tt.NumAtomic-1 {
			w.WriteString(`-+-`)
		}
	}

	w.WriteString("\n")

	//the actual table

	for rownum, row := range tt.Rows {
		for i, val := range row {

			text := "T"

			if val == false {
				text = "F"
			}

			if i == tt.MainConnective {
				text = "|" + text + "|"
			}

			w.WriteString(center(text, colWidths[i]+2))

			if i == tt.NumAtomic-1 {
				w.WriteString(` | `)
			}

		}
		w.WriteString("\n")

		if rowsep && rownum < len(tt.Rows)-1 {
			for i, wdth := range colWidths {
				w.WriteString(strings.Repeat("-", wdth+2))

				if i == tt.NumAtomic-1 {
					w.WriteString(`-+-`)
				}
			}
			w.WriteString("\n")
		}
	}
	return w.String()
}

// SetColumnTitles sets the column titles for tt. Only useful for narrow
// format tables.
func (tt *TruthTable) SetColumnTitles(mode PrintMode) {

	if !tt.Narrow {
		return
	}

	n, err := ParseStrict(tt.Formula, !allowGreekUpper)

	if err != nil {

		return

	}
	tt.ColumnTitles = tt.ColumnTitles[:tt.NumAtomic]

	txtTokens := infixTextTokens(n, mode)

	b0 := 0

	for _, b := range tt.Boundary {

		tt.ColumnTitles = append(tt.ColumnTitles, strings.Join(txtTokens[b0:b], ""))

		b0 = b
	}
}

// PrintTruthTableNarrowLatex returns the LaTeX code for printing  tt
// in narrow format.
func printTruthTableNarrowLatex(tt *TruthTable, rowsep bool) string {

	if !tt.Narrow {
		return ""
	}

	mode := O_Latex

	tt.SetColumnTitles(mode)

	w := new(strings.Builder)

	w.WriteString(`\begin{tabular}`)

	w.WriteString(`{`)

	for i := range tt.ColumnTitles {

		if i == tt.NumAtomic {
			w.WriteString(`||`)
		}

		w.WriteString(`c`)

	}

	w.WriteString(`}` + "\n")

	for i, f := range tt.ColumnTitles {

		w.WriteString(`\p{` + f + `}`)

		if i < len(tt.ColumnTitles)-1 {
			w.WriteString(` & `)
		} else {
			w.WriteString(` \\`)
		}
	}

	w.WriteString("\n")

	w.WriteString(`\hline`)

	w.WriteString("\n")

	//the actual table

	for _, row := range tt.Rows {
		for i, val := range row {

			text := "T"

			if val == false {
				text = "F"
			}

			if i == tt.MainConnective {
				text = `\textbf ` + text
			}

			w.WriteString(`\emph{` + text + `}`)

			if i != len(row)-1 {
				w.WriteString(" & ")
			} else {
				w.WriteString(` \\`)
			}
		}
		w.WriteString("\n")
	}

	w.WriteString(`\end{tabular}`)
	w.WriteString("\n")

	return w.String()
}
