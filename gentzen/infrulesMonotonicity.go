package gentzen

func monotonicity(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Monotonicity depends on one line")
		return false
	}

	ini := d.supportingLines[0].line.seq
	res := d.line.seq
	n := d.line.lines[0]

	Debug("checking datum 1: ", ini.datumSlice(), " against: ", res.datumSlice())

	if isSeqAddition(ini, res) {
		return true
	}

	logger.Print("Monotonicity requires adding item to datum", n)
	return false

}
