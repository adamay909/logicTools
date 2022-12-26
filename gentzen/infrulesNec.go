package gentzen

func necE(seq1, seq2 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}

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
func necI(seq1, seq2 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
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

	if len(seq1.datumSlice()) > 0 {
		logger.Print("must start with theorem")
		return false
	}

	if len(seq2.datumSlice()) > 0 {
		logger.Print("datum cannot change")
		return false
	}

	return true

}

func necI_T(seq1, seq2 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
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

	if Parse(seq1.succedent().String()).IsModal() {
		logger.Print("cannot necessitate modal formulas")
		return false
	}

	for _, d := range seq1.datumSlice() {
		if len(d) == 0 {
			continue
		}
		if isFormulaSet(d.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
		if Parse(d.String()).MainConnective() != nec {
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

func necI_S4(seq1, seq2 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
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

	for _, d := range seq1.datumSlice() {
		if len(d) == 0 {
			continue
		}
		if isFormulaSet(d.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
		if Parse(d.String()).MainConnective() != nec {
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
func necI_S5(seq1, seq2 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
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

	for _, d := range seq1.datumSlice() {
		if len(d) == 0 {
			continue
		}
		if isFormulaSet(d.String()) {
			logger.Print("all datum items must be modal claims")
			return false
		}
		if !isModalClaim(d.String()) {
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
