package gentzen

func SLform(n *Node) *Node {

	if oPL {
		logger.Print("some functionality not available for Predicate Logic")
		return n
	}

	atomic := n.AtomicSentences()

	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "w", "y", "z"}
	count := 0

	var nextLetter func() string
	nextLetter = func() string {
		ret := letters[count]
		count++
		if count > len(letters)-1 {
			return "K"
		}
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

			if e.IsQuantifier() || e.IsModal() {
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

	if oPL {
		logger.Print("some functionality not available for Predicate Logic")
		return false
	}
	s1 = Parse(s1).Formula()
	s2 = Parse(s2).Formula()

	n1 := SLform(Parse(lconj + lcond + s1 + s2 + lcond + s2 + s1))

	if !IsTautology(n1.Formula()) {
		return false
	}

	return true
}

func sententialLogic(s1, s2 sequent) bool {

	if !oML {
		logger.Print("appeal to SL only allowed with Modal Logic.")
		return false
	}

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
