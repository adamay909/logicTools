package gentzen

import "strconv"

func SLform(n *Node) *Node {

	if oPL {
		logger.Print("some functionality not available for Predicate Logic")
		return n
	}

	atomic := n.AtomicSentences()

	count := 0

	var nextLetter func() string
	nextLetter = func() string {
		ret := "p_" + strconv.Itoa(count)
		count++
		if slicesContains(atomic, ret) {
			ret = nextLetter()
		}
		return ret
	}

	ns := getSubnodes(n)
	found := true
	for found {
		found = false
		for _, e := range ns {

			if !e.HasFlag("c") && (e.IsQuantifier() || e.IsModal() || e.IsAtomic()) {
				s := nextLetter()
				target := e.Formula()
				for j := range ns {
					if ns[j].Formula() == target && !ns[j].HasFlag("c") {
						ns[j].SetAtomic()
						ns[j].SetFormula(s)
						ns[j].SetFlag("c")
					}
				}
				found = true
				break
			}
		}
		ns = getSubnodes(ns[0])
	}
	return ns[0]
}

func equivSL(s1, s2 string) bool {

	restorePL := false

	n1 := Parse(s1)
	n2 := Parse(s2)

	if oPL {
		restorePL = true
		oPL = false
	}

	s1 = SLform(n1).Formula()
	s2 = SLform(n2).Formula()

	s3 := lconj + lcond + s1 + s2 + lcond + s2 + s1

	if !IsTautology(s3) {
		oPL = restorePL
		return false
	}

	oPL = restorePL
	return true
}

func sententialLogic(d *derivNode) bool {

	if !oML {
		logger.Print("Appeal to logic only allowed with Modal Logic")
		return false
	}

	if len(d.supportingLines) != 1 {
		logger.Print("Appeal to logic depends on one line")
		return false
	}

	s1 := d.supportingLines[0].line.seq
	s2 := d.line.seq

	if !equivSL(s1.succedent().String(), s2.succedent().String()) {
		logger.Print("succedents not equivalent in Sentential Logic")
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
