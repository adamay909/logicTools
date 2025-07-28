package gentzen

import (
	"strings"
)

// ConvertToDNF returns a dnf formula equivalent to s.
func ConvertToDNF(s string) string {

	s = strings.ReplaceAll(s, "C", "AN")

	tks, err := tokenize(s, !allowGreekUpper, !allowSpecial)

	if err != nil {
		return s
	}

	var repl tokenStr

	var subs []tokenStr

	replaced := true

	for replaced {

		repl = nil

		subs = nil

		replaced = false

		for i := 0; i < len(tks)-2; i++ {

			//NNp => p
			if tks[i].tokenType == tNeg && tks[i+1].tokenType == tNeg {

				tks = tks.replaceFormulaAt(i, tks[i+1:].subFormulas()[0])

				replaced = true

				break
			}

			//NKpq => ANpNq
			if tks[i].tokenType == tNeg && tks[i+1].tokenType == tConj {

				subs = tks[i+1:].subFormulas()

				repl = subs[0].negate().disjoin(subs[1].negate())

				tks = tks.replaceFormulaAt(i, repl)

				replaced = true

				break
			}

			//NApq => KNpNq
			if tks[i].tokenType == tNeg && tks[i+1].tokenType == tDisj {

				subs = tks[i+1:].subFormulas()

				repl = subs[0].negate().conjoin(subs[1].negate())

				tks = tks.replaceFormulaAt(i, repl)

				replaced = true

				break

			}

			if tks[i].isBinary() {
				subs = tks[i:].subFormulas()
			}

			// KpAqr =>KAqrp
			if tks[i].tokenType == tConj {

				if subs[0][0].tokenType != tDisj && subs[1][0].tokenType == tDisj {

					repl = subs[1].conjoin(subs[0])

					tks = tks.replaceFormulaAt(i, repl)

					replaced = true

					break

				}
			}

			//ApKqr => AKqrp
			if tks[i].tokenType == tDisj {

				if subs[0][0].tokenType != tConj && subs[1][0].tokenType == tConj {

					repl = subs[1].disjoin(subs[0])

					tks = tks.replaceFormulaAt(i, repl)

					replaced = true

					break

				}
			}
			//KApqr => AKprKqr
			if tks[i].tokenType == tConj && tks[i+1].tokenType == tDisj {

				subs2 := subs[0].subFormulas()

				repl = subs2[0].conjoin(subs[1]).disjoin(subs2[1].conjoin(subs[1]))

				tks = tks.replaceFormulaAt(i, repl)

				replaced = true

				break

			}

		}
	}

	return tks.String()
}

func isTautologyDNF(s string) bool {

	d := ConvertToDNF("N" + s)

	tks, _ := tokenize(d, !allowGreekUpper, !allowSpecial)

	var idx int

	var conjuncts, disjuncts, subs []tokenStr

	idx = tks.index(tDisj)

	if idx == -1 {
		disjuncts = append(disjuncts, tks)

	} else {
		for pos := 0; pos < len(tks); pos++ {
			idx = tks[pos:].index(tDisj)
			if idx == -1 {
				break
			}
			f := tks.wffAt(pos + idx)
			subs = f.subFormulas()
			for _, str := range subs {
				if str[0].tokenType != tDisj {
					disjuncts = append(disjuncts, str)
				}
			}
			pos = pos + idx
		}
	}

	var f tokenStr
	var inconsistent bool

	for _, disjunct := range disjuncts {
		conjuncts = nil

		idx = disjunct.index(tConj)
		if idx == -1 {
			conjuncts = append(conjuncts, disjunct)
		} else {
			for pos := 0; pos < len(disjunct); pos++ {

				if disjunct[pos].tokenType != tConj {
					continue
				}

				f = disjunct.wffAt(pos)

				subs = f.subFormulas()

				for _, sub := range subs {
					if sub[0].tokenType != tConj {
						conjuncts = append(conjuncts, sub)
					}
				}
			}
		}

		inconsistent = false
		for _, e := range conjuncts {
			if e[0].tokenType == tAtomicSentence {
				for _, c := range conjuncts {
					if c.String() == "N"+e[0].str {
						inconsistent = true
					}
				}
			}
		}

		if !inconsistent {
			return false
		}
	}

	return true
}
