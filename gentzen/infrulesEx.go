package gentzen

func exI(seq1, seq2 sequent) bool {

	if Parse(seq2.succedent().String()).MainConnective() != ex {
		logger.Print("conclusion must be existentially quantified")
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

	val, _, _ := isInstanceOf(seq1.succedent().String(), seq2.succedent().String())
	if !val {
		logger.Print("conclusion must be existential generalization of premise")
	}
	return val
}

func exE(seq1, seq2, seq3 sequent) bool {

	v1, msg1 := exEhelper(seq1, seq2, seq3)

	v2, msg2 := exEhelper(seq2, seq1, seq3)

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

func exEhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	v = false

	if Parse(seq1.succedent().String()).MainConnective() != ex {
		msg = "no existential quantification in premises"
		return
	}

	if seq2.s != seq3.s {
		msg = "conclusion does not match premises"
		return
	}

	found := false
	var kappa string

	datum2 := seq2.datumSlice()

	for _, d := range datum2 {
		if d[:1] == `\` {
			continue
		}
		found, _, kappa = isInstanceOf(d.String(), seq1.succedent().String())
		if found {
			datum2 = datumRm(datum2, d)
			break
		}
	}

	if !found {
		msg = "no datum item found as instance of existential claim"
		return
	}

	datum1 := seq1.datumSlice()
	for _, d := range datum1 {
		if len(d) == 0 {
			continue
		}
		if d[:1] == `\` {
			continue
		}
		if Parse(d).hasTerm(kappa) {
			msg = kappa + " may not appear in any datum items"
			return
		}
	}

	for _, d := range datum2 {
		if len(d) == 0 {
			continue
		}
		if d[:1] == `\` {
			continue
		}
		if Parse(d).hasTerm(kappa) {
			msg = kappa + " may not appear in any datum items"
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
