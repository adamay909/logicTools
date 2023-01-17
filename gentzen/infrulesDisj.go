package gentzen

import (
	"errors"
)

func disjI(d *derivNode) bool {

	if len(d.supportingLines) != 1 {
		logger.Print("Disjunction Introduction depends on one line")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.line.seq

	n1 := Parse(seq1.succedent().String())
	n2 := Parse(seq2.succedent().String())

	if n2.MainConnective() != disj {
		logger.Print("conclusion must be a disjunction")
		return false
	}

	if n2.Child1Must().Formula() != n1.Formula() && n2.Child2Must().Formula() != n1.Formula() {
		logger.Print("premise is not one of disjuncts")
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

	return true
}

func disjE(d *derivNode) bool {

	if len(d.supportingLines) != 3 {
		logger.Print("Disjunction Introduction depends on three lines")
		return false
	}

	seq1 := d.supportingLines[0].line.seq
	seq2 := d.supportingLines[1].line.seq
	seq3 := d.supportingLines[2].line.seq
	seq4 := d.line.seq

	var err error

	seq := []sequent{seq1, seq2, seq3}

	//check if there are premises of the right form
	for i := 0; i < 3; i++ {
		seq1, seq2, seq3, err = disjEhelper3(seq, i)
		if err == nil {
			break
		}
	}

	if err != nil {
		logger.Print(err.Error())
		return false
	}

	//check if non-disjunction premises have the right succedent
	err = disjEhelper4(seq2, seq3, seq4)
	if err != nil {
		logger.Print(err.Error())
		return false
	}

	//check datum
	ok := disjEhelper5(seq1, seq2, seq3, seq4)
	if !ok {
		logger.Print("check your datums")
		return false
	}
	return true

}

func disjEhelper3(seq []sequent, i int) (seq1, seq2, seq3 sequent, err error) {

	var j int

	ok := false
	for j = range seq {

		if Parse(seq[j].succedent()).MainConnective() == disj {
			seq1 = seq[j]
			switch j {
			case 0:
				seq2 = seq[1]
				seq3 = seq[2]
			case 1:
				seq2 = seq[0]
				seq3 = seq[2]
			case 2:
				seq2 = seq[0]
				seq3 = seq[1]
			}
			d1 := datum(Parse(seq1.succedent()).Child1Must().Formula())
			d2 := datum(Parse(seq1.succedent()).Child2Must().Formula())
			if !datumIncludes(seq2.datumSlice(), d1) && !datumIncludes(seq3.datumSlice(), d1) {
				continue
			}
			if !datumIncludes(seq2.datumSlice(), d2) && !datumIncludes(seq3.datumSlice(), d2) {
				continue
			}
			ok = true
			break
		}
	}

	if !ok {
		err = errors.New("check forms of your premises")
		return
	}

	return
}

func disjEhelper4(seq2, seq3, seq4 sequent) (err error) {

	want := Parse(seq4.succedent()).Formula()

	have1 := Parse(seq2.succedent()).Formula()

	have2 := Parse(seq3.succedent()).Formula()

	if want != have1 {
		err = errors.New("succedents of premises do not match succedent of conclusion")
		return
	}

	if want != have2 {
		err = errors.New("succedents of premises do not match succedent of conclusion")
		return
	}

	return
}

func disjEhelper5(seq ...sequent) bool {

	datumU := datumUnion(seq[0].datumSlice(), seq[1].datumSlice(), seq[2].datumSlice())

	d1 := Parse(seq[0].succedent()).Child1Must().Formula()
	d2 := Parse(seq[0].succedent()).Child2Must().Formula()

	want := datumRm(datumU, d1, d2)
	have := seq[3].datumSlice()
	if strictCheck {
		return datumsEqual(want, have)
	}
	return datumsEquiv(want, have)
}
