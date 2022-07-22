package gentzen

func assumption(seq sequent) bool {

	if seq.datum().String() != seq.succedent().String() {
		logger.Print("datum and subseqent cannot differ for assumption")
		return false
	}
	return true

}
