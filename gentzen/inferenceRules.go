package gentzen

import (
	"sort"
	"strings"
)

type sequent struct {
	datum, succedent string
}

var strictCheck bool

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

func checkLineRefs(cur, need int, lines []int) bool {

	return true
}

func assumption(seq sequent) bool {

	if strings.TrimSpace(seq.datum) != strings.TrimSpace(seq.succedent) {
		logger.Print("datum and subseqent cannot differ for assumption")
		return false
	}
	return true

}

func condE(seq1, seq2, seq3 sequent) bool {

	v1, msg1 := condEhelper(seq1, seq2, seq3)

	v2, msg2 := condEhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "no conditional found" {

		logger.Print(msg2)
		return false

	}

	logger.Print(msg1)
	return false

}

func condEhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)
	n3 := Parse(seq3.succedent)
	v = false

	if n1.MainConnective() != cond {
		msg = "no conditional found"
		return
	}

	if n1.subnode1.Formula() != n2.Formula() {
		msg = "antecedent of conditional does not match succedent of other line"
		return
	}

	if n3.Formula() != n1.subnode2.Formula() {
		msg = "consequent of conditional does not match conclusion"
		return
	}

	if strictCheck {
		if !equal(add(seq2.datum, seq1.datum), seq3.datum) {
			msg = "datum of conclusion must be union of datums of premises"
			return
		}
	} else {
		if !equal(datumReduce(add(seq2.datum, seq1.datum)), datumReduce(seq3.datum)) {
			msg = "datum of conclusion must be union of datums of premises"
			return
		}
	}
	v = true

	return

}

func condI(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)

	if n2.MainConnective() != cond {
		logger.Print("main connective of conclusion must be conditional")
		return false
	}

	if !contains(n2.subnode1.Formula(), seq1.datum) {
		logger.Print("antecedent of conclusion must be in datum of premise")
		return false
	}

	if n2.subnode2.Formula() != n1.Formula() {
		logger.Print("consequent of conclusion must be succedent of premise")
		return false
	}
	if strictCheck {
		if !equal(rm(seq1.datum, n2.subnode1.Formula()), seq2.datum) {
			logger.Print("must remove one datum item")
			return false
		}
	} else {
		if !equal(datumReduce(rm(seq1.datum, n2.subnode1.Formula())), datumReduce(seq2.datum)) {
			logger.Print("must remove one datum item")
			return false
		}
	}
	return true
}

func conjE(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)

	if n1.MainConnective() != conj {
		logger.Print("must start with conjunction")
		return false
	}

	if n2.Formula() != n1.subnode1.Formula() {
		if n2.Formula() != n1.subnode2.Formula() {
			logger.Print("conclusion not one of conjuncts")
			return false
		}
	}
	if strictCheck {
		if !equal(seq1.datum, seq2.datum) {
			logger.Print("datum of conclusion must be same as datum of premise")
			return false
		}
	} else {
		if !equal(datumReduce(seq1.datum), datumReduce(seq2.datum)) {
			logger.Print("datum of conclusion must be same as datum of premise")
			return false
		}
	}
	return true
}

func conjI(seq1, seq2, seq3 sequent) bool {

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)
	n3 := Parse(seq3.succedent)

	if n3.MainConnective() != conj {
		logger.Print("conclusion must be a conjunction")
		return false
	}

	if n1.Formula() != n3.subnode1.Formula() && n1.Formula() != n3.subnode2.Formula() {
		logger.Print("first premise not conjunct of conclusion")

		return false
	}

	if n2.Formula() != n3.subnode1.Formula() && n2.Formula() != n3.subnode2.Formula() {
		logger.Print("second premise not conjunct of conclusion")
		return false
	}

	if strictCheck {
		if !equal(add(seq1.datum, seq2.datum), seq3.datum) {
			logger.Print("datum of conclusion must be union of datums of premise")
			return false
		}
	} else {
		if !equal(datumReduce(add(seq1.datum, seq2.datum)), datumReduce(seq3.datum)) {
			logger.Print("datum of conclusion must be union of datums of premise")
			return false
		}
	}

	return true
}

func disjI(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)

	if n2.MainConnective() != disj {
		logger.Print("conclusion must be a disjunction")
		return false
	}

	if n2.subnode1.Formula() != n1.Formula() && n2.subnode2.Formula() != n1.Formula() {
		logger.Print("premise is not one of disjuncts")
		return false
	}
	if strictCheck {
		if !equal(seq1.datum, seq2.datum) {
			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !equal(datumReduce(seq1.datum), datumReduce(seq2.datum)) {
			logger.Print("datum cannot change")
			return false
		}
	}

	return true
}

func disjE(seq1, seq2, seq3, seq4 sequent) bool {

	v1, msg1 := disjEhelper1(seq1, seq2, seq3, seq4)
	v2, msg2 := disjEhelper1(seq2, seq1, seq3, seq4)
	v3, msg3 := disjEhelper1(seq3, seq1, seq2, seq4)

	if !v1 && !v2 && !v3 {
		if msg1 == "must have disjunction among premises" {

			if msg2 == msg1 {
				logger.Print(msg3)
				return false
			}

			logger.Print(msg2)
			return false
		}
		logger.Print(msg1)
		return false
	}

	var v bool
	switch {

	case v1:
		v, msg2 = disjEhelper2(seq1, seq2, seq3, seq4)

	case v2:
		v, msg2 = disjEhelper2(seq2, seq1, seq3, seq4)

	case v3:
		v, msg2 = disjEhelper2(seq3, seq1, seq2, seq4)

	}

	if !v {
		logger.Print(msg2)
		return false
	}

	return true
}

func disjEhelper1(seq1, seq2, seq3, seq4 sequent) (v bool, msg string) {

	v = false

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)
	n3 := Parse(seq3.succedent)
	n4 := Parse(seq4.succedent)

	if n1.MainConnective() != disj {
		msg = "must have disjunction among premises"
		return
	}

	if n2.Formula() != n3.Formula() {
		msg = "must have two identical succedents"
		return
	}

	if n3.Formula() != n4.Formula() {
		msg = "conclusion must be identical to two of the succedents"
		return
	}

	d1 := n1.subnode1.Formula()
	d2 := n1.subnode2.Formula()

	if !contains(d1, seq2.datum) && !contains(d1, seq3.datum) {
		msg = "one of disjuncts not in datums"
		return
	}

	if !contains(d2, seq2.datum) && !contains(d2, seq3.datum) {
		msg = "one of disjuncts not in datums"
		return
	}

	if contains(d1, seq2.datum) {
		if !contains(d2, seq3.datum) {
			msg = "one of disjuncts not in datums"
			return
		}
	}

	if contains(d2, seq2.datum) {
		if !contains(d1, seq3.datum) {
			msg = "one of disjuncts not in datums"
			return
		}
	}

	v = true
	return
}

func disjEhelper2(seq1, seq2, seq3, seq4 sequent) (v bool, msg string) {

	v = false

	ndatum := add(add(seq1.datum, seq2.datum), seq3.datum)

	n1 := Parse(seq1.succedent)
	d1 := n1.subnode1.Formula()
	d2 := n1.subnode2.Formula()
	ndatum = rm(rm(ndatum, d1), d2)

	if strictCheck {
		if !equal(ndatum, seq4.datum) {
			msg = "datum of conclusion incorrect"
			return
		}
	} else {
		if !equal(datumReduce(ndatum), datumReduce(seq4.datum)) {
			msg = "datum of conclusion incorrect"
			return
		}
	}
	v = true

	return
}

func negE(seq1, seq2 sequent) bool {

	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)

	if n1.MainConnective() != neg {
		logger.Print("premise must be double negation")
		return false
	}

	if n1.subnode1.MainConnective() != neg {
		logger.Print("premise must be double negation")
		return false
	}

	if n1.subnode1.subnode1.Formula() != n2.Formula() {
		logger.Print("conclusion is not the elimnation of double negation")
		return false
	}
	if strictCheck {
		if !equal(seq1.datum, seq2.datum) {
			logger.Print("datum must remain same")
			return false
		}
	} else {
		if !equal(datumReduce(seq1.datum), datumReduce(seq2.datum)) {
			logger.Print("datum must remain same")
			return false
		}
	}

	return true
}

func negI(seq1, seq2, seq3 sequent) bool {

	v1, msg1 := negIhelper(seq1, seq2, seq3)

	v2, msg2 := negIhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "conclusion must be negation" {
		logger.Print(msg2)
		return false
	}
	logger.Print(msg1)

	return false
}

func negIhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	v = false
	n1 := Parse(seq1.succedent)
	n2 := Parse(seq2.succedent)
	n3 := Parse(seq3.succedent)

	if n3.MainConnective() != neg {
		msg = "conclusion must be negation"
		return
	}

	if n1.MainConnective() != neg {
		msg = `premises's succedents must negate each other`
		return
	}

	if n1.subnode1.Formula() != n2.Formula() {
		msg = `premises's succedents must negate each other`
		return
	}

	f := n3.subnode1.Formula()

	if !contains(f, seq1.datum) {
		msg = "conclusion must be negation of something in common between the datums of premises"

		return
	}

	if !contains(f, seq2.datum) {
		msg = "conclusion must be in datums of both premises"
		return
	}
	if strictCheck {
		if !equal(add(rm(seq1.datum, f), rm(seq2.datum, f)), seq3.datum) {
			msg = "datum of conclusion incorrect"
			return
		}
	} else {
		if !equal(datumReduce(add(rm(seq1.datum, f), rm(seq2.datum, f))), datumReduce(seq3.datum)) {
			msg = "datum of conclusion incorrect"
			return
		}
	}
	v = true
	return
}

func uniE(seq1, seq2 sequent) bool {

	d1 := seq1.datum
	d2 := seq2.datum

	if strictCheck {
		if !equal(d1, d2) {
			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !equal(datumReduce(d1), datumReduce(d2)) {
			logger.Print("datum cannot change")
			return false
		}
	}

	if Parse(seq1.succedent).MainConnective() != uni {
		logger.Print("premise must be universally quantified")
		return false
	}

	val, _, _ := isInstanceOf(seq2.succedent, seq1.succedent)
	if !val {
		logger.Print("conclusion not an instance of premise")
	}
	return val
}

func exI(seq1, seq2 sequent) bool {

	d1 := seq1.datum
	d2 := seq2.datum

	if strictCheck {
		if !equal(d1, d2) {
			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !equal(datumReduce(d1), datumReduce(d2)) {
			logger.Print("datum cannot change")
			return false
		}
	}

	if Parse(seq2.succedent).MainConnective() != ex {
		logger.Print("conclusion must be existentially quantified")
		return false
	}

	val, _, _ := isInstanceOf(seq1.succedent, seq2.succedent)
	if !val {
		logger.Print("conclusion must be existential generalization of premise")
	}
	return val
}

func uniI(seq1, seq2 sequent) bool {

	d1 := seq1.datum
	d2 := seq2.datum

	if strictCheck {
		if !equal(d1, d2) {

			logger.Print("datum cannot change")
			return false
		}
	} else {
		if !equal(datumReduce(d1), datumReduce(d2)) {
			logger.Print("datum cannot change")
			return false
		}
	}

	if Parse(seq2.succedent).MainConnective() != uni {
		logger.Print("conclusion must be universally quantified")
		return false
	}

	val, _, term := isInstanceOf(seq1.succedent, seq2.succedent)

	if !val {
		logger.Print("conclusion not a universalization of premise")
		return false
	}

	datums := strings.Split(seq1.datum, ",")

	for _, d := range datums {
		if len(d) == 0 {
			continue
		}
		if string(d[0]) == `\` {
			continue
		}
		if Parse(d).hasTerm(term) {
			logger.Print(term, " cannot appear in datum")
			return false
		}
	}

	return true

}

func exE(seq1, seq2, seq3 sequent) bool {

	v1, msg1 := exEhelper(seq1, seq2, seq3)

	v2, msg2 := exEhelper(seq2, seq1, seq3)

	if v1 || v2 {
		return true
	}

	if msg1 == "no existential generalization in premise" {
		logger.Print(msg2)
		return false
	}
	logger.Print(msg1)
	return false

}

func exEhelper(seq1, seq2, seq3 sequent) (v bool, msg string) {

	v = false

	if Parse(seq1.succedent).MainConnective() != ex {
		msg = "no existential generalization in premise"
		return
	}

	if seq2.succedent != seq3.succedent {
		msg = "conclusion does not match premises"
		return
	}

	found := false
	var kappa string

	datums2 := strings.Split(seq2.datum, ",")

	for _, d := range datums2 {
		if d[:1] == `\` {
			continue
		}
		found, _, kappa = isInstanceOf(d, seq1.succedent)
		if found {
			datums2 = strings.Split(rm(seq2.datum, d), ",")
			break
		}
	}

	if !found {
		msg = "no datum item found as instance of existential claim"
		return
	}

	datums1 := strings.Split(seq1.datum, ",")
	for _, d := range datums1 {
		if len(d) == 0 {
			continue
		}
		if d[:1] == `\` {
			continue
		}
		if Parse(d).hasTerm(kappa) {
			msg = kappa + " may not appear in any datum items"

			return
		}
	}

	for _, d := range datums2 {
		if len(d) == 0 {
			continue
		}
		if d[:1] == `\` {
			continue
		}
		if Parse(d).hasTerm(kappa) {
			msg = kappa + " may not appear in any datum items"
			return
		}
	}
	if strictCheck {
		if !equal(add(strings.Join(datums1, ","), strings.Join(datums2, ",")), seq3.datum) {
			msg = "datum of conclusion must be union of datums of premise"
			v = false
			return
		}
	} else {
		if !equal(datumReduce(add(strings.Join(datums1, ","), strings.Join(datums2, ","))), datumReduce(seq3.datum)) {
			msg = "datum of conclusion must be union of datums of premise"
			v = false
			return
		}
	}
	v = true
	msg = ""
	return

}

func idI(seq sequent) bool {

	if strings.TrimSpace(seq.datum) != "" {
		logger.Print("datum must be empty")
		return false
	}

	n := Parse(seq.succedent)

	if !n.IsAtomic() {
		logger.Print("must be atomic identity statement")
		return false
	}

	if n.predicateLetter != "=" {
		logger.Print("must be atomic statement")
		return false
	}

	if len(n.Terms()) != 2 {
		logger.Print("identity is a 2-place relation")
		return false
	}

	if n.Terms()[0] != n.Terms()[1] {
		logger.Print("must assert identity with self")
		return false
	}

	return true

}

func idE(seq sequent) bool {

	if strings.TrimSpace(seq.datum) != "" {
		logger.Print("datum must be empty")
		return false
	}

	n := Parse(seq.succedent)

	if n.MainConnective() != cond {
		logger.Print("main connective must be conditional")
		return false
	}

	if n.subnode1.MainConnective() != conj {
		logger.Print("antecedent must be conjunction")
		return false
	}

	if !n.subnode1.subnode1.IsIdentity() {
		logger.Print("first conjunct must be identity")
		return false
	}

	k1 := n.subnode1.subnode1.Terms()[0]
	k2 := n.subnode1.subnode1.Terms()[1]

	s1 := n.subnode1.subnode2
	s2 := n.subnode2

	s3 := replaceTerms(s2, k2, k1)

	if s1.String() != s3.String() {
		logger.Print("consequent and second conjunct don't match up in the right way")
		return false
	}
	return true
}

func seqRewrite(seq1, seq2 sequent, n int) bool {

	d1 := datumReduce(seq1.datum)
	d2 := datumReduce(seq2.datum)
	if equal(d1, d2) {
		return true
	}
	if contains(d1, d2) {
		return true
	}
	if contains(d2, d1) {
		return true
	}
	logger.Print("not a rewrite of line ", n)
	return false

}

func datumReduce(d string) (out string) {

	d = strings.ReplaceAll(d, " ", "")
	data := strings.Split(d, ",")
	if len(data) == 0 {
		return d
	}

	sort.Strings(data)

	var ndata []string

	ndata = append(ndata, data[0])

	for i := 1; i < len(data); i++ {
		if ndata[len(ndata)-1] != data[i] {
			ndata = append(ndata, data[i])
		}
	}
	for i := range ndata {
		out = out + ndata[i] + ","
	}

	return strings.TrimRight(out, ",")
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

	if !sameStructure(tf, seq.succedent) {
		logger.Print("not instance of ", inf)
		return false
	}
	return true

}

func standardize(f string) string {

	var atomic []string

	var repl = []string{"p", "q", "r", "s", "t", "u", "x", "y", "z"}
	for i := 0; i < len(f); i++ {

		if !isConnective(f[i : i+1]) {
			if !in(f[i:i+1], atomic) {
				atomic = append(atomic, f[i:i+1])
			}
		}
	}
	if len(atomic) > len(repl) {
		return f
	}

	for i, x := range atomic {
		f = strings.ReplaceAll(f, x, repl[i])
	}
	return f
}

func isSame(d1, d2 []string) bool {

	if len(d1) != len(d2) {
		return false
	}

	for i := range d1 {
		if d1[i] != d2[i] {
			return false
		}
	}

	return true
}

func equal(d1, d2 string) bool {

	d1 = strings.ReplaceAll(d1, " ", "")
	d2 = strings.ReplaceAll(d2, " ", "")

	data1 := strings.Split(d1, ",")
	data2 := strings.Split(d2, ",")

	if len(data1) != len(data2) {
		return false
	}

	sort.Strings(data1)
	sort.Strings(data2)

	for i := 0; i < len(data1); i++ {
		if data1[i] != data2[i] {
			return false
		}
	}

	return true
}

//remove d2 from d1
func rm(d1, d2 string) string {

	d1 = strings.ReplaceAll(d1, " ", "")
	d2 = strings.ReplaceAll(d2, " ", "")

	data1 := strings.Split(d1, ",")
	data2 := strings.Split(d2, ",")

	var data3 []string

	for i := range data1 {
		if !in(data1[i], data2) {
			data3 = append(data3, data1[i])
		}
	}
	var d3 string

	for _, j := range data3 {
		d3 = d3 + j + ","
	}
	d3 = strings.TrimRight(d3, ",")
	return d3
}

func add(d1, d2 string) string {

	d1 = strings.ReplaceAll(d1, " ", "")
	d2 = strings.ReplaceAll(d2, " ", "")
	if d1 == "" {
		return d2
	}

	if d2 == "" {
		return d1
	}

	data1 := strings.Split(d1, ",")
	data2 := strings.Split(d2, ",")

	for _, i := range data2 {
		data1 = append(data1, i)
	}

	sort.Strings(data1)

	var d3 string

	for _, j := range data1 {
		d3 = d3 + j + ","
	}
	d3 = strings.TrimRight(d3, ",")
	return d3
}

func contains(d string, dg string) bool {

	dg = strings.ReplaceAll(dg, " ", "")
	data1 := strings.Split(dg, ",")

	return in(d, data1)

}

func in(s string, g []string) bool {
	if s == "" {
		return true
	}
	for _, i := range g {
		if s == i {
			return true
		}
	}
	return false
}

//Take n and reduce it to class-c
func reduceClass(n *Node, c int) *Node {

	var nodes []*Node

	nodes = getSubnodes(n)

	nodes = rmNodes(nodes, c)

	return nodes[0]

}

func rmNodes(in []*Node, c int) (out []*Node) {

	for _, e := range in {
		if e.Generation()+1 > c {
			continue
		}
		if e.Generation()+1 == c {
			e.subnode1 = nil
			e.subnode2 = nil
			e.SetConnective("")
		}
		out = append(out, e)
	}
	return out
}

func getSubnodes(n *Node) []*Node {

	var gs func(n *Node, list []*Node) []*Node

	gs = func(n *Node, list []*Node) []*Node {

		list = append(list, n)

		if n.IsAtomic() {
			return list
		}

		if n.subnode1 != nil {
			list = gs(n.subnode1, list)
		}

		if n.subnode2 != nil {
			list = gs(n.subnode2, list)
		}

		return list
	}

	var list []*Node

	return gs(n, list)
}

func nAtomicElements(nodes []*Node) (n int) {

	n = 0

	for _, e := range nodes {
		if e.IsAtomic() {
			n++
		}
	}
	return n
}

func nAtomicElementsDistinct(nodes []*Node) int {

	var collected []string

	for _, e := range nodes {
		if e.IsAtomic() && !in(e.Formula(), collected) {
			collected = append(collected, e.Formula())
		}
	}
	return len(collected)
}

func replaceAtomic(nodes []*Node, old, repl string) []*Node {

	for _, n := range nodes {
		if n.IsAtomic() {
			if n.Formula() == old {
				n.SetFormula(repl)
			}
		}
	}
	return nodes
}

//order nodes by depth
func reorderNodes(nodes []*Node) (out []*Node) {

	d := findMaxDepth(nodes)

	for i := 0; i <= d; i++ {

		for _, j := range nodes {
			if j.Generation() == i {
				out = append(out, j)
			}
		}
	}
	return out
}

func findMaxDepth(nodes []*Node) int {

	var ds []int

	for _, n := range nodes {
		ds = append(ds, n.Generation())
	}

	sort.Ints(ds)

	return ds[len(ds)-1]
}

//check if s2 is instance of s1
func sameStructure(s1, s2 string) bool {

	ns1 := getSubnodes(Parse(s1))
	ns2 := getSubnodes(Parse(s2))

	ns1 = reorderNodes(ns1)
	ns2 = reorderNodes(ns2)

	if len(ns2) < len(ns1) {
		return false
	}

	for i := range ns1 {

		if ns1[i].IsAtomic() {
			old := ns1[i].Formula()
			repl := ns2[i].Formula()
			for _, n := range ns1 {
				if n.IsAtomic() && !n.HasFlag("c") {
					if n.Formula() == old {
						n.SetFormula(repl)
						n.SetFlag("c")
					}
				}
			}
		}
		continue

		if ns1[i].MainConnective() != ns2[i].MainConnective() {
			return false
		}
	}

	return printNodePolish(ns1[0]) == printNodePolish(ns2[0])

}
