package gentzen

import "strings"

var theorems = [][]string{
	{"Identity", "ID", ">pp"},
	{"Non-Contradiction", "NC", "-^p-p"},
	{"Excluded Middle", "EM", "Vp-p"},
	{"Contraposition", "CP", ">>pq>-q-p"},
	{"Implication", "IM", ">>pqV-pq"},
	{"Elimination", "EL", ">Vpq>-pq"},
	{"DeMorgan 1", "DM1", ">-Vpq^-p-q"},
	{"DeMorgan 2", "DM2", ">-^pqV-p-q"},
	{"DeMorgan 3", "DM3", ">V-p-q-^pq"},
	{"DeMorgan 4", "DM4", ">^-p-q-Vpq"},
	{"Commutativity of Conjunction", "CC", ">^pq^qp"},
	{"Commutatitivity of Disjunction", "CD", ">VpqVpq"},
	{"Associativity of Conjunction", "AC", ">^^pqr^p^qr"},
	{"Associativity of Disjunction", "AD", ">VVpqrVpVqr"},
}

func theorem(seq sequent, inf string) bool {
	var tf string
	thm := theorems

	if oPL {
		for i := range thm {
			thm[i][2] = strings.ReplaceAll(thm[i][2], "p", "Fx")
			thm[i][2] = strings.ReplaceAll(thm[i][2], "q", "Gx")
			thm[i][2] = strings.ReplaceAll(thm[i][2], "r", "Hx")
		}
	}
	inf = strings.TrimSpace(inf)
	for i := range thm {
		if inf == theorems[i][0] || inf == theorems[i][1] {
			tf = theorems[i][2]
		}
	}

	if tf == "" {
		logger.Print(inf, "is not a theorem")
		return false
	}

	if !sameStructure(tf, seq.succedent().String()) {
		logger.Print("not instance of ", inf)
		return false
	}
	return true

}
