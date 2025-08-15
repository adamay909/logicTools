package gentzen

import (
	"strings"
)

func theoremsInUse() [][]string {

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
		/*		{"Distribution", "K", ">[>pq>[p[q"},
				{"S4", "S4", ">[p[[p"},
				{"S4", "S4", "><<p<p"},
				{"S5", "S5", "><p[<p"},
				{"S5", "S5", "><[p[p"},
				{"Duality", "DL", ">[p-<-p"},
				{"Duality", "DL", "><p-[-p"},
				{"Duality", "DL", ">-[-p<p"},
				{"Duality", "DL", ">-<-p[p"},
				{"Duality", "DL", ">[-p-<p"},
				{"Duality", "DL", "><-p-[p"},
				{"Duality", "DL", ">-<p[-p"},
				{"Duality", "DL", ">-[p<-p"},
				{"M1", "M1", ">[p-<-p"},
				{"M3", "M3", "><p-[-p"},
				{"M4", "M4", ">-[-p<p"},
				{"M2", "M2", ">-<-p[p"},
				{"M5", "M5", ">[pp"},
				{"\u25c7 Distribution", `\lposD`, ">[>pq><p<q"},
		*/
	}

	var axiomsSL = [][]string{
		{"", "A1", ">p>q^pq"},
		{"", "A2", ">^pqp"},
		{"", "A3", ">^pqq"},
		{"", "A4", ">pVpq"},
		{"", "A5", ">pVqp"},
		{"", "A6", ">Vpq>>pr>>qrr"},
		{"", "A7", ">>pq>>p-q-p"},
		{"", "A8", ">--pp"},
		{"", "A9", ">>p>qr>>pq>pr"},
		{"", "A10", ">p>qp"},
	}

	var thm [][]string

	thm = append(thm, theorems...)

	if oAX {
		thm = append(thm, axiomsSL...)
	}

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

	if !oPL {
		for i := range thm {
			thm[i][2] = strings.ReplaceAll(thm[i][2], "Fx", "p")
			thm[i][2] = strings.ReplaceAll(thm[i][2], "Gx", "q")
			thm[i][2] = strings.ReplaceAll(thm[i][2], "Hx", "r")
		}
	}

	return thm
}

func matchingTheorems(inf string) []string {

	var tf []string

	thms := theoremsInUse()

	inf = strings.TrimSpace(inf)
	for _, thm := range thms {
		if inf == thm[1] {
			tf = append(tf, thm[2])
		}
	}

	return tf
}

func foundMatch(c string, tf []string) bool {

	for _, thc := range tf {
		Debug("<--Theorem check: ", c, " against: ", Parse(thc, !allowGreekUpper).StringF(O_PlainText))
		if sameStructure(thc, c) {
			Debug("ok")
			Debug("--done theorem check-->")

			return true
		}
	}
	return false

}

func theorem(d *derivNode) bool {

	var tf []string

	inf := d.line.inf
	seq := d.line.seq

	tf = matchingTheorems(inf)

	if len(tf) == 0 {
		logger.Print(inf, " is not a theorem")
		return false
	}

	if foundMatch(seq.succedent().String(), tf) {
		return true
	}

	logger.Print("not instance of ", inf)

	Debug("fail")
	Debug("--done theorem check-->")

	return false

}
