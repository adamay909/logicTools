package main

import (
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func getTable(tt gentzen.TruthTable) *ttable {

	var tr ttable

	tr.Formula = tt.Formula
	tr.ColumnTitles = tt.ColumnTitles
	tr.Rows = tt.Rows
	tr.NumAtomic = tt.NumAtomic
	tr.Narrow = tt.Narrow
	tr.MainConnective = tt.MainConnective
	tr.Boundary = tr.Boundary

	for range tr.Rows {

		tr.hide = append(tr.hide, make([]bool, len(tr.ColumnTitles)))
		tr.mistake = append(tr.mistake, make([]bool, len(tr.ColumnTitles)))

	}

	return &tr

}

func printtable(tt *ttable) string {

	out := ""

	if tt.Narrow {
		out = tableheadnarrow(tt)
	} else {
		out = tablehead(tt)
	}

	for i, r := range tt.Rows {

		for j, val := range r {
			txt := `\emph{T}`
			if val == false {
				txt = `\emph{F}`
			}
			switch {
			case tt.mistake[i][j]:
				if val == true {
					txt = `\error{\emph{F}}`
				} else {
					txt = `\error{\emph{T}}`
				}

			case tt.hide[i][j]:
				txt = `\cover{\textcircled{` + txt + `}}`

			}

			out = out + txt
			if j != len(r)-1 {
				out = out + ` & `
			}
		}

		out = out + `\\` + "\n"
		out = out + `\hdashline` + "\n"
	}
	out = out + `\end{tabular}` + "\n\n"

	return out

}

func tablehead(tt *ttable) string {

	out := "% " + tt.Formula + "\n"

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
		out = out + `\p{` + gentzen.Parse(f, false).StringF(gentzen.O_Latex) + `}`
		if i != len(tt.ColumnTitles)-1 {
			out = out + ` & `
		}
	}

	out = out + `\\` + "\n"
	out = out + `\hline` + "\n"

	return out
}

func tableheadnarrow(tt *ttable) string {

	w := new(strings.Builder)

	w.WriteString("% " + tt.Formula + "\n")

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

	return w.String()
}
