package gentzen

func necE(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}

	if len(d.supportingLines) != 1 {
		logger.Print("Necessity Elimination depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	if Parse(seq1.succedent().String()).MainConnective() != nec {
		logger.Print("premise must be a necessary truth")
		return false
	}

	if !isModalInstanceOf(seq2.succedent().String(), seq1.succedent().String()) {
		logger.Print("conclusion does match premise")
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
func necI(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}

	if len(d.supportingLines) != 1 {
		logger.Print("Necessity Introduction depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	dep := d.supportingLines[0]
	if !dep.isTheorem() {
		logger.Print("Necessity Introduction depends on a theorem")
		return false
	}

	if Parse(seq2.succedent().String()).MainConnective() != nec {
		logger.Print("conclusion must be necessary truth")
		return false
	}
	if !isModalInstanceOf(seq1.succedent().String(), seq2.succedent().String()) {
		logger.Print("conclusion not a necessitation of premise")
		return false
	}

	if len(seq2.datumSlice()) > 0 {
		logger.Print("datum cannot change")
		return false
	}

	return true

}

func necI_T(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}

	if len(d.supportingLines) != 1 {
		logger.Print("T Necessity Introduction depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	if Parse(seq2.succedent().String()).MainConnective() != nec {
		logger.Print("conclusion must be necessary truth")
		return false
	}
	if !isModalInstanceOf(seq1.succedent().String(), seq2.succedent().String()) {
		logger.Print("conclusion not a necessitation of premise")
		return false
	}

	if Parse(seq1.succedent().String()).IsModal() {
		logger.Print("cannot necessitate modal formulas")
		return false
	}

	for _, datum := range seq1.datumSlice() {
		if len(datum) == 0 {
			continue
		}
		if isFormulaSet(datum.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
		if Parse(datum.String()).MainConnective() != nec {
			logger.Print("all datum items must be necessity claims")
			return false
		}
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

func necI_S4(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}
	if len(d.supportingLines) != 1 {
		logger.Print("S4 Necessity Introduction depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	if Parse(seq2.succedent().String()).MainConnective() != nec {
		logger.Print("conclusion must be necessary truth")
		return false
	}
	if !isModalInstanceOf(seq1.succedent().String(), seq2.succedent().String()) {
		logger.Print("conclusion not a necessitation of premise")
		return false
	}

	for _, datum := range seq1.datumSlice() {
		if len(datum) == 0 {
			continue
		}
		if isFormulaSet(datum.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
		if Parse(datum.String()).MainConnective() != nec {
			logger.Print("all datum items must be necessity claims")
			return false
		}
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
func necI_S5(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}
	if len(d.supportingLines) != 1 {
		logger.Print("S5 Necessity Introduction depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	if Parse(seq2.succedent().String()).MainConnective() != nec {
		logger.Print("conclusion must be necessary truth")
		return false
	}
	if !isModalInstanceOf(seq1.succedent().String(), seq2.succedent().String()) {
		logger.Print("conclusion not a necessitation of premise")
		return false
	}

	for _, datum := range seq1.datumSlice() {
		if len(datum) == 0 {
			continue
		}
		if isFormulaSet(datum.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
		if !isModalClaim(datum.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
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

func isModalClaim(s string) bool {

	n := Parse(s)

	if n.MainConnective().isModalOperator() {
		return true
	}

	if n.MainConnective().isNegation() {
		return n.Child1Must().MainConnective().isModalOperator()
	}

	return false
}
