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

func genRand(m, d int, fixed bool) string {

	if m > len(atomicE) {
		m = len(atomicE)
	}

	cand := generateCandidates(m, fixed)

	var s, sNew, sOld string

	s = "P"
	for Parse(s).ConnectiveCount() < d {
		sNew = ""
		sOld = s
		for _, c := range s {
			sNew = sNew + replace(string(c), cand)
		}

		if s == sNew {
			break
		}
		s = sNew
	}
	if Parse(s).ConnectiveCount() > d {
		s = sOld
	}
	return s

}

func generateCandidates(m int, fixed bool) []string {

	var r []string
	var perm []int

	if fixed {
		for i := 0; i < len(atomicE); i++ {
			perm = append(perm, i)
		}
	} else {
		perm = rand.Perm(len(atomicE))
	}
	for _, n1 := range perm[:m] {
		s1 := atomicE[n1]

		for _, n2 := range perm[:m] {
			s2 := atomicE[n2]
			for _, c := range connectivesSL {

				if c[0] == lneg {
					r = append(r, string(c[0])+s1)
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
