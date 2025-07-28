package gentzen

import "strings"

func derivR(d *derivNode) bool {

	if !oDR {
		logger.Print("Unknown inference rule")
		return false
	}

	var tf []string

	if len(d.supportingLines) != 1 {
		logger.Print("Derived Rule depends on a single line")
		return false
	}

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
		Debug("<--Derived Rule check: ", Parse(sn, !allowGreekUpper).StringF(O_PlainText), " against: ", Parse(thc, !allowGreekUpper).StringF(O_PlainText))
		if sameStructure(thc, sn) {
			Debug("ok")
			return true
		}
	}
	logger.Print("not valid application of ", infrule)
	return false
}
