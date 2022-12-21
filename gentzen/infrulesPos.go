package gentzen

func posI(seq1, seq2 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}
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

func posE(seq1, seq2, seq3 sequent) bool {

	if !oML {
		logger.Print("Modal Logic not allowed")
		return false
	}
	v1, msg1 := posEhelper(seq1, seq2, seq3)

	v2, msg2 := posEhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "no existential generalization in premises" {
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

	if seq2.s != seq3.s {
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

	if !isModalClaim(seq3.s.String()) {
		msg = "target conclusion must be a modal claim"
		return
	}

	datum1 := seq1.datumSlice()
	for _, d := range datum1 {
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
