package gentzen

func negE(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())

	if n1.MainConnective() != neg {
		logger.Print("premise must be double negation")
		return false
	}

	if n1.subnode1.MainConnective() != neg {
		logger.Print("premise must be double negation")
		return false
	}

	if n1.subnode1.subnode1.Formula() != n2.Formula() {
		logger.Print("conclusion is not the elimnation of double negation")
		return false
	}
	if strictCheck {
		if !datumsEqual(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum must remain same")
			return false
		}
	} else {
		if !datumsEquiv(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum must remain same")
			return false
		}
	}

	return true
}

func negI(seq1, seq2, seq3 sequent) bool {

	v1, msg1 := negIhelper(seq1, seq2, seq3)

	v2, msg2 := negIhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "conclusion must be negation" {
		logger.Print(msg2)
		return false
	}
	logger.Print(msg1)

	return false
}

func negIhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	v = false
	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())
	n3 := Parse(seq3.succedent().String())

	if n3.MainConnective() != neg {
		msg = "conclusion must be negation"
		return
	}

	if n1.MainConnective() != neg {
		msg = `rule misapplied`
		return
	}

	if n1.subnode1.Formula() != n2.Formula() {
		msg = `rule misapplied`
		return
	}

	f := datum(n3.subnode1.Formula())

	if !datumIncludes(seq1.datumSlice(), f) {
		msg = "conclusion must be negation of something in common between the datums of premises"

		return
	}

	if !datumIncludes(seq2.datumSlice(), f) {
		msg = "conclusion must be in datums of both premises"
		return
	}

	canonicalDatum := datumUnion(datumRm(seq1.datumSlice(), f), datumRm(seq2.datumSlice(), f))

	if strictCheck {
		if !datumsEqual(canonicalDatum, seq3.datumSlice()) {
			msg = "datum of conclusion incorrect"
			return
		}
	} else {
		if !datumsEquiv(canonicalDatum, seq3.datumSlice()) {
			msg = "datum of conclusion incorrect"
			return
		}
	}
	v = true
	return
}
