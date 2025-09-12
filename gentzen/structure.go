package gentzen

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// check if s0 and s1 have same structure. s0 is the template. Make sure s0 and s1 are already
// appropriately normalized and are wff. If repl is non-empty, its mapping will be respected and
// also appended.
func issubstitutioninstance(s0, s1 string, repl map[string]string, exclude []string) bool {

	n0 := Parse(s0, !allowGreekUpper)
	n1 := Parse(s1, !allowGreekUpper)

	//strategy is to compare the two trees level by level and see if we
	//can find replacement subtrees for n0 so that n0 can be turned into n1.

	nodes0 := levelWalk(n0)
	nodes1 := levelWalk(n1)

	for i := range nodes0 {
		if i > len(nodes1)-1 {
			fmt.Println("insufficient depth")
			return false
		}
		for j := range nodes0[i] {
			k := 0
			c := 0
			for c = range nodes1[i] {
				if nodes1[i][c].HasFlag("ignore") {
					continue
				}
				if k == j {
					break
				}
				k++
			}
			if k != j {
				return false
			}
			k = c

			if !nodes0[i][j].IsLeafNode() {
				if nodes0[i][j].Val.connective != nodes1[i][k].Val.connective {
					return false
				}
				continue
			}

			if nodes0[i][j].Check(isAtomic) {
				if v, ok := repl[nodes0[i][j].String()]; ok {
					if v != nodes1[i][k].String() {
						return false
					}
				} else {
					repl[nodes0[i][j].String()] = nodes1[i][k].String()
				}
			}

			if nodes0[i][j].Check(isPredicate) {

				if f, ok := repl[nodes0[i][j].Val.predicateLetter+"x_1"]; ok {
					//this is for the case we are looking for an instantiation of an open formula
					//in the template formula.
					var success bool
					var arg string
					if success, arg = findArg(f, nodes1[i][k].String(), exclude); !success {
						return false
					}
					repl["const"] = arg

				}
				if _, ok := repl[nodes0[i][j].String()]; !ok {
					repl[nodes0[i][j].String()] = nodes1[i][k].String()
				}
				if repl[nodes0[i][j].String()] != nodes1[i][k].String() {
					return false
				}
			}
			for l := range nodes1[i][k].Linearize() {
				if l == 0 {
					continue
				}
				nodes1[i][k].Linearize()[l].SetFlag("ignore")
			}
			continue
		}
	}

	return true
}

// replace sentential atomic subformulas within s according to replacement pattern provided
// by repl. The replacements must be wff of sentential logic. The returned sentence is also
// a wff of sentential logic.
func replaceAtomicWith(s string, repl map[string]string) string {

	n, err := ParseStrict(s, !allowGreekUpper)
	if err != nil {
		return ""
	}

	//the remainder of processing is in sentential logic
	//but ensure we go back to the original setting
	ropl := oPL
	defer func() {
		oPL = ropl
	}()
	SetPL(false)

	nodes := sententialAtoms(n)

	for j := range nodes {
		if v, ok := repl[nodes[j].String()]; ok {
			nodes[j].ReplaceWith(Parse(v, !allowGreekUpper))
		}
	}

	return n.String()
}

// replace sentential atomic subformulas within s according to replacement pattern provided
// by repl. The replacements must be wff of sentential logic. The returned sentence is also
// a wff of sentential logic.
func replaceSubformulas(s string, repl map[string]string) string {

	n, err := ParseStrict(s, !allowGreekUpper)
	if err != nil {
		return ""
	}

	ingressFunc := func(e *Node) {
		if e.HasFlag("ignore") {
			return
		}
		if v, ok := repl[e.String()]; ok {
			e.ReplaceWith(Parse(v, !allowGreekUpper))
			for i := range e.Linearize() {
				e.Linearize()[i].SetFlag("ignore")
			}
		}
	}

	donothing := func(*Node) {}

	n.Walk(ingressFunc, donothing, donothing)

	return n.String()
}

// replace atomic subformulas within s according to replacement pattern provided
// by repl. This is for predicate logic.Repl should map predicate letters to new letters,
// term letters to term letters.
func replaceAtomicPL(s string, repl map[string]string) string {

	n, err := ParseStrict(s, !allowGreekUpper)
	if err != nil {
		return ""
	}

	nodes := n.Linearize()

	for j := range nodes {
		if !nodes[j].Check(isPredicate) {
			continue
		}
		if v, ok := repl[nodes[j].Val.predicateLetter]; ok {
			nodes[j].Val.predicateLetter = v
		}
		for k := range nodes[j].Val.term {
			if v, ok := repl[nodes[j].Val.term[k]]; ok {
				nodes[j].Val.term[k] = v
			}
		}
		setraw(nodes[j])
	}

	return n.String()
}

func normalizeQuantifiers(s string) string {

	repl := make(map[string]string)
	ingressFunc := func(e *Node) {
		if e.Check(isQuantifier) {
			repl[e.Val.variable] = "var_" + strconv.Itoa(len(repl))
			e.Val.variable = repl[e.Val.variable]
			return
		}
		if e.Check(isPredicate) {
			for i := range e.Val.term {
				if v, ok := repl[e.Val.term[i]]; ok {
					e.Val.term[i] = v
				}
			}
			setraw(e)
		}
	}

	egressFunc := func(e *Node) {
		if !e.Check(isQuantifier) {
			return
		}
		if quantifierDepth(e) == 1 {
			clear(repl)
		}
	}

	donothing := func(e *Node) {}

	n, err := ParseStrict(s, !allowGreekUpper)
	if err != nil {
		return s
	}
	n.Walk(ingressFunc, donothing, egressFunc)
	s = strings.ReplaceAll(n.String(), "var_", "x_")
	return s
}

func normalizePL(s string, l ...string) string {

	n := Parse(s, !allowGreekUpper)

	preds := make(map[string]string)
	terms := make(map[string]string)
	safifyNames(n, preds, terms)

	toNormalNames(n, l...)

	return n.String()
}

func toNormalNames(n *Node, l ...string) {

	p := "F_"
	if len(l) > 0 {
		p = l[0] + "_"
	}

	t := "t_"
	if len(l) > 1 {
		t = l[1] + "_"
	}

	pred := make(map[string]string)
	terms := make(map[string]string)

	ingressFunc := func(e *Node) {
		if e.Check(isQuantifier) {
			e.Val.variable = "x_" + strings.TrimPrefix(e.Val.variable, "var_")
			return
		}
		if !e.Check(isPredicate) {
			return
		}
		if _, ok := pred[e.Val.predicateLetter]; !ok {
			pred[e.Val.predicateLetter] = p + strings.TrimPrefix(e.Val.predicateLetter, "pr_")
		}
		e.Val.predicateLetter = pred[e.Val.predicateLetter]
		for i := range e.Val.term {
			if strings.HasPrefix(e.Val.term[i], "var_") {
				e.Val.term[i] = "x_" + strings.TrimPrefix(e.Val.term[i], "var_")
				continue
			}
			if _, ok := terms[e.Val.term[i]]; !ok {
				terms[e.Val.term[i]] = t + strings.TrimPrefix(e.Val.term[i], "co_")
			}
			e.Val.term[i] = terms[e.Val.term[i]]
		}
		setraw(e)
	}

	donothing := func(*Node) {}

	n.Walk(ingressFunc, donothing, donothing)

	return
}

// change predicate and terms/variables so that there are no clashes
// you cannot get a string out of this! repl needs to be an initialized map
// and will be appended. If you need a clean replacement map, ensure it
// manually.
func safifyNames(n *Node, preds, terms map[string]string) {

	repl := make(map[string]string)
	ingressFunc := func(e *Node) {
		if e.Check(isQuantifier) {
			repl[e.Val.variable] = "var_" + strconv.Itoa(len(repl)+1)
			e.Val.variable = repl[e.Val.variable]
			return
		}
		if e.Check(isPredicate) {
			for i := range e.Val.term {
				if v, ok := repl[e.Val.term[i]]; ok {
					e.Val.term[i] = v
				}
			}
		}
	}

	egressFunc := func(e *Node) {
		if !e.Check(isQuantifier) {
			return
		}
		if quantifierDepth(e) == 1 {
			clear(repl)
		}
	}

	donothing := func(*Node) {}

	n.Walk(ingressFunc, donothing, egressFunc)

	p := "pr_"
	t := "co_"

	ingressFunc = func(e *Node) {
		if !e.Check(isPredicate) {
			return
		}
		if _, ok := preds[e.Val.predicateLetter]; !ok {
			preds[e.Val.predicateLetter] = p + strconv.Itoa(len(preds)+1)
		}
		e.Val.predicateLetter = preds[e.Val.predicateLetter]
		for i := range e.Val.term {
			if strings.HasPrefix(e.Val.term[i], "var_") {
				continue
			}
			if _, ok := terms[e.Val.term[i]]; !ok {
				terms[e.Val.term[i]] = t + strconv.Itoa(len(terms)+1)
			}
			e.Val.term[i] = terms[e.Val.term[i]]
		}
	}

	n.Walk(ingressFunc, donothing, donothing)

}

// Check whether s2 can be got by replacing terms in s1 in the way specified by repl.
// If once is specified, replace only once for each term. Other wise all instances
// of a term specified in repl are replaced.
func isInstanceOf(s1, s2 string, repl map[string]string) bool {

	ingressFunc := func(e *Node) {
		if !e.Check(isPredicate) {
			return
		}
		for i := range e.Val.term {
			if v, ok := repl[e.Val.term[i]]; ok {
				e.Val.term[i] = v
			}
		}
		setraw(e)
	}

	donothing := func(*Node) {}

	n := Parse(s1, !allowGreekUpper)
	n.Walk(ingressFunc, donothing, donothing)
	return n.String() == s2
}

// f is a single variable open formula f(x). look for an argument c such that s is f(c). If none can be
// found, returns false and an empty string. f and s need to be normalized. In particular, the variable in f must
// be named x_1.
func findArg(f, s string, exclude []string) (ok bool, arg string) {

	var terms []string

	//we do not allow trivial satisfaction by not having
	//x_1 in f.
	found := false
	for _, e := range Parse(f, !allowGreekUpper).Linearize() {
		if !e.Check(isPredicate) {
			continue
		}
		for _, t := range e.Val.term {
			if t == "x_1" {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return false, ""
	}

	ns := Parse(s, !allowGreekUpper)

	for _, e := range ns.Linearize() {
		if !e.Check(isPredicate) {
			continue
		}
		for _, t := range e.Val.term {
			terms = append(terms, t)
		}
	}

	for _, arg = range terms {
		if isInstanceOf(f, s, map[string]string{"x_1": arg}) {
			if !slices.Contains(slices.DeleteFunc(exclude, deleteOnceFunc(exclude, arg)), arg) {
				return true, arg
			}
		}
	}

	return false, ""
}

// f is a single variable open sentences f(x). Return value is
// f(c)
func instantiate(f string, c string) string {

	n := Parse(f, !allowGreekUpper).Linearize()

	for i := range n {
		if !n[i].Check(isPredicate) {
			continue
		}
		for j := range n[i].Val.term {
			if n[i].Val.term[j] == "x_1" {
				n[i].Val.term[j] = c
			}
		}
		setraw(n[i])
	}
	return n[0].String()

}

func matchPattern(series, canonical []sequent, reverseOrder bool, exclude []string) bool {

	if len(series) != len(canonical) {
		return false
	}

	repl := make(map[string]string)

	//deal with back formula side first.
	var j int
	for i := range canonical {
		j = i
		if reverseOrder {
			j = len(canonical) - 1 - i
		}
		if !issubstitutioninstance(canonical[j].back, series[j].back, repl, exclude) {
			fmt.Println("fail here 1")
			return false
		}
	}

	//deal with front formula side

	replSet := make(map[string][]string)

	delete(repl, "const")

	for i := range canonical {
		//figure out replacement maps for sets
		//first, remove all elements that are not in the
		//range of the replacement map
		var dfront []string
		dfront = append(dfront, series[i].front...)
		for j := range canonical[i].front {
			if Parse(canonical[i].front[j], allowGreekUpper).Check(isSet) {
				//remove elements of known sets
				if _, ok := replSet[canonical[i].front[j]]; ok {
					dfront = removeSubslice(dfront, replSet[canonical[i].front[j]])
				}
				continue
			}
			//remove non-set specified elements
			for k := range series[i].front {
				if Parse(series[i].front[k], allowGreekUpper).Check(isSet) {
					continue
				}
				fmt.Println("check front item:", series[i].front[k])
				if issubstitutioninstance(canonical[i].front[j], series[i].front[k], repl, exclude) {
					dfront = slices.DeleteFunc(dfront, deleteOnceFunc(dfront, series[i].front[k]))
				}
			}
		}
		//now we can tell the substitution map for
		//remaining set
		for j := range canonical[i].front {
			if !Parse(canonical[i].front[j], allowGreekUpper).Check(isSet) {
				continue
			}
			if _, ok := replSet[canonical[i].front[j]]; ok {
				continue
			}
			replSet[canonical[i].front[j]] = append(replSet[canonical[i].front[j]], dfront...)
			break
		}
	}

	fmt.Println(repl, replSet)

	//check if the replacements result in what we want
	canonicalSub := make([]sequent, 0, len(canonical))
	canonicalSub = append(canonicalSub, canonical...)
	for i := range canonicalSub {
		canonicalSub[i].back = replaceSubformulas(canonicalSub[i].back, repl)
		var newfront []string
		for j := range canonicalSub[i].front {
			n := Parse(canonicalSub[i].front[j], allowGreekUpper)
			if !n.Check(isSet) {
				newfront = append(newfront, replaceSubformulas(n.String(), repl))
				continue
			}
			newfront = append(newfront, replSet[n.String()]...)
		}
		canonicalSub[i].front = newfront

		if canonicalSub[i].back != series[i].back {
			fmt.Println("fail here 2")
			return false
		}

		fmt.Println("compare", canonicalSub[i], series[i])

		if !isRewriteFront(canonicalSub[i], series[i]) {
			fmt.Println("fail here 3")
			return false
		}

	}
	return true
}

func getConstants(series []sequent) []string {

	var resp []string

	ingressFunc := func(e *Node) {
		if !e.Check(isPredicate) {
			return
		}
		for _, t := range e.Val.term {
			if strings.HasPrefix(t, "x_") {
				continue
			}
			//		if slices.Contains(resp, t) {
			//			continue
			//		}
			resp = append(resp, t)
		}
	}

	donothing := func(*Node) {}

	for i := range series {

		for _, d := range series[i].front {
			e := Parse(d, allowGreekUpper)
			if e.Check(isSet) {
				continue
			}

			e.Walk(ingressFunc, donothing, donothing)

		}

		Parse(series[i].back, !allowGreekUpper).Walk(ingressFunc, donothing, donothing)

	}

	return resp
}
