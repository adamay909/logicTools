package gentzen

func assumption(d *derivNode) bool {

	if len(d.supportingLines) > 0 {
		logger.Print("Assumption does not depend on other lines")
		return false
	}

	seq := d.line.seq

	if seq.datum().String() != seq.succedent().String() {
		logger.Print("datum and subseqent cannot differ for assumption")
		return false
	}
	return true

}
