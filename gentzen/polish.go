package gentzen

// IsWff returns whether s is a wff.
func IsWff(s string) bool {

	t, err := tokenize(s, !allowGreekUpper, !allowSpecial)

	if err != nil {
		return false
	}

	tk := t.wffAt(0)

	if len(tk) == 0 {
		return false
	}

	return s == tk.String()

}

// IsTautology returns whether s is a tautology.
func IsTautology(s string) bool {

	return isTautologyTableaux(s)

}

// Class returns the (minimum) class number of s.
func Class(s string) int {

	tk, err := tokenize(s, !allowGreekUpper, !allowSpecial)

	if err != nil {
		panic(s + " is not well formed")
	}

	if !tk.isWff() {
		panic(s + " is not well formed")
	}

	class := 0
	subclass := 0

	for i := len(tk) - 1; i > -1; i-- {

		if tk[i].tokenType != tAtomicSentence {
			continue
		}
		subclass = -1
		for j := i; j > -1; j-- {
			w := tk.wffAt(j)

			if j+len(w) > i {

				subclass++
			}
		}
		if subclass > class {
			class = subclass
		}
	}

	return class
}

// Normalize normalizes s by assigning subscripted sentence letters to
// each atomic sentence.
func Normalize(s string) string {

	if oPL {
		return s
	}

	tk, err := tokenize(s, !allowGreekUpper, !allowSpecial)

	if err != nil {
		panic(s + " is not wff")
	}

	tk.normalize()

	return tk.String()
}

// SubSentences returns the subsentences of s. If s is atomic,
// the returned slice has length 0.
func SubSentences(s string) []string {

	tk, err := tokenize(s, !allowGreekUpper, !allowSpecial)

	if err != nil {
		panic(s + " is not wff")
	}

	resp := make([]string, 0, 2)

	for _, e := range tk.subFormulas() {

		resp = append(resp, e.String())

	}

	return resp
}

// Equivalent returns whether s1 and s2 are logically equivalent.
func Equivalent(s1, s2 string) bool {

	return IsTautology("K" + "C" + s1 + s2 + "C" + s2 + s1)

}
