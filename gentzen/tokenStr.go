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
func (tks tokenStr) wffAt(i int) (resp tokenStr) {

	openNode := 1
	j := i
	for ; j < len(tks); j++ {

		if tks[j].isBinary() {
			openNode++
		}

		if tks[j].isAtomicSentence() {
			openNode--
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

func (tks tokenStr) subFormulas() (resp []tokenStr) {

	if tks[0].isAtomicSentence() {
		return
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

func (tks tokenStr) isWff() bool {

	t := tks.wffAt(0)

	return len(t) == len(tks)

}

// Replace atomic sentences with subscripted sentence variables.
// The left most atomic sentence is p_1 and the others are named in
// ascending order.
func (tks tokenStr) normalize() {

	if oPL {
		return
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
		repl[o] = `p_` + strconv.Itoa(i+1)
	}

	for i := range tks {

		if tks[i].tokenType != tAtomicSentence {
			continue
		}

		tks[i].str = repl[tks[i].str]

	}

	return
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
