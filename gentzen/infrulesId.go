package gentzen

import "strings"

func idI(d *derivNode) bool {

	if len(d.supportingLines) != 0 {
		logger.Print("Identity Introduction depends on no other lines")
		return false
	}

	seq := d.line.seq

	if strings.TrimSpace(seq.datum().String()) != "" {
		logger.Print("datum must be empty")
		return false
	}

	n := Parse(seq.succedent().String(), !allowGreekUpper)

	if !n.Check(isAtomic) {
		logger.Print("must be atomic identity statement")
		return false
	}

	if n.Val.predicateLetter != "=" {
		logger.Print("must be atomic statement")
		return false
	}

	if len(n.Val.term) != 2 {
		logger.Print("identity is a 2-place relation")
		return false
	}

	if n.Val.term[0] != n.Val.term[1] {
		logger.Print("must assert identity with self")
		return false
	}

	return true

}

func idE(d *derivNode) bool {

	if len(d.supportingLines) != 0 {
		logger.Print("Identity Elimination depends on no other lines")
		return false
	}

	seq := d.line.seq

	if strings.TrimSpace(seq.datum().String()) != "" {
		logger.Print("datum must be empty")
		return false
	}

	n := Parse(seq.succedent().String(), !allowGreekUpper)

	if n.Val.connective != Cond {
		logger.Print("main connective must be conditional")
		return false
	}

	if n.Child(0).Val.connective != Conj {
		logger.Print("antecedent must be conjunction")
		return false
	}

	if !n.Child(0).Child(0).Check(isIdentity) {
		logger.Print("first conjunct must be identity")
		return false
	}

	k1 := n.Child(0).Child(0).Val.term[0]
	k2 := n.Child(0).Child(0).Val.term[1]

	s1 := n.Child(0).Child(1)
	s2 := n.Child(1)

	s3 := replaceTerms(s2, k2, k1)

	if s1.String() != s3.String() {
		logger.Print("consequent and second conjunct don't match up in the right way")
		return false
	}
	return true
}
