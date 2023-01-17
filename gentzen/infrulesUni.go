package gentzen

func uniE(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Universal Elimination depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	if Parse(seq1.succedent().String()).MainConnective() != uni {
		logger.Print("premise must be universally quantified")
		return false
	}

	val, _, _ := isInstanceOf(seq2.succedent().String(), seq1.succedent().String())
	if !val {
		logger.Print("conclusion not an instance of premise")
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

	return true
}

func uniI(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Universal Elimination depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

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
		if isFormulaSet(d.String()) {
			continue
		}
		if Parse(d).hasTerm(term) {
			logger.Print(term, " cannot appear in datum")
			return false
		}
	}

	return true

}
