package gentzen

import (
	"slices"
	"strconv"
	"strings"
)

func (tks tokenStr) String() string {

	w := new(strings.Builder)

	for _, e := range tks {
		w.WriteString(e.String())
	}

	return w.String()
}

func (tks tokenStr) StringF(m PrintMode) string {

	w := new(strings.Builder)

	for _, e := range tks {
		w.WriteString(e.String())
	}

	return w.String()
}

func (tks tokenStr) wffAt(i int) tokenStr {

	if !oPL {
		return tks.wffAtSimple(i)
	}

	return tks.wffAtPL(i)
}

func (tks tokenStr) wffAtSimple(i int) (resp tokenStr) {

	openNode := 1
	j := i
	for ; j < len(tks); j++ {

		if tks[j].isBinary() {
			openNode++
		}

		if tks[j].isAtomicSentence() {
			openNode--
		}

		if tks[j].isTerm() {
			break
		}

		if tks[j].isPredicate() {
			k := 1
			for ; j+k < len(tks) && tks[j+k].isTerm(); k++ {
			}
			k--

			if k == 0 { //handle special case of Greek lower case letter as function letter
				ch, _ := getFirstChar(tks[j].str, !allowSubscr, !allowNumeral, !allowGreekUpper, !allowIdentity, !allowSpecial)
				if !isGreekLower(ch) {
					break
				}
				openNode--
				break
			}

			j = j + k
			openNode--
		}

		if tks[j].isQuantifier() {
			if j == len(tks)-1 {
				break
			}
			if !tks[j+1].isTerm() {
				break
			}
			j++
		}

		if openNode == 0 {
			break
		}

	}

	if openNode != 0 {
		return
	}

	return tks[i : j+1]
}

func (tks tokenStr) wffAtPL(i int) tokenStr {

	r := tks.wffAtSimple(i)

	for i := range r {

		e := r[i:]

		if !e[0].isQuantifier() {
			continue
		}

		v := r[i+1].str

		scope := e.subFormulas()[0]

		for j := range scope {

			e2 := scope[j:]

			if !e2[0].isQuantifier() {
				continue
			}
			if e2[1].str == v {
				return nil
			}
		}
	}

	return r
}

func (tks tokenStr) subFormulas() (resp []tokenStr) {

	if tks[0].isAtomicSentence() {
		return
	}

	if tks[0].isPredicate() {
		return
	}

	if tks[0].isTerm() {
		return
	}

	if tks[0].isQuantifier() {
		sub1 := tks.wffAt(2)
		if sub1 == nil {
			return
		}
		return append(resp, sub1)
	}

	sub1 := tks.wffAt(1)

	resp = append(resp, sub1)

	if tks[0].isUnary() {
		return resp
	}

	sub2 := tks.wffAt(1 + len(sub1))

	resp = append(resp, sub2)

	return resp
}

func (tks tokenStr) replaceFormulaAt(i int, repl tokenStr) tokenStr {

	var resp tokenStr

	oldFormula := tks.wffAt(i)

	resp = append(resp, tks[:i]...)

	resp = append(resp, repl...)

	resp = append(resp, tks[i+len(oldFormula):]...)

	return resp

}

// return index of first token of type t. Returns -1 if not found.
func (tks tokenStr) index(t tokenID) int {

	for i := range tks {

		if tks[i].tokenType == t {
			return i
		}

	}

	return -1

}

func (tks tokenStr) _isWff() bool {

	t := tks.wffAt(0)

	return len(t) == len(tks)

}

// Replace atomic sentences with subscripted sentence variables. If
// l is not supplied, p is used as the common sentence letter. Else the
// the first letter in l is used (so just supply one argument at most).
// The left most atomic sentence is p_1 and the others are named in
// ascending order.
func (tks tokenStr) normalize(l ...string) {

	if oPL {
		tks.normalizePL(l...)
	}

	var v string

	if len(l) == 0 {
		v = "p_"
	} else {
		v = l[0] + "_"
	}

	atomicS := make([]string, 0, len(tks))

	for i := range tks {

		if tks[i].tokenType != tAtomicSentence {
			continue
		}

		if slices.Contains(atomicS, tks[i].String()) {
			continue
		}

		atomicS = append(atomicS, tks[i].String())

	}

	repl := make(map[string]string, len(atomicS))

	for i, o := range atomicS {
		repl[o] = v + strconv.Itoa(i+1)
	}

	for i := range tks {

		if tks[i].tokenType != tAtomicSentence {
			continue
		}

		tks[i].str = repl[tks[i].str]

	}

	return
}

func (tks tokenStr) normalizePL(l ...string) {

	if !oPL {
		panic("normalizePL can only be used in predicate logic mode")
	}

	defP := "F_"
	defC := "c_"

	if len(l) > 0 {
		defP = l[0] + "_"
	}

	if len(l) > 1 {
		defC = l[1] + "_"
	}

	for i := range tks {

		if !tks[i].isQuantifier() {
			continue
		}

		d := tks.quantifierDepth(i)

		ov := tks[i+1].str

		nv := "var_" + strconv.Itoa(d)

		tks[i+1].str = nv

		l := len(tks[i:].subFormulas()[0])

		for j := i + 2; j < i+2+l; j++ {
			if !tks[j].isTerm() {
				continue
			}

			if tks[j-1].isQuantifier() {
				continue
			}

			if tks[j].str == ov {
				tks[j].str = nv
			}
		}
	}

	predS := make([]string, 0, len(tks))

	constS := make([]string, 0, len(tks))

	for i := range tks {

		if tks[i].isPredicate() {
			if slices.Contains(predS, tks[i].str) {
				continue
			}
			predS = append(predS, tks[i].str)
		}

		if tks[i].isTerm() {
			if strings.HasPrefix(tks[i].str, "var_") {
				continue
			}
			if slices.Contains(constS, tks[i].str) {
				continue
			}
			constS = append(constS, tks[i].str)
		}
	}

	replP := make(map[string]string, len(predS))

	for i, o := range predS {
		replP[o] = defP + strconv.Itoa(i+1)
	}

	replC := make(map[string]string, len(constS))

	for i, o := range constS {
		replC[o] = defC + strconv.Itoa(i+1)
	}

	for i := range tks {

		if tks[i].isPredicate() {
			tks[i].str = replP[tks[i].str]
		}

		if tks[i].isTerm() {
			if strings.HasPrefix(tks[i].str, "var_") {
				tks[i].str = "x_" + strings.Split(tks[i].str, "_")[1]
			} else {
				tks[i].str = replC[tks[i].str]
			}
		}
	}
}

// report the quantifier depth of quantifier token at i
func (tks tokenStr) quantifierDepth(i int) int {

	if !tks[i].isQuantifier() {
		return 0
	}

	d := 1

	j := i - 1
	for ; j > -1; j-- {
		if !tks[j].isQuantifier() {
			continue
		}
		if len(tks.wffAt(j))+j > i {
			d++
		}
	}

	return d
}

func (tks tokenStr) negate() (tkn tokenStr) {

	tkn = append(tkn, token{tokenType: tNeg, str: "N"})

	tkn = append(tkn, tks...)

	return
}

func (tks tokenStr) disjoin(tks2 tokenStr) (tkn tokenStr) {

	tkn = append(tkn, token{tokenType: tDisj, str: "A"})

	tkn = append(tkn, tks...)

	tkn = append(tkn, tks2...)

	return

}

func (tks tokenStr) conjoin(tks2 tokenStr) (tkn tokenStr) {

	tkn = append(tkn, token{tokenType: tConj, str: "K"})

	tkn = append(tkn, tks...)

	tkn = append(tkn, tks2...)

	return

}

func (tks tokenStr) isNeg() bool {

	return tks[0].tokenType == tNeg

}

func (tks tokenStr) isConj() bool {

	return tks[0].tokenType == tConj

}

func (tks tokenStr) isDisj() bool {

	return tks[0].tokenType == tDisj

}

func (tks tokenStr) isCond() bool {

	return tks[0].tokenType == tCond

}

func (tks tokenStr) isAtomic() bool {

	if len(tks) > 1 {
		return false
	}

	return tks[0].tokenType == tAtomicSentence

}

func equaltkstr(a, b tokenStr) bool {
	return a.String() == b.String()
}

func ordertkstr(a, b tokenStr) int {

	if a.String() == b.String() {
		return 0
	}

	if a.String() > b.String() {
		return -1
	}

	if a.String() < b.String() {
		return 1
	}

	return 0
}

func (tks tokenStr) isBasic() bool {

	if tks.isAtomic() {
		return true
	}

	if tks.isNeg() && tks[1:].isAtomic() {
		return true
	}

	return false
}

func (tks tokenStr) height() int {

	h := -1

	var subs, subst []tokenStr

	for subs = append(subs, tks); len(subs) != 0; {

		h++

		subst = nil

		for _, e := range subs {
			subst = append(subst, e.subFormulas()...)
		}

		subs = nil
		subs = append(subs, subst...)
	}
	return h
}
