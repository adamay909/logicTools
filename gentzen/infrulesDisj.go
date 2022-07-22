package gentzen

func disjI(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())

	if n2.MainConnective() != disj {
		logger.Print("conclusion must be a disjunction")
		return false
	}

	if n2.Child1Must().Formula() != n1.Formula() && n2.Child2Must().Formula() != n1.Formula() {
		logger.Print("premise is not one of disjuncts")
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

func disjE(seq1, seq2, seq3, seq4 sequent) bool {

	v1, msg1 := disjEhelper1(seq1, seq2, seq3, seq4)
	v2, msg2 := disjEhelper1(seq2, seq1, seq3, seq4)
	v3, msg3 := disjEhelper1(seq3, seq1, seq2, seq4)

	if !v1 && !v2 && !v3 {
		if msg1 == "must have disjunction among premises" {

			if msg2 == msg1 {
				logger.Print(msg3)
				return false
			}

			logger.Print(msg2)
			return false
		}
		logger.Print(msg1)
		return false
	}

	var v bool
	switch {

	case v1:
		v, msg2 = disjEhelper2(seq1, seq2, seq3, seq4)

	case v2:
		v, msg2 = disjEhelper2(seq2, seq1, seq3, seq4)

	case v3:
		v, msg2 = disjEhelper2(seq3, seq1, seq2, seq4)

	}

	if !v {
		logger.Print(msg2)
		return false
	}

	return true
}

func disjEhelper1(seq1, seq2, seq3, seq4 sequent) (v bool, msg string) {

	v = false

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())
	n3 := Parse(seq3.succedent().String())
	n4 := Parse(seq4.succedent().String())

	if n1.MainConnective() != disj {
		msg = "must have disjunction among premises"
		return
	}

	if n2.Formula() != n3.Formula() {
		msg = "must have two identical succedents"
		return
	}

	if n3.Formula() != n4.Formula() {
		msg = "conclusion must be identical to two of the succedents"
		return
	}

	d1 := datum(n1.Child1Must().Formula())
	d2 := datum(n1.Child2Must().Formula())
	/*
		if !datumIncludes(seq2.datumSlice(), d1) && !datumIncludes(seq3.datumSlice(), d1) {
			msg = "one of disjuncts not in datums"
			return
		}

		if !datumIncludes(seq2.datumSlice(), d2) && !datumIncludes(seq3.datumSlice(), d2) {
			msg = "one of disjuncts not in datums"
			return
		}
	*/
	if !(datumIncludes(seq2.datumSlice(), d1) && datumIncludes(seq3.datumSlice(), d2)) && !(datumIncludes(seq2.datumSlice(), d2) && datumIncludes(seq3.datumSlice(), d1)) {
		msg = "one of disjuncts not in datums"
		return
	}

	v = true
	return
}

func disjEhelper2(seq1, seq2, seq3, seq4 sequent) (v bool, msg string) {

	v = false

	n1 := Parse(seq1.succedent().String())
	d1 := datum(n1.Child1Must().Formula())
	d2 := datum(n1.Child2Must().Formula())

	datumCanonical := datumRm(datumUnion(seq1.datumSlice(), seq2.datumSlice(), seq3.datumSlice()), d1, d2)

	if strictCheck {
		if !datumsEqual(datumCanonical, seq4.datumSlice()) {
			msg = "datum of conclusion incorrect"
			return
		}
	} else {
		if !datumsEquiv(datumCanonical, seq4.datumSlice()) {
			msg = "datum of conclusion incorrect"
			return
		}
	}
	v = true

	return
}
