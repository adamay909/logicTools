package gentzen

//seq1 is concluding sequent
func seqRewrite(seq1, seq2 sequent, n int) bool {

	if isSeqReduce(seq2, seq1) {
		return true
	}

	if isSeqReorder(seq2, seq1) {
		return true
	}

	if isSeqAddition(seq2, seq1) {
		return true
	}

	logger.Print("not a rewrite of line ", n)
	return false

}

//check if datum of seq1 is a reduction of seq2
func isSeqReduce(seq1, seq2 sequent) bool {

	if seq1.s != seq2.s {
		return false
	}

	datum1 := seq1.datumSlice()
	datum2 := seq2.datumSlice()

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
func isSeqReorder(seq1, seq2 sequent) bool {

	if seq1.s != seq2.s {
		return false
	}

	datum1 := seq1.datumSlice()
	datum2 := seq2.datumSlice()

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

//check if datum of seq2 is addtion to seq1
func isSeqAddition(seq1, seq2 sequent) bool {

	if seq1.s != seq2.s {
		return false
	}

	datum1 := seq1.datumSlice()
	datum2 := seq2.datumSlice()

	if !(len(datum2) > len(datum1)) {
		return false
	}

	for _, e := range datum1 {

		if !slicesContains(datum2, e) {
			return false
		}
	}

	return true
}
