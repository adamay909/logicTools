package gentzen

func posI(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}
	if len(d.supportingLines) != 1 {
		logger.Print("Possibility Introduction depends on two lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	if Parse(seq2.succedent().String()).MainConnective() != pos {
		logger.Print("conclusion must be possibility")
		return false
	}

	if strictCheck {
		if !datumsEqual(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !datumsEquiv(seq1.datumSlice(), seq2.datumSlice()) {
			logger.Print("datum cannot change")
			return false
		}
	}

	if isModalInstanceOf(Parse(seq2.succedent().String()).Child1Must().String(), seq1.succedent().String()) {
		logger.Print("conclusion does not match premise")
		return false
	}

	return true

}

func posE(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}
	if len(d.supportingLines) != 2 {
		logger.Print("Possibility Elimination depends on two lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.supportingLines[1].line.seq
	seq3 := d.line.seq

	v1, msg1 := posEhelper(seq1, seq2, seq3)

	v2, msg2 := posEhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "no possibility in premises" {
		logger.Print(msg2)
		return false
	}
	logger.Print(msg1)
	return false

}

func posEhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	v = false

	if Parse(seq1.succedent().String()).MainConnective() != pos {
		msg = "no possibility in premises"
		return
	}

	if Parse(seq3.succedent().String()).MainConnective() != pos {
		msg = "conclusison must be possibility claim"
		return
	}

	if seq2.succedent().String() != Parse(seq3.succedent()).Child1Must().Formula() {
		msg = "conclusion does not match premises"
		return
	}

	found := false

	datum2 := seq2.datumSlice()

	for _, d := range datum2 {
		if isFormulaSet(d.String()) {
			continue
		}
		found = isModalInstanceOf(d.String(), seq1.succedent().String())
		if found {
			datum2 = datumRm(datum2, d)
			break
		}
	}

	if !found {
		msg = "no datum item found as instance of modal claim"
		return
	}
	/**
		if !isModalClaim(seq3.s.String()) {
			msg = "target conclusion must be a modal claim"
			return
		}
	**/

	datum1 := seq1.datumSlice()
	/**	for _, d := range datum1 {
			if len(d) == 0 {
				continue
			}

			if Parse(d.String()).Formula() == Parse(seq1.succedent()).Formula() {
				continue
			}

			if isFormulaSet(d.String()) {
				msg = "all datum items must be modal claims"
				return
			}
			if Parse(d.String()).MainConnective() != nec {
				msg = "all datum items must be necessity claims"
				return
			}
		}
	**/
	for _, d := range datum2 {
		if len(d) == 0 {
			continue
		}
		if isFormulaSet(d.String()) {
			msg = "proviso about datum not respected"
			return
		}
		if Parse(d.String()).MainConnective() != nec {
			msg = "proviso about datum not respected"
			return
		}
	}

	if strictCheck {
		if !datumsEqual(datumUnion(datum1, datum2), seq3.datumSlice()) {
			msg = "datum of conclusion must be union of datums of premise"
			v = false
			return
		}
	} else {
		if !datumsEquiv(datumUnion(datum1, datum2), seq3.datumSlice()) {
			msg = "datum of conclusion must be union of datums of premise"
			v = false
			return
		}
	}
	v = true
	msg = ""
	return

}
func posE_S5(d *derivNode) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}

	if len(d.supportingLines) != 2 {
		logger.Print("S5 Possibility Elimination depends on two lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.supportingLines[1].line.seq
	seq3 := d.line.seq

	v1, msg1 := posEhelper_S5(seq1, seq2, seq3)

	v2, msg2 := posEhelper_S5(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "no possibility in premises" {
		logger.Print(msg2)
		return false
	}
	logger.Print(msg1)
	return false

}

func posEhelper_S5(seq1, seq2, seq3 sequent) (v bool, msg string) {

	v = false

	if Parse(seq1.succedent().String()).MainConnective() != pos {
		msg = "no possibility in premises"
		return
	}

	if seq2.succedent().String() != seq3.succedent().String() {
		msg = "conclusion does not match premises"
		return
	}

	if !isModalClaim(seq2.succedent().String()) {
		msg = "proviso not met"
		return
	}

	found := false

	datum2 := seq2.datumSlice()

	for _, d := range datum2 {
		if isFormulaSet(d.String()) {
			continue
		}
		found = isModalInstanceOf(d.String(), seq1.succedent().String())
		if found {
			datum2 = datumRm(datum2, d)
			break
		}
	}

	if !found {
		msg = "no datum item found as instance of modal claim"
		return
	}

	datum1 := seq1.datumSlice()
	for _, d := range datum1 {
		if len(d) == 0 {
			continue
		}

		if Parse(d.String()).Formula() == Parse(seq1.succedent()).Formula() {
			continue
		}

		if isFormulaSet(d.String()) {
			msg = "all datum items must be modal claims"
			return
		}
		if !isModalClaim(d.String()) {
			msg = "all datum items must be modal claims"
			return
		}
	}

	for _, d := range datum2 {
		if len(d) == 0 {
			continue
		}
		if isFormulaSet(d.String()) {
			msg = "all datum items must be modal claims"
			return
		}
		if !isModalClaim(d.String()) {
			msg = "all datum items must be modal claims"
			return
		}
	}

	if strictCheck {
		if !datumsEqual(datumUnion(datum1, datum2), seq3.datumSlice()) {
			msg = "datum of conclusion must be union of datums of premise"
			v = false
			return
		}
	} else {
		if !datumsEquiv(datumUnion(datum1, datum2), seq3.datumSlice()) {
			msg = "datum of conclusion must be union of datums of premise"
			v = false
			return
		}
	}
	v = true
	msg = ""
	return

}
