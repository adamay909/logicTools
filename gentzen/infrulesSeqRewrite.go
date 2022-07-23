package gentzen

//seq1 is concluding sequent
func seqRewrite(have, want sequent, n int) bool {

	if isSeqReduce(have, want) {
		return true
	}

	if isSeqReorder(have, want) {
		return true
	}

	if isSeqAddition(have, want) {
		return true
	}

	logger.Print("not a rewrite of line ", n)
	return false

}

//check if datum of seq1 is a reduction of seq2
func isSeqReduce(have, want sequent) bool {

	if have.succedent() != want.succedent() {
		return false
	}

	datum1 := have.datumSlice()
	datum2 := want.datumSlice()

	if !(len(datum1) < len(datum2)) {
		return false
	}

	for _, e := range datum2 {
		if !slicesContains(datum1, e) {
			return false
		}
	}

	return true
}

//check if datum of seq1 is a reordering of seq2
func isSeqReorder(have, want sequent) bool {

	if have.succedent() != want.succedent() {
		return false
	}

	datum1 := have.datumSlice()
	datum2 := want.datumSlice()

	if len(datum1) != len(datum2) {
		return false
	}

	datumSort(datum1)
	datumSort(datum2)

	for i := range datum1 {
		if datum1[i] != datum2[i] {
			return false
		}
	}

	return true
}

//check if datum of have is addtion to want
func isSeqAddition(have, want sequent) bool {

	if have.succedent() != have.succedent() {
		return false
	}

	datum1 := have.datumSlice()
	datum2 := want.datumSlice()

	if !(len(datum1) > len(datum2)) {
		return false
	}

	for _, e := range datum2 {

		if !slicesContains(datum1, e) {
			return false
		}
	}

	return true
}
