package gentzen

import "strings"

var theorems = [][]string{
	{"Identity", "ID", ">pp"},
	{"Non-Contradiction", "NC", "-^p-p"},
	{"Excluded Middle", "EM", "Vp-p"},
	{"Contraposition", "CP", ">>pq>-q-p"},
	{"Implication", "IM", ">>pqV-pq"},
	{"Elimination", "EL", ">Vpq>-pq"},
	{"DeMorgan", "DM", ">-Vpq^-p-q"},
	{"DeMorgan", "DM", ">-^pqV-p-q"},
	{"DeMorgan", "DM", ">V-p-q-^pq"},
	{"DeMorgan", "DM", ">^-p-q-Vpq"},
	{"Commutativity of Conjunction", "CC", ">^pq^qp"},
	{"Commutatitivity of Disjunction", "CD", ">VpqVqp"},
	{"Associativity of Conjunction", "AC", ">^^pqr^p^qr"},
	{"Associativity of Conjunction", "AC", ">^p^qr^^pqr"},
	{"Associativity of Disjunction", "AD", ">VVpqrVpVqr"},
	{"Associativity of Disjunction", "AD", ">VpVqrVVpqr"},
	{"Double Negation Introduction", "DN", ">p--p"},
}

var quantifierRules = [][]string{
	{"Quantifier Exchange", "QE", ">UxFx-Xx-Fx"},
	{"Quantifier Exchange", "QE", ">XxFx-Ux-Fx"},
	{"Quantifier Exchange", "QE", ">-Xx-FxUxFx"},
	{"Quantifier Exchange", "QE", ">-Ux-FxXxFx"},
	{"Quantifier Exchange", "QE", ">-UxFxXx-Fx"},
	{"Quantifier Exchange", "QE", ">-XxFxUx-Fx"},
	{"Quantifier Exchange", "QE", ">Xx-Fx-UxFx"},
	{"Quantifier Exchange", "QE", ">Ux-Fx-XxFx"},

	//	{"Confinement", "CF", ">^UxFxUxGxUx^FxGx"},
	//	{"Confinement", "CF", ">Ux^FxGx^UxFxUxGx"},
	//	{"Confinement", "CF", ">VUxFxUxGxUxVFxGx"},
	//	{"Confinement", "CF", ">UxVFxGxVUxFxUxGx"},
}

var modalTheorems = [][]string{
	{"Distribution", "K", ">[>pq>[p[q"},
	{"S4", "S4", ">[p[[p"},
	{"S5", "S5", "><p[<p"},
	{"Duality", "DL", ">[p-<-p"},
	{"Duality", "DL", "><p-[-p"},
	{"Duality", "DL", ">-[-p<p"},
	{"Duality", "DL", ">-<-p[p"},
	{"Duality", "DL", ">[-p-<p"},
	{"Duality", "DL", "><-p-[p"},
	{"Duality", "DL", ">-<p[-p"},
	{"Duality", "DL", ">-[p<-p"},
}

func theoremsInUse() [][]string {

	thm := theorems

	if oML {
		thm = append(thm, modalTheorems...)
	}

	if oPL {
		for i := range thm {
			thm[i][2] = strings.ReplaceAll(thm[i][2], "p", "Fx")
			thm[i][2] = strings.ReplaceAll(thm[i][2], "q", "Gx")
			thm[i][2] = strings.ReplaceAll(thm[i][2], "r", "Hx")
		}
		thm = append(thm, quantifierRules...)
	}

	return thm
}

func theorem(seq sequent, inf string) bool {
	var tf []string

	inf = strings.TrimSpace(inf)
	for _, thm := range theoremsInUse() {
		if inf == thm[1] {
			tf = append(tf, thm[2])
		}
	}

	if len(tf) == 0 {
		logger.Print(inf, "is not a theorem")
		return false
	}

	//	tf = nil
	//	tf = append(tf, quantifierRules[0][2])

	for _, thc := range tf {
		if sameStructure(thc, seq.succedent().String()) {
			return true
		}
	}

	logger.Print("not instance of ", inf)
	return false

}
