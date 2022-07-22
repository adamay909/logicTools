package gentzen

func uniE(seq1, seq2 sequent) bool {

	if Parse(seq1.succedent().String()).MainConnective() != uni {
		logger.Print("premise must be universally quantified")
		return false
	}

	val, _, _ := isInstanceOf(seq2.succedent().String(), seq1.succedent().String())
	if !val {
		logger.Print("conclusion not an instance of premise")
	}

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

	return true
}

func uniI(seq1, seq2 sequent) bool {

	if Parse(seq2.succedent().String()).MainConnective() != uni {
		logger.Print("conclusion must be universally quantified")
		return false
	}

	val, _, term := isInstanceOf(seq1.succedent().String(), seq2.succedent().String())

	if !val {
		logger.Print("conclusion not a universalization of premise")
		return false
	}

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

	for _, d := range seq1.datumSlice() {
		if len(d) == 0 {
			continue
		}
		if string(d[0]) == `\` {
			continue
		}
		if Parse(d).hasTerm(term) {
			logger.Print(term, " cannot appear in datum")
			return false
		}
	}

	return true

}
