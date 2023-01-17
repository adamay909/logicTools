package gentzen

func condE(d *derivNode) bool {

	if len(d.supportingLines) != 2 {
		logger.Print("Conditional Elimination depends on two lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.supportingLines[1].line.seq
	seq3 := d.line.seq

	v1, msg1 := condEhelper(seq1, seq2, seq3)

	v2, msg2 := condEhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "no conditional found" {

		logger.Print(msg2)
		return false

	}

	logger.Print(msg1)
	return false

}

func condEhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())
	n3 := Parse(seq3.succedent().String())
	v = false

	if n1.MainConnective() != cond {
		msg = "no conditional found"
		return
	}

	if n2.Formula() != n1.Child1Must().Formula() {
		msg = "mismatch between conditional and other premise"
		return
	}

	d3 := datumUnion(seq1.datumSlice(), seq2.datumSlice())
	canonicalSeq := mkSequent(d3, n1.Child2Must())
	//canonicalSeq := sequent{datum(d3.String()), plshFormula(n1.Child2Must().Formula())}
	if canonicalSeq.succedent().String() != n3.Formula() {
		msg = "conclusion does not match consequent of conditional"
		return
	}

	if strictCheck {
		if !equalSequents(canonicalSeq, seq3) {
			msg = "datum of conclusion must be union of datums of premises"
			return
		}
	} else {
		if !equivSequents(canonicalSeq, seq3) {
			msg = "datum of conclusion must be union of datums of premises"
			return
		}
	}
	v = true

	return

}

func condI(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Conditional Introduction depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())

	if n2.MainConnective() != cond {
		logger.Print("main connective of conclusion must be conditional")
		return false
	}

	if n2.Child2Must().Formula() != n1.Formula() {
		logger.Print("consequent of conclusion must be succedent of premise")
		return false
	}

	if !datumIncludes(seq1.datumSlice(), datum(n2.Child1Must().Formula())) {
		logger.Print("antecedent of conditional must be in datum of premise")
		return false
	}

	d1 := datumAdd(seq2.datumSlice(), datum(n2.Child1Must().Formula()))
	canonicalPrem := mkSequent(d1, n2.Child2Must())
	//canonicalPrem := sequent{datum(d1.String()), plshFormula(n2.Child2Must().Formula())}

	if strictCheck {
		if !equalSequents(canonicalPrem, seq1) {
			logger.Print("must remove exactly one datum item")
			return false
		}
	} else {
		if !equivSequents(canonicalPrem, seq1) {
			logger.Print("must remove exactly one datum item")
			return false
		}
	}
	return true
}
