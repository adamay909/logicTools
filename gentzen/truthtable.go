package gentzen

import (
	"slices"
	"strconv"
	"strings"
)

// TruthTable holds information about a truth table.
type TruthTable struct {
	Formula        string
	ColumnTitles   []string
	Rows           [][]bool
	NumAtomic      int   //number of atomic sentences
	Narrow         bool  //true if the truth table uses the narrow format
	MainConnective int   //index within infix tokenstring of formula (used for narrow fomat
	Boundary       []int //index of boundaries between tokens to determine column titles for narrow format
}

func isTautologyTT(s string) bool {

	tt, err := GenerateTruthTable(s)

	if err != nil {
		return false
	}

	last := len(tt.ColumnTitles) - 1

	for _, row := range tt.Rows {

		if row[last] == false {
			return false
		}
	}

	return true
}

func isTautologyTTfast(s string) bool {

	_, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {
		panic(err)
	}

	cols := getColumnTitles(s)

	numAtomic := len(getAtomicColumns(cols))

	last := len(cols) - 1

	for rownum := 0; rownum < 1<<numAtomic; rownum++ {

		if rowValues(rownum, cols)[last] == false {
			return false
		}

	}

	return true
}

/*
GenerateTruthTable returns the TruthTable of s where the table
takes the wide format. For narrow format truth tables, use
[GenerateTruthTableNarrow]. Here is a wide truth table for [(p⊃q)∧¬p]⊃¬q presented in O_Polish mode:

	   	   p  q  |  Cpq  Np  KCpqNp  Nq  CKCpqNpNq
		  -------+---------------------------------
	  	   T  T  |   T   F     F     F       T
		   F  T  |   T   T     T     F       F
		   T  F  |   F   F     F     T       T
		   F  F  |   T   T     T     T       T
*/
func GenerateTruthTable(s string) (tt TruthTable, err error) {

	_, err = ParseStrict(s, !allowGreekUpper)

	if err != nil {
		return
	}

	tt.ColumnTitles = getColumnTitles(s)

	tt.NumAtomic = len(getAtomicColumns(tt.ColumnTitles))

	for rownum := 0; rownum < 1<<tt.NumAtomic; rownum++ {

		tt.Rows = append(tt.Rows, rowValues(rownum, tt.ColumnTitles))

	}

	return
}

func valuation(rownumber, numAtomic int) (val []bool, err error) {

	var v int
	for i := 0; i < numAtomic; i++ {

		v = rownumber

		v = v >> i

		val = append(val, v%2 == 0)

	}

	return
}

func getColumnTitles(s string) []string {

	n, err := ParseStrict(s, !allowGreekUpper)

	if err != nil {
		panic(s + "is not well-formed")
	}

	nslice := linearizeBT(n)

	var atomicCol []string

	var regCol []string

	for _, n = range nslice {

		nstr := n.String()

		if n.IsAtomic() {

			if !slices.Contains(atomicCol, nstr) {
				atomicCol = append(atomicCol, nstr)
			}
		} else {
			if !slices.Contains(regCol, nstr) {
				regCol = append(regCol, nstr)
			}
		}
	}

	slices.SortFunc(atomicCol, sortAtomicSentences)

	return slices.Concat(atomicCol, regCol)
}

func getAtomicColumns(column []string) []string {

	var resp []string

	for _, c := range column {

		n, err := ParseStrict(c, !allowGreekUpper)

		if err != nil {
			panic(c + "is not well-formed")
		}

		if n.IsAtomic() {

			resp = append(resp, c)

		}

	}

	return resp

}

func rowValues(rownumber int, columns []string) []bool {

	colval := make(map[string]bool)

	var resp []bool

	ac := getAtomicColumns(columns)

	val, _ := valuation(rownumber, len(ac))

	for i, a := range ac {

		colval[a] = val[i]

		resp = append(resp, colval[a])
	}

	var v1, v2, ok bool

	for i := len(ac); i < len(columns); i++ {

		n, err := ParseStrict(columns[i], !allowGreekUpper)

		if err != nil {

			panic(columns[i] + " is not well formed")

		}

		v1, ok = colval[n.children[0].String()]
		if !ok {
			panic("NEG. something wrong " + n.children[0].String())
		}

		if n.IsBinary() {

			v2, ok = colval[n.children[1].String()]

			if !ok {
				panic("BINARY. something wrong " + n.children[1].String())
			}

		}

		switch n.MainConnective() {

		case Neg:

			colval[n.String()] = !v1

		case Conj:

			colval[n.String()] = v1 && v2

		case Disj:

			colval[n.String()] = v1 || v2

		case Cond:

			colval[n.String()] = !v1 || v2

		}

		resp = append(resp, colval[n.String()])

	}

	return resp

}

// PrintTruthTable prints the truth table.
func (tt *TruthTable) PrintTruthTable(mode PrintMode, rowsep bool) string {

	if tt.Narrow {
		return printTruthTableNarrow(tt, mode, rowsep)
	}

	if mode == O_Latex {
		return printTruthTableLatex(tt, rowsep)
	}

	w := new(strings.Builder)

	numCols := len(tt.ColumnTitles)
	colWidths := make([]int, numCols)

	//title line
	for i, f := range tt.ColumnTitles {

		n, _ := ParseStrict(f, !allowGreekUpper)

		title := n.StringF(mode)

		colWidths[i] = len([]rune(title))

		w.WriteString(center(title, colWidths[i]+2))

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

// center centers text within a given width
func center(text string, width int) string {
	if len([]rune(text)) >= width {
		return text
	}
	padding := width - len([]rune(text))
	leftPad := padding / 2
	rightPad := padding - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// PrintTruthTableLatex returns LaTeX code for printing tt in wide format.
func printTruthTableLatex(tt *TruthTable, rowsep bool) string {

	out := `%truth table of:` + "\n"

	out = out + "% " + tt.ColumnTitles[len(tt.ColumnTitles)-1] + "\n"

	out = out + `%generated by gentzen` + "\n"

	out = out + `\begin{tabular}{`

	for i := 0; i < len(tt.ColumnTitles); i++ {
		out = out + `c`
		if i < tt.NumAtomic-1 {
			continue
		}

		if i > tt.NumAtomic-2 && i < len(tt.ColumnTitles)-1 {
			out = out + `|`
		}

		if i == len(tt.ColumnTitles)-2 {
			out = out + `|`
		}
	}

	out = out + "}" + "\n"

	for i, f := range tt.ColumnTitles {
		out = out + `\p{` + Parse(f, !allowGreekUpper).StringF(O_Latex) + `}`
		if i != len(tt.ColumnTitles)-1 {
			out = out + ` & `
		}
	}

	out = out + `\\` + "\n"
	out = out + `\hline` + "\n"

	for _, r := range tt.Rows {

		for j, val := range r {
			txt := "T"
			if val == false {
				txt = "F"
			}
			out = out + `\emph{` + txt + `}`
			if j != len(r)-1 {
				out = out + ` & `
			}
		}

		out = out + `\\` + "\n"
		if rowsep {
			out = out + `\hdashline` + "\n"
		}
	}
	out = out + `\end{tabular}` + "\n\n"

	return out

}

func sortAtomicSentences(s1, s2 string) int {

	if s1 == s2 {
		return 0
	}

	idx1 := strings.Index(s1, "_")

	idx2 := strings.Index(s2, "_")

	if idx1 == -1 || idx2 == -1 {

		if s1 < s2 {
			return -1
		}
		return 1

	}

	if s1[:idx1] < s2[:idx2] {
		return -1
	}

	if s1[:idx1] > s2[:idx2] {
		return 1
	}

	sub1, err := strconv.Atoi(s1[idx1+1:])

	if err != nil {
		if s1 < s2 {
			return -1
		}
		return 1

	}

	sub2, err := strconv.Atoi(s2[idx2+1:])

	if err != nil {
		if s1 < s2 {
			return -1
		}
		return 1

	}

	if sub1 < sub2 {
		return -1
	}

	if sub1 > sub2 {
		return 1
	}

	return 0
}

func isTautologySub(s string) bool {

	var atomic []string

	tks, err := tokenize(s, !allowGreekUpper, !allowSpecial)

	if err != nil {
		panic(err)
	}

	for _, t := range tks {

		if t.tokenType != tAtomicSentence {
			continue
		}

		if !slices.Contains(atomic, t.str) {
			atomic = append(atomic, t.str)
		}
	}

	s = strings.ReplaceAll(s, "C", "AN")
	for r := 0; r < 1<<len(atomic); r++ {

		s0 := s

		val, _ := valuation(r, len(atomic))
		for i := range atomic {
			if val[i] == true {
				s0 = strings.ReplaceAll(s0, atomic[i], "T")
			} else {
				s0 = strings.ReplaceAll(s0, atomic[i], "F")
			}
		}

		for len(s0) > 1 {

			s0 = strings.ReplaceAll(s0, "NT", "F")
			s0 = strings.ReplaceAll(s0, "NF", "T")

			s0 = strings.ReplaceAll(s0, "KTT", "T")
			s0 = strings.ReplaceAll(s0, "KTF", "F")
			s0 = strings.ReplaceAll(s0, "KFT", "F")
			s0 = strings.ReplaceAll(s0, "KFF", "F")

			s0 = strings.ReplaceAll(s0, "ATT", "T")
			s0 = strings.ReplaceAll(s0, "ATF", "T")
			s0 = strings.ReplaceAll(s0, "AFT", "T")
			s0 = strings.ReplaceAll(s0, "AFF", "F")

		}

		if s0 == "F" {
			return false
		}
	}

	return true
}
