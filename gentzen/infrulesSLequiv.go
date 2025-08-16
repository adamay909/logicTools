package gentzen

import "strconv"

func slForm(n *Node) *Node {

	if oPL {
		logger.Print("some functionality not available for Predicate Logic")
		return n
	}

	atomic := atomicSentences(n)

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

	ns := n.Linearize()
	found := true
	for found {
		found = false
		for _, e := range ns {

			if !e.HasFlag("c") && (e.Check(isQuantifier) || e.Check(isAtomic)) {
				s := nextLetter()
				target := e.String()
				for j := range ns {
					if ns[j].String() == target && !ns[j].HasFlag("c") {
						ns[j].Val.connective = None
						ns[j].Val.raw = s
						ns[j].SetFlag("c")
					}
				}
				found = true
				break
			}
		}
		ns = ns[0].Linearize()
	}
	return ns[0]
}

func equivSL(s1, s2 string) bool {

	restorePL := false

	n1 := Parse(s1, !allowGreekUpper)
	n2 := Parse(s2, !allowGreekUpper)

	if oPL {
		restorePL = true
		oPL = false
	}

	s1 = slForm(n1).String()
	s2 = slForm(n2).String()

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

	if !datumsEquiv(s1.datumSlice(), s2.datumSlice()) {
		logger.Print("datum cannot change")
		return false
	}

	return true
}
