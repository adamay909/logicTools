package gentzen

import "strings"

func derivR(d *derivNode) bool {

	var tf []string

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq
	infrule := d.line.inf
	thms := theoremsInUse()

	//check name of derived rule
	t := strings.TrimSuffix(strings.TrimSpace(infrule), "R")

	for i := range thms {
		if t == thms[i][1] {
			tf = append(tf, thms[i][2])
		}
	}
	if len(tf) == 0 {
		logger.Print(infrule, "does not match any theorems")
		return false
	}

	if len(d.supportingLines) != 1 {
		logger.Print("Derived Rule depends on a single line")
		return false
	}

	//check datums
	if strictCheck {
		if !datumsEqual(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !datumsEquiv(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum cannot change")
			return false
		}
	}

	//check for match with theorem
	s1 := seq1.succedent().String()
	s2 := seq2.succedent().String()

	sn := `>` + s1 + s2

	for _, thc := range tf {
		Debug("<--Derived Rule check: ", Parse(sn).StringPlain(), " against: ", Parse(thc).StringPlain())
		if sameStructure(thc, sn) {
			Debug("ok")
			return true
		}
	}
	logger.Print("not valid application of ", infrule)
	return false
}
