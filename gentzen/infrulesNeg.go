package gentzen

func negE(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Negation Elimination depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	n1 := Parse(seq1.succedent().String(), !allowGreekUpper)
	n2 := Parse(seq2.succedent().String(), !allowGreekUpper)

	if n1.MainConnective() != Neg {
		logger.Print("premise must be double negation")
		return false
	}

	if n1.subnode1.MainConnective() != Neg {
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

func negI(d *derivNode) bool {

	if len(d.supportingLines) != 2 {
		logger.Print("Negation Introduction depends on two lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.supportingLines[1].line.seq
	seq3 := d.line.seq

	n1 := Parse(seq1.succedent().String(), !allowGreekUpper)
	n2 := Parse(seq2.succedent().String(), !allowGreekUpper)
	n3 := Parse(seq3.succedent().String(), !allowGreekUpper)

	if n3.MainConnective() != Neg {
		logger.Print("conclusion must be negation")
		return false
	}

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

	wantDatum1 := datumRm(seq1.datumSlice(), datum(n3.Child1Must().Formula()))
	wantDatum2 := datumRm(seq2.datumSlice(), datum(n3.Child1Must().Formula()))

	wantDatum := datumUnion(wantDatum1, wantDatum2)

	if strictCheck {
		if !datumsEqual(wantDatum, seq3.datumSlice()) {
			logger.Print("check datum of conclusion")
			return false
		}
	} else {
		if !datumsEquiv(wantDatum, seq3.datumSlice()) {
			logger.Print("check datum of conclusion")
			return false
		}
	}

	return true
}
