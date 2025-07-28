package gentzen

func scopeReplacement(d *derivNode) bool {

	if oPL {
		logger.Print("Scope Replacement not implemented for Predicate Logic")
		return false
	}

	if !oML {
		logger.Print("Scope Replament only allowed with Modal Logic")
		return false
	}

	if len(d.supportingLines) != 1 {
		logger.Print("Scope Replacement depends on one line")
		return false
	}

	s1 := d.supportingLines[0].line.seq
	s2 := d.line.seq

	m1 := Parse(s1.succedent(), !allowGreekUpper).MainConnective()
	if m1 != Neg && m1 != Nec && m1 != Pos {
		logger.Print("Scope Replacement only works for", lneg, ", ", lnec, ", ", lpos)
		return false
	}

	if m1 != Parse(s2.succedent(), !allowGreekUpper).MainConnective() {
		logger.Print("Main connective cannot change")
		return false
	}

	sc1 := Parse(s1.succedent().String(), !allowGreekUpper).Child1Must().Formula()
	sc2 := Parse(s2.succedent().String(), !allowGreekUpper).Child1Must().Formula()

	if !Parse(sc1, !allowGreekUpper).IsPureSL() {
		logger.Print("formula inside scope must be sentence of pure Sentential Logic.")
		return false
	}

	if !Parse(sc2, !allowGreekUpper).IsPureSL() {
		logger.Print("formula inside scope must be sentence of pure Sentential Logic.")
		return false
	}

	n1 := ">" + sc1 + sc2
	n2 := ">" + sc1 + sc2

	if !IsTautology(n1) {
		logger.Print("Scope Replacement requires logical equivalences")
		return false
	}

	if !IsTautology(n2) {
		logger.Print("Scope Replacement requires logical equivalences")
		return false
	}

	if strictCheck {
		if !datumsEqual(s1.datumSlice(), s2.datumSlice()) {
			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !datumsEquiv(s1.datumSlice(), s2.datumSlice()) {
			logger.Print("datum cannot change")
			return false
		}
	}

	return true

}
