package gentzen

import (
	"slices"
	"strconv"
	"strings"
)

type psequent struct {
	datum []tokenStr
	succ  tokenStr
}

type reductionNode struct {
	seq      psequent
	child    *reductionNode
	parent   *reductionNode
	depend   []*reductionNode
	terminal bool
	number   int
}

/*
PrintProofOutline returns a string representing a derivation
of ⊢ s. If the derivation is not a valid proof, you will see
"Premise" (as opposed to "Assumption") in the annotation of at least one line.
*/
func PrintProofOutline(s string, mode PrintMode) string {

	n := genProofTree(s)

	return n.printDeriv(mode)

}

func genProofTree(s string) *reductionNode {

	var seq sequent

	seq.d = datum("")
	seq.s = plshFormula(s)

	n := generateProofTree(seq)

	return n.top()
}

func isTautologySequentReduction(s string) bool {

	var seq sequent

	seq.d = datum("")
	seq.s = plshFormula(s)

	n := generateProofTree(seq)

	for e := n.bottom(); e != nil; e = e.parent {

		if !e.terminal {
			continue
		}

		if !e.isAssumption() {
			return false
		}

	}
	return true
}

func generateProofTree(s sequent) *reductionNode {

	var err error

	var sq psequent

	datums := s.datumSlice()

	tks := make(tokenStr, 0, 10000)

	for _, d := range datums {
		tks, err = tokenize(string(d), allowGreekUpper, !allowSpecial)

		if err != nil {
			panic(s.String() + " is not a well formed sequent")
		}

		sq.datum = append(sq.datum, tks)

		clear(tks)
	}

	tks, err = tokenize(string(s.s), !allowGreekUpper, !allowSpecial)

	if err != nil {
		panic(s.String() + " is not a well formed sequent")
	}

	sq.succ = tks

	n := new(reductionNode)

	n.seq = sq

	tks = nil

	for e := n; ; e = e.child {

		e.growProofTree()
		if e.child == nil {
			break
		}
	}

	return n.top()
}

func (n *reductionNode) top() *reductionNode {

	for e := n; e != nil; e = e.parent {
		if e.parent == nil {
			return e
		}
	}

	return nil
}

func (n *reductionNode) bottom() *reductionNode {

	for e := n; e != nil; e = e.child {
		if e.child == nil {
			return e
		}
	}

	return nil
}

func (n *reductionNode) linearize() []*reductionNode {

	var resp []*reductionNode

	for e := n; e != nil; e = e.child {

		resp = append(resp, e)

	}
	return resp

}

func (n *reductionNode) isBasic() bool {

	sq := n.seq

	if len(sq.datum) > 1 {
		return false
	}

	if !n.seq.succ.isBasic() {
		return false
	}

	if len(sq.datum) == 0 || sq.datum[0].isBasic() {
		return true
	}

	return false
}

func (n *reductionNode) growProofTree() {

	if n.isBasic() {
		n.terminal = true
	}

	if n.terminal {
		return
	}

	var cseq1, cseq2 psequent

	// remove redundant datum items
	target := n.seq.succ.String()

	for _, d := range n.seq.datum {

		if d.String() == target {

			cseq1.addDatum(d)

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		}
	}

	//ECQ on datum side

	for _, d := range n.seq.datum {

		if !d.isAtomic() {
			continue
		}

		for _, e := range n.seq.datum {

			if equaltkstr(d.negate(), e) {

				cseq1.addDatum(d)

				cseq1.succ = d

				cseq2.addDatum(d.negate())

				cseq2.succ = d.negate()

				n.addChild(cseq1)

				n.addChild(cseq2)

				return
			}
		}
	}

	//handle datum side

	for i, d := range n.seq.datum {

		switch {

		case d.isConj():

			subs := d.subFormulas()

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(subs...)

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		case d.isDisj():

			subs := d.subFormulas()

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(subs[0])

			cseq2.setDatum(removeElement(n.seq.datum, i))

			cseq2.addDatum(subs[1])

			cseq1.succ = n.seq.succ

			cseq2.succ = n.seq.succ

			n.addChild(cseq1)
			n.addChild(cseq2)

			return

		case d.isCond():

			subs := d.subFormulas()

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(subs[0].negate().disjoin(subs[1]))

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		case d.isNeg() && d[1:].isNeg():

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(d[2:])

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		case d.isNeg() && d[1:].isConj():

			subs := d[1:].subFormulas()

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(subs[0].negate().disjoin(subs[1].negate()))

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		case d.isNeg() && d[1:].isDisj():

			subs := d[1:].subFormulas()

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(subs[0].negate().conjoin(subs[1].negate()))

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		case d[0].tokenType == tNeg && d[1:].isCond():

			subs := d[1:].subFormulas()

			cseq1.setDatum(removeElement(n.seq.datum, i))

			cseq1.addDatum(subs[0].conjoin(subs[1].negate()))

			cseq1.succ = n.seq.succ

			n.addChild(cseq1)

			return

		}
	}

	//handle succedent side
	switch {

	case n.seq.succ.isConj():

		subs := n.seq.succ.subFormulas()

		cseq1.addDatum(n.seq.datum...)

		cseq1.succ = subs[0]

		cseq2.addDatum(n.seq.datum...)

		cseq2.succ = subs[1]

		n.addChild(cseq1)

		n.addChild(cseq2)

		return

	case n.seq.succ.isCond():

		subs := n.seq.succ.subFormulas()

		cseq1.setDatum(n.seq.datum)

		cseq1.succ = subs[0].negate().disjoin(subs[1])

		n.addChild(cseq1)

		return

	case n.seq.succ.isDisj():

		subs := n.seq.succ.subFormulas()

		cseq1.setDatum(n.seq.datum)

		cseq1.succ = subs[0].negate().conjoin(subs[1].negate()).negate()

		n.addChild(cseq1)

		return

	case n.seq.succ.isNeg():

		if n.seq.succ[1:].isAtomic() {
			break
		}

		cseq1.setDatum(n.seq.datum)
		cseq2.setDatum(n.seq.datum)

		cseq1.addDatum(n.seq.succ[1:])
		cseq2.addDatum(n.seq.succ[1:])

		for _, e := range n.seq.succ {

			if e.tokenType == tAtomicSentence {
				cseq1.succ = []token{e}
				cseq2.succ = cseq1.succ.negate()
				break
			}
		}

		n.addChild(cseq1)
		n.addChild(cseq2)

		return

	}
	n.terminal = true

}

func (n *reductionNode) addChild(c psequent) {

	gp := n.bottom()

	nc := newreductionnode(c)
	gp.child = nc
	nc.parent = gp
	nc.cleanDatum()

	n.depend = append(n.depend, nc)

}

func newreductionnode(pseq psequent) *reductionNode {

	n := new(reductionNode)

	n.seq = pseq

	return n
}

func (n *reductionNode) String() string {

	w := new(strings.Builder)

	ds := make([]string, 0, len(n.seq.datum))

	for _, d := range n.seq.datum {

		ds = append(ds, d.String())
	}

	w.WriteString(strings.Join(ds, ","))

	w.WriteString(" ; ")

	w.WriteString(n.seq.succ.String())

	return w.String()

}

func (n *reductionNode) printDeriv(mode PrintMode) string {

	w := new(strings.Builder)

	var derivlines []string

	lines := n.linearize()

	slices.Reverse(lines)

	for i, l := range lines {

		l.number = i + 1

	}

	for _, l := range lines {

		w.WriteString(l.String())

		w.WriteString(";")

		lannot := make([]int, 0, len(l.depend))

		for _, e := range l.depend {

			lannot = append(lannot, e.number)

		}

		slices.Sort(lannot)

		for k, e := range lannot {
			w.WriteString(strconv.Itoa(e))
			if k < len(lannot)-1 {
				w.WriteString(",")
			}
		}

		if len(l.depend) == 0 {

			if l.isAssumption() {
				w.WriteString("Assumption")
			} else {
				w.WriteString("Premise")
			}
		}

		derivlines = append(derivlines, w.String())

		w.Reset()

	}

	return PrintDerivation(derivlines, 1, mode)
}

func (n *reductionNode) isAssumption() bool {

	if len(n.seq.datum) != 1 {
		return false
	}

	return n.seq.datum[0].String() == n.seq.succ.String()
}

func (n *reductionNode) cleanDatum() {

	n.seq.cleanDatum()

}

func (p *psequent) cleanDatum() {

	slices.SortFunc(p.datum, ordertkstr)

	nd := slices.CompactFunc(p.datum, equaltkstr)

	p.datum = nd
}

func (p *psequent) setDatum(d []tokenStr) {

	p.datum = nil

	p.datum = append(p.datum, d...)

}

func (p *psequent) addDatum(d ...tokenStr) {

	p.datum = append(p.datum, d...)

	p.cleanDatum()

}

func removeElement[S ~[]E, E any](s S, i int) (resp S) {

	if i == 0 {

		resp = append(resp, s[1:]...)

		return

	}

	resp = append(resp, s[:i]...)

	if i < len(s)-1 {

		resp = append(resp, s[i+1:]...)

	}

	return resp

}

func (p *psequent) removeDatum(d tokenStr) {

	f := func(e tokenStr) bool {
		return equaltkstr(e, d)
	}
	p.datum = slices.DeleteFunc(p.datum, f)

	p.cleanDatum()

}

func (n *reductionNode) isTerminal() bool {

	for _, d := range n.seq.datum {
		if !d.isBasic() {
			return false
		}
	}

	return n.seq.succ.isBasic()

}
