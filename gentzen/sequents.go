package gentzen

import (
	"sort"
	"strings"
)

type datum string
type plshFormula string
type datumSlice []datum

type sequent struct {
	d datum
	s plshFormula
}

func (d datum) String() string {
	return string(d)
}

func (f plshFormula) String() string {
	return string(f)
}

func (s sequent) datumSlice() datumSlice {

	var r datumSlice
	dstr := string(s.d)
	if dstr == "" {
		return r
	}
	d := strings.Split(dstr, ",")

	for i := range d {
		if len(strings.TrimSpace(d[i])) == 0 {
			continue
		}
		r = append(r, datum(strings.TrimSpace(d[i])))
	}

	return r
}

func (s sequent) datum() datum {
	return s.d
}

func (s sequent) succedent() plshFormula {
	return s.s
}

func mkSequent[dat datumSlice | datum | string, fml plshFormula | *Node | string](d dat, s fml) sequent {

	var seq sequent

	switch {
	case isDatumSlice(d):
		seq.d = datum(any(d).(datumSlice).String())
	case isDatum(d):
		seq.d = any(d).(datum)
	case isString(d):
		seq.d = datum(any(d).(string))
	}

	switch {
	case isPlshFormula(s):
		seq.s = any(s).(plshFormula)
	case isString(s):
		seq.s = plshFormula(any(s).(string))
	case isNode(s):
		seq.s = plshFormula(any(s).(*Node).Formula())
	}

	return seq
}

func equalSequents(seq1, seq2 sequent) bool {

	if Parse(seq1.succedent()).Formula() != Parse(seq2.succedent()).Formula() {
		return false
	}

	return slicesEqual(seq1.datumSlice(), seq2.datumSlice())

}

func equivSequents(canonical, target sequent) bool {

	if Parse(canonical.succedent()).Formula() != Parse(target.succedent()).Formula() {
		return false
	}

	if isSeqReduce(canonical, target) {
		return true
	}

	if isSeqReorder(target, canonical) {
		return true
	}

	return false
}

func datumUnion(d ...datumSlice) datumSlice {

	var r datumSlice
	for _, e := range d {
		r = append(r, e...)
	}

	return r
}

func datumAdd[dat datum | string](d1 datumSlice, d ...dat) datumSlice {

	var r datumSlice

	r = append(r, d1...)

	for _, e := range d {

		r = append(r, datum(e))

	}

	return r
}

func datumRm[dat datum | string](d1 datumSlice, d ...dat) datumSlice {

	var r datumSlice
	r = append(r, d1...)
	for _, e := range d {
		r = slicesRemove(r, datum(e))
	}

	return r

}

// check if d1 contains d2
func datumContains(d1, d2 datumSlice) bool {

	return slicesSupset(d1, d2)

}

func datumIncludes(d1 datumSlice, d datum) bool {

	return slicesContains(d1, d)

}

func (d datumSlice) StringSlice() []string {

	var r []string

	for _, e := range d {
		r = append(r, string(e))
	}

	return r
}

func (d datumSlice) String() string {

	r := strings.Join(d.StringSlice(), ",")

	r = strings.TrimPrefix(r, ",")

	r = strings.TrimSuffix(r, ",")

	return r
}

func (d datumSlice) datum() datum {

	return datum(d.String())

}

func datumsEqual(d1, d2 datumSlice) bool {

	return slicesEqual(d1, d2)

}

func datumsEquiv(want, have datumSlice) bool {

	dummy := plshFormula("P")
	if oPL {
		dummy = plshFormula("Fx")
	}

	w1 := sequent{want.datum(), dummy}
	h2 := sequent{have.datum(), dummy}

	return equivSequents(w1, h2)

}

func (d datumSlice) Less(i, j int) bool {
	return string(d[i]) < string(d[j])
}

func datumSort(d datumSlice) {

	sort.Slice(d, func(i, j int) bool { return string(d[i]) < string(d[j]) })

	return
}
