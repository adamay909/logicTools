package gentzen

import (
	"math/rand"
)

var (
	atomicE = []string{
		"P",
		"Q",
		"R",
		"S",
		"T",
		"W",
		"Y",
		"Z",
		"M",
		"L",
	}
)

func genRand(m, d int) string {

	if m > len(atomicE) {
		m = len(atomicE)
	}

	cand := generateCandidates(m)

	var s, s2 string

	s = "P"

	for Parse(s).ConnectiveCount() < (d / 2) {
		s2 = ""
		for _, c := range s {
			s2 = s2 + replace(string(c), cand)
		}

		if s == s2 {
			break
		}
		s = s2
	}

	return s

}

func generateCandidates(m int) []string {

	var r []string

	for _, s1 := range atomicE[:m] {

		r = append(r, string(lneg)+s1)
		for _, s2 := range atomicE[:m] {

			for _, c := range connectivesSL {

				if c[0] == lneg {
					continue
				}
				r = append(r, string(c[0])+s1+s2)
			}
		}
	}
	return r
}

func replace(s string, cand []string) string {

	if isConnective(s) {
		return s
	}
	d := rand.Intn(2)

	if d == 0 {
		return s
	}

	return cand[rand.Intn(len(cand))]

}
