package main

import (
	"math/rand"

	"github.com/adamay909/logicTools/gentzen"
)

func processFormula(l string, narrow bool) (output string, err error) {

	var tt gentzen.TruthTable

	if narrow {
		tt, err = gentzen.GenerateTruthTableNarrow(l)
	} else {
		tt, err = gentzen.GenerateTruthTable(l)

	}
	if err != nil {
		return
	}

	if narrow {
		tt.SetColumnTitles(gentzen.O_Latex)
	}

	output = output + `%%% ` + l
	table := getTable(tt)
	switch {
	case *mistakes:
		insertErrors(table, *number, l, *val)
	case *clearCol:
		clearColumns(table, *number)
	case *clearRow:
		clearRows(table, *number)
	case *clearCell:
		clearCells(table, *number)
	case *clearRandom:
		clearXY(table, *number)
	case *empty:
		emptyTable(table, l)
	case *correct:
	}

	output = output + "\n\n" + printtable(table)

	return
}

func emptyTable(table *ttable, l string) {

	v := table.NumAtomic

	if *val {
		v = 0
	}

	for i, r := range table.hide {

		for j := range r {
			if j < v {
				continue
			}
			table.hide[i][j] = true
		}
	}
}

func clearXY(table *ttable, n int) {

	if n < 1 {
		return
	}

	d1 := rand.Intn(2)

	d2 := rand.Intn(n + 1)

	if d1 == 0 {
		clearColumns(table, d2)
		return
	}

	clearRows(table, d2)

	return
}

func clearColumns(table *ttable, n int) {

	if n < 1 {
		return
	}

	if n > len(table.ColumnTitles)-1 {
		n = len(table.ColumnTitles) - 1
	}

	cs := rand.Perm(len(table.ColumnTitles))

	for hiddenCount := 0; hiddenCount < n; hiddenCount++ {

		for _, r := range table.hide {

			r[cs[hiddenCount]] = true

		}

	}
}

func clearRows(table *ttable, n int) {

	if n < 1 {
		return
	}

	if n > len(table.Rows)-1 {
		n = len(table.Rows) - 1
	}

	rs := rand.Perm(len(table.Rows) - 1)

	for hiddenCount := 0; hiddenCount < n; hiddenCount++ {

		for k := range table.hide[rs[hiddenCount]] {

			table.hide[rs[hiddenCount]][k] = true

		}
	}
}

func clearCells(table *ttable, n int) {

	if n < 1 {
		return
	}

	nrows := len(table.Rows) - 1
	ncols := len(table.ColumnTitles) - 1

	rs := rand.Perm(nrows * ncols)

	if n > len(rs) {
		n = len(rs)
	}

	for hiddenCount := 0; hiddenCount < n; hiddenCount++ {

		r := rs[hiddenCount]/ncols + 1
		c := rs[hiddenCount] % ncols
		table.hide[r][c] = true

	}

	return
}

func insertErrors(table *ttable, n int, l string, valerr bool) {

	if n < 1 {
		return
	}

	nrows := len(table.Rows) - 1
	ncols := len(table.ColumnTitles)
	if !valerr {
		ncols = ncols - table.NumAtomic
	}

	rs := rand.Perm(nrows * ncols)

	m := rand.Intn(n) + 1
	if m > len(rs) {
		m = len(rs)
	}

	errn := 0

	for i := 0; errn < m; i++ {

		r := rs[i]/ncols + 1
		c := rs[i] % ncols
		if !valerr {
			c = c + table.NumAtomic
		}
		table.mistake[r][c] = true

		errn++
		if i == len(rs)-1 {
			break
		}
	}

	return
}
