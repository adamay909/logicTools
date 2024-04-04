package gentzen

import (
	"strings"
)

func SeqReductionString(s string) string {

	if !isReducibleSeq(s) {
		return s
	}

	parts := strings.Split(s, ";")
	seq := mkSequent(parts[0], parts[1])

	rseries := sequentReductionSeries(seq)

	out := s + "\n"

	for _, row := range rseries {

		for _, e := range row {

			out = out + e.String() + ", "

		}
		out = strings.TrimRight(out, ", ") + "\n"
	}

	return out
}

func SeqReductionLatexString(s string) string {

	if !isReducibleSeq(s) {
		return s
	}

	parts := strings.Split(s, ";")
	seq := mkSequent(parts[0], parts[1])

	rseries := sequentReductionSeries(seq)

	out := seq.StringLatex() + "\n\n"

	for _, row := range rseries {

		for _, e := range row {

			out = out + e.StringLatex() + ";~~ "

		}
		out = strings.TrimRight(out, ";~~ ") + "\n\n"
	}

	return out
}

func IsProvableSequent(s string) bool {

	if !isReducibleSeq(s) {
		return false
	}

	parts := strings.Split(s, ";")
	seq := mkSequent(parts[0], parts[1])

	rseries := sequentReductionSeries(seq)

	rset := rseries[len(rseries)-1]

	for _, seq := range rset {

		if len(seq.datumSlice()) != 1 {
			return false
		}

		if Parse(string(seq.datumSlice()[0])).String() != Parse(seq.succedent()).String() {

			return false

		}
	}
	return true
}

func sequentReductionSeries(s sequent) [][]sequent {

	var out [][]sequent

	var rset []sequent

	rset = append(rset, s)
	changed := true
	for changed {

		rset, changed = sequentReduction(rset)
		if changed {
			out = append(out, rset)
		}
	}

	return out
}

func sequentReduction(seqSet []sequent) (rset []sequent, change bool) {

	redRule := []func(sequent) ([]sequent, bool){
		redr1,  //chech if base sequent
		redr2,  //back formula among front formulas
		redr3,  //contradicting front formulas
		redr10, //all basic sentences
		redr11, //back formula Cond
		redr12, //back formula Disj
		redr7,  //back formula DN
		redr15, //back formula negated Cond
		redr16, //back formula negated Disj
		redr4,  //front formula has DN
		redr5,  //front formulas has Conj
		redr17, //front formula negated Cond
		redr18, //front formula negated Disj
		redr8,  //back formula Conj
		redr9,  //back formula Neg
		redr13, //front formula Cond
		redr14, //front formula Disj
		redr6,  //front formulas has negated Conj
	}
	var applied bool
	var acount int

	for _, s := range seqSet {

		var rs []sequent

		for rn, rule := range redRule {
			rs, applied = rule(s)

			if applied {
				rset = append(rset, rs...)
				acount++
				if rn == 0 || rn == 3 {
					acount--
				}
				break
			}
		}
	}

	change = acount > 0
	nrset := slicesCleanDuplicates(rset)
	rset = nil
	rset = append(rset, nrset...)

	return

}

// remove duplications in datum
func datumClean(s sequent) (o sequent) {

	old := s.datumSlice()
	newd := slicesCleanDuplicates(old)

	if len(old) == len(newd) {
		o = s
		return
	}

	o = mkSequent(datumSlice(newd), s.succedent())
	return
}

// if base sequent
func redr1(s sequent) (o []sequent, applied bool) {

	applied = false
	if len(s.datumSlice()) != 1 {
		o = append(o, datumClean(s))
		return
	}

	if Parse(string(s.datumSlice()[0])).String() != s.succedent().String() {
		o = append(o, datumClean(s))
		return
	}

	o = append(o, datumClean(s))

	applied = true

	return
}

// If back formula in front formulas
func redr2(s sequent) (o []sequent, applied bool) {

	suc := s.succedent()

	if !slicesContains(s.datumSlice(), datum(suc.String())) {
		applied = false
		o = append(o, datumClean(s))
		return
	}

	o = append(o, datumClean(mkSequent(suc.String(), suc.String())))
	applied = true

	return
}

// If all basic sentences
func redr10(s sequent) (o []sequent, applied bool) {

	applied = false

	for _, d := range s.datumSlice() {
		if !Parse(d.String()).IsBasic() {
			o = append(o, datumClean(s))
			return
		}
	}

	if !Parse(s.succedent()).IsBasic() {
		o = append(o, datumClean(s))
		return
	}

	applied = true
	o = append(o, datumClean(s))
	return
}

// If front formulas contains contradicting sentences
func redr3(s sequent) (o []sequent, applied bool) {

	applied = false

	for _, i := range s.datumSlice() {

		ns := Negate(Parse(i.String()))

		for _, j := range s.datumSlice() {

			if ns.String() == j.String() {
				o = nil
				o = append(o, datumClean(mkSequent(i.String(), i.String())))
				o = append(o, datumClean(mkSequent(ns.String(), ns.String())))
				applied = true
				return
			}
		}
	}
	o = append(o, datumClean(s))
	return
}

// If front formulas contains a double negation
func redr4(s sequent) (o []sequent, applied bool) {

	applied = false

	var dslice datumSlice

	dslice = append(dslice, s.datumSlice()...)
	nd := ""

	for _, d := range dslice {

		if !Parse(d.String()).IsDoubleNegation() {
			nd = nd + string(d) + ","
			continue
		}
		applied = true
		nd = nd + Parse(d.String()).Child1Must().Child1Must().String() + ","

	}

	nd = strings.TrimRight(nd, ",")
	o = append(o, datumClean(mkSequent(nd, s.succedent())))

	return
}

// If front formulas contains a conjunction
func redr5(s sequent) (o []sequent, applied bool) {

	var newdatumSlice, repls datumSlice
	applied = false

	for k, i := range s.datumSlice() {
		n := Parse(i.String())
		if n.IsConjunction() {
			d1 := n.Child1Must()
			d2 := n.Child2Must()
			repls = append(repls, datum(d1.String()))
			repls = append(repls, datum(d2.String()))
			newdatumSlice = slicesReplace(s.datumSlice(), k, repls)
			seq1 := mkSequent(newdatumSlice, s.succedent())
			o = append(o, datumClean(seq1))
			applied = true
			return
		}
	}
	o = append(o, datumClean(s))
	return
}

//If front formulas contain a negated conjunction

func redr6(s sequent) (o []sequent, applied bool) {

	applied = false

	var newds1, newds2 datumSlice

	for k, i := range s.datumSlice() {

		n := Parse(i.String())
		if n.IsNegation() {
			if n.Child1Must().IsConjunction() {
				d1 := Negate(n.Child1Must().Child1Must())
				d2 := Negate(n.Child1Must().Child2Must())

				newds1 = slicesReplace(s.datumSlice(), k, []datum{datum(d1.String())})
				newds2 = slicesReplace(s.datumSlice(), k, []datum{datum(d2.String())})
				seq1 := mkSequent(newds1, s.succedent())
				seq2 := mkSequent(newds2, s.succedent())
				o = append(o, datumClean(seq1), datumClean(seq2))
				applied = true
				return
			}
		}
	}
	o = append(o, s)
	return
}

//If back formula is double negation

func redr7(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsDoubleNegation() {

		o = append(o, datumClean(mkSequent(s.datum(), Parse(s.succedent()).Child1Must().Child1Must())))
		applied = true
		return
	}

	o = append(o, datumClean(s))
	return
}

//If back formula is a conjunction

func redr8(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsConjunction() {
		seq1 := mkSequent(s.datum(), Parse(s.succedent()).Child1Must())
		seq2 := mkSequent(s.datum(), Parse(s.succedent()).Child2Must())
		o = append(o, datumClean(seq1), datumClean(seq2))
		applied = true
		return
	}

	o = append(o, datumClean(s))
	return
}

// If back formula is negation
func redr9(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsNegation() {

		s2 := Parse(s.succedent()).AtomicSentences()[0]
		nd1 := append(s.datumSlice(), datum(Parse(s.succedent()).Child1Must().String()))
		seq1 := mkSequent(nd1, Parse(s2))
		seq2 := mkSequent(nd1, Negate(Parse(s2)))

		o = append(o, datumClean(seq1), datumClean(seq2))
		applied = true
		return
	}

	o = append(o, datumClean(s))
	return
}

// If back formula is conditional
func redr11(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsConditional() {

		d1 := Parse(s.succedent()).Child1Must()
		d2 := Parse(s.succedent()).Child2Must()

		ns := Negate(Conjoin(d1, Negate(d2)))

		o = append(o, datumClean(mkSequent(s.datumSlice(), ns)))
		applied = true
		return
	}

	o = append(o, datumClean(s))
	return
}

// If back formula is disjunction
func redr12(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsDisjunction() {

		d1 := Parse(s.succedent()).Child1Must()
		d2 := Parse(s.succedent()).Child2Must()

		ns := Negate(Conjoin(Negate(d1), Negate(d2)))

		o = append(o, datumClean(mkSequent(s.datumSlice(), ns)))
		applied = true
		return
	}

	o = append(o, datumClean(s))
	return
}

// If front formulas contains a conditional
func redr13(s sequent) (o []sequent, applied bool) {

	applied = false
	var newds1, newds2 datumSlice

	for k, i := range s.datumSlice() {

		n := Parse(i.String())
		if n.IsConditional() {
			d1 := Negate(n.Child1Must())
			d2 := n.Child2Must()

			newds1 = slicesReplace(s.datumSlice(), k, []datum{datum(d1.String())})
			newds2 = slicesReplace(s.datumSlice(), k, []datum{datum(d2.String())})
			seq1 := mkSequent(newds1, s.succedent())
			seq2 := mkSequent(newds2, s.succedent())
			o = append(o, datumClean(seq1), datumClean(seq2))
			applied = true
			return
		}
	}
	o = append(o, s)
	return
}

// If front formulas contains a disjunction
func redr14(s sequent) (o []sequent, applied bool) {

	applied = false
	var newds1, newds2 datumSlice

	for k, i := range s.datumSlice() {

		n := Parse(i.String())
		if n.IsDisjunction() {
			d1 := n.Child1Must()
			d2 := n.Child2Must()

			newds1 = slicesReplace(s.datumSlice(), k, []datum{datum(d1.String())})
			newds2 = slicesReplace(s.datumSlice(), k, []datum{datum(d2.String())})
			seq1 := mkSequent(newds1, s.succedent())
			seq2 := mkSequent(newds2, s.succedent())
			o = append(o, datumClean(seq1), datumClean(seq2))
			applied = true
			return
		}
	}
	o = append(o, s)
	return
}

// If back formula is a negated conditional
func redr15(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsNegation() {

		if Parse(s.succedent()).Child1Must().IsConditional() {
			applied = true
			d1 := Parse(s.succedent()).Child1Must().Child1Must()
			d2 := Parse(s.succedent()).Child1Must().Child2Must()

			ns := Conjoin(d1, Negate(d2))

			o = append(o, datumClean(mkSequent(s.datumSlice(), ns)))
			return
		}
	}
	o = append(o, datumClean(s))
	return
}

// If back formula is a negated disjunction
func redr16(s sequent) (o []sequent, applied bool) {

	applied = false

	if Parse(s.succedent()).IsNegation() {

		if Parse(s.succedent()).Child1Must().IsDisjunction() {
			applied = true
			d1 := Parse(s.succedent()).Child1Must().Child1Must()
			d2 := Parse(s.succedent()).Child1Must().Child2Must()

			ns := Conjoin(Negate(d1), Negate(d2))

			o = append(o, datumClean(mkSequent(s.datumSlice(), ns)))
			return
		}
	}
	o = append(o, datumClean(s))
	return
}

// If front formula has a negated conditional
func redr17(s sequent) (o []sequent, applied bool) {

	var newdatumSlice, repls datumSlice
	applied = false

	for k, i := range s.datumSlice() {
		n := Parse(i.String())
		if n.IsNegation() {
			if n.Child1Must().IsConditional() {
				d1 := n.Child1Must().Child1Must()
				d2 := Negate(n.Child1Must().Child2Must())
				repls = append(repls, datum(d1.String()))
				repls = append(repls, datum(d2.String()))
				newdatumSlice = slicesReplace(s.datumSlice(), k, repls)
				seq1 := mkSequent(newdatumSlice, s.succedent())
				o = append(o, datumClean(seq1))
				applied = true
				return
			}
		}
	}
	o = append(o, datumClean(s))
	return
}

// If front formula has a negated disjunction
func redr18(s sequent) (o []sequent, applied bool) {

	var newdatumSlice, repls datumSlice
	applied = false

	for k, i := range s.datumSlice() {
		n := Parse(i.String())
		if n.IsNegation() {
			if n.Child1Must().IsDisjunction() {
				d1 := Negate(n.Child1Must().Child1Must())
				d2 := Negate(n.Child1Must().Child2Must())
				repls = append(repls, datum(d1.String()))
				repls = append(repls, datum(d2.String()))
				newdatumSlice = slicesReplace(s.datumSlice(), k, repls)
				seq1 := mkSequent(newdatumSlice, s.succedent())
				o = append(o, datumClean(seq1))
				applied = true
				return
			}
		}
	}
	o = append(o, datumClean(s))
	return
}

func isReducibleSeq(s string) bool {

	parts := strings.Split(s, ";")

	if len(parts) != 2 {
		return false
	}

	for _, e := range strings.Split(parts[0], ",") {
		if e == "" {
			continue
		}
		for _, c := range greekUCBindings {
			if e == c[0] {
				return false
			}
		}
		if _, err := ParseStrict(e); err != nil {
			return false
		}
	}

	if _, err := ParseStrict(parts[1]); err != nil {
		return false
	}
	return true
}
