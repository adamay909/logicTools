package gentzen

var emptySet = "\u2300"

// seq1 is concluding sequent
func seqRewrite(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Sequent rewrite depends on one line")
		return false
	}

	ini := d.supportingLines[0].line.seq
	res := d.line.seq
	n := d.line.lines[0]

	if isSeqReduce(ini, res) {
		return true
	}

	if isSeqReorder(ini, res) {
		return true
	}

	if isSeqAddition(ini, res) {
		return true
	}

	logger.Print("not a rewrite of line ", n)
	return false

}

// check if datum of res is a reduction of ini
func isSeqReduce(ini, res sequent) bool {

	if ini.succedent() != res.succedent() {
		return false
	}

	datum1 := ini.datumSlice()
	datum2 := res.datumSlice()

	if !(len(datum2) < len(datum1)) {
		return false
	}

	for _, e := range datum1 {
		if string(e) == emptySet {
			continue
		}
		if !slicesContains(datum2, e) {
			return false
		}
	}
	return true
}

// check if datum of seq1 is a reordering of seq2
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

// check if datum of have is addtion to want
func isSeqAddition(ini, res sequent) bool {

	if ini.succedent() != res.succedent() {
		return false
	}

	datum1 := ini.datumSlice()
	datum2 := res.datumSlice()

	if !(len(datum2) > len(datum1)) {
		return false
	}

	for _, e := range datum1 {
		if string(e) == emptySet {
			continue
		}
		if !slicesContains(datum2, e) {
			return false
		}
	}

	return true
}
