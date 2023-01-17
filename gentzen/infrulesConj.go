package gentzen

func conjE(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Conjunction Elimination depends on a single line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

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

func conjI(d *derivNode) bool {

	if len(d.supportingLines) != 2 {
		logger.Print("Conjunction Introduction depends on two lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.supportingLines[1].line.seq
	seq3 := d.line.seq

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
