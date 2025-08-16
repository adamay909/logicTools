package gentzen

func conjE(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Conjunction Elimination depends on a single line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	n1 := Parse(seq1.succedent().String(), !allowGreekUpper)
	n2 := Parse(seq2.succedent().String(), !allowGreekUpper)

	if n1.Val.connective != Conj {
		logger.Print("must start with conjunction")
		return false
	}

	if n2.String() != n1.Child(0).String() && n2.String() != n1.Child(1).String() {
		logger.Print("conclusion not one of conjuncts")
		return false
	}

	if !datumsEquiv(seq1.datumSlice(), seq2.datumSlice()) {
		logger.Print("datum of conclusion must be same as datum of premise")
		return false
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

	n1 := Parse(seq1.succedent().String(), !allowGreekUpper)
	n2 := Parse(seq2.succedent().String(), !allowGreekUpper)
	n3 := Parse(seq3.succedent().String(), !allowGreekUpper)

	if n3.Val.connective != Conj {
		logger.Print("conclusion must be a conjunction")
		return false
	}

	if n1.String() != n3.Child(0).String() && n1.String() != n3.Child(1).String() {
		logger.Print("succedent of conclusion must be conjunction of succedents of premises")

		return false
	}

	if n2.String() != n3.Child(0).String() && n2.String() != n3.Child(1).String() {
		logger.Print("succedent of conclusion must be conjunction of succedents of premises")
		return false
	}

	datumCanonical := datumUnion(seq1.datumSlice(), seq2.datumSlice())
	if !datumsEquiv(datumCanonical, seq3.datumSlice()) {
		logger.Print("datum of conclusion must be union of datums of premises")
		return false
	}

	return true
}
