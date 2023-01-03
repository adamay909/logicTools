package gentzen

import "strings"

func derivR(infrule string, seq1, seq2 sequent) bool {

	var tf []string

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
		Debug("<--Derived Rule check: ", Parse(sn).StringPlain(), " against: ", Parse(thc).StringPlain())
		if sameStructure(thc, sn) {

			Debug("ok")
			Debug("--done derived rule check-->")

			return true
		}
	}
	logger.Print("not valid application of ", infrule)

	Debug("fail")
	Debug("--done derived rule check-->")

	return false
}
