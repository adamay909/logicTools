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

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())
	n3 := Parse(seq3.succedent().String())

	if n3.MainConnective() != neg {
		logger.Print("conclusion must be negation")
		return false
	}
	/*
		if n1.Class() > n2.Class() {
			if n1.Child1Must().Formula() != n2.Formula() {
				logger.Print("succedents of premises must be negations of each other")
				return false
			}
		}

		if n2.Class() >= n1.Class()+1 {
			if n2.Child1Must().Formula() != n1.Formula() {
				logger.Print("succedents of premises must be negations of each other")
				return false
			}
		}

		if !datumIncludes(seq1.datumSlice(), datum(n3.Child1Must().Formula())) {
			logger.Print("conclusion's negation must be in datums of both premises")
			return false
		}
	*/

	if lneg+n1.Formula() != n2.Formula() && n1.Formula() != lneg+n2.Formula() {
		logger.Print("succedents of premises must be negations of each other")
		return false
	}

	if !datumIncludes(seq2.datumSlice(), datum(n3.Child1Must().Formula())) {
		logger.Print("conclusion's negation must be in datums of both premises")
		return false
	}
	if !datumIncludes(seq1.datumSlice(), datum(n3.Child1Must().Formula())) {
		logger.Print("conclusion's negation must be in datums of both premises")
		return false
	}

	return true
}
