package main

func slicesIndex[c comparable](s []c, e c) int {

	for i := range s {
		if s[i] == e {
			return i
		}
	}

	return -1
}

func slicesContains[c comparable](s []c, e c) bool {

	return slicesIndex(s, e) > -1

}

func slicesEqual[c comparable](s1, s2 []c) bool {

	if len(s1) != len(s2) {
		return false
	}

	for _, e := range s1 {
		if !slicesContains(s2, e) {
			return false
		}
	}

	for _, e := range s2 {
		if !slicesContains(s1, e) {
			return false
		}
	}

	return true
}

func slicesRemove[c comparable](s []c, e c) (r []c) {

	idx := slicesIndex(s, e)

	if idx == -1 {
		return s
	}

	r = append(r, s[:idx]...)

	if idx < len(s)-1 {

		r = append(r, s[idx+1:]...)

	}

	return r
}

// check if s1 is proper super set of s2. Repeat elements count as distinct.
func slicesSupset[c comparable](s1, s2 []c) bool {

	if len(s1) <= len(s2) {
		return false
	}

	for _, e := range s2 {
		if !slicesContains(s1, e) {
			return false
		}
	}
	return true
}

func slicesCleanDuplicates[t comparable](old []t) (resp []t) {

	if len(old) == 0 {
		return old
	}

	resp = append(resp, old[0])

	for i := 1; i < len(old); i++ {
		if slicesContains(resp, old[i]) {
			continue
		}
		resp = append(resp, old[i])
	}

	return
}

func insertString(s1, s2 string, n int) string {

	if n == 0 {
		return s2 + s1
	}

	o := []rune(s1)

	o1 := o[:n]

	o2 := o[n:]

	ins := []rune(s2)

	var r []rune

	r = append(r, o1...)

	r = append(r, ins...)

	r = append(r, o2...)

	return string(r)

}

func rmString(s1, s2 string) string {

	o := []rune(s1)

	n := []rune(s2)

	r := o

	if len(o) < len(n) {
		return string(o)
	}
	for i := range o {
		match := false
		if o[i] == n[0] {
			if len(o[i:]) < len(n) {
				break
			}
			for j := range n {
				if o[i+j] != n[j] {
					match = false
					break
				}
				match = true
			}
			if match {
				r = append(o[:i], o[i+len(n):]...)
			}
		}
	}
	return string(r)

}
