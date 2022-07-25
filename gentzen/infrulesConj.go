package gentzen

func conjE(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())

	if n1.MainConnective() != conj {
		logger.Print("must start with conjunction")
		return false
	}

	if n2.Formula() != n1.subnode1.Formula() {
		if n2.Formula() != n1.subnode2.Formula() {
			logger.Print("conclusion not one of conjuncts")
			return false
		}
	}
	if strictCheck {
		if !datumsEqual(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum of conclusion must be same as datum of premise")
			return false
		}
	} else {
		if !datumsEquiv(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum of conclusion must be same as datum of premise")
			return false
		}
	}
	return true
}

func conjI(seq1, seq2, seq3 sequent) bool {

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())
	n3 := Parse(seq3.succedent().String())

	if n3.MainConnective() != conj {
		logger.Print("conclusion must be a conjunction")
		return false
	}

	if n1.Formula() != n3.subnode1.Formula() && n1.Formula() != n3.subnode2.Formula() {
		logger.Print("succedent of conclusion must be conjunction of succedents of premises")

		return false
	}

	if n2.Formula() != n3.subnode1.Formula() && n2.Formula() != n3.subnode2.Formula() {
		logger.Print("succedent of conclusion must be conjunction of succedents of premises")
		return false
	}

	datumCanonical := datumUnion(seq1.datumSlice(), seq2.datumSlice())
	if strictCheck {
		if !datumsEqual(datumCanonical, seq3.datumSlice()) {
			logger.Print("datum of conclusion must be union of datums of premises")
			return false
		}
	} else {
		if !datumsEquiv(datumCanonical, seq3.datumSlice()) {
			logger.Print("datum of conclusion must be union of datums of premises")
			return false
		}
	}

	return true
}
