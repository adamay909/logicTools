package gentzen

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type sequent struct {
	front []string
	back  string
}

var (
	errNotSequent error = errors.New("not a sequent")
)

func mkSequent(s string, mode PrintMode) (seq sequent, err error) {

	idx := strings.Index(s, ":")

	backstr := ""
	frontstr := ""

	if idx == -1 {
		fmt.Println("1")
		err = errNotSequent
		return
	}

	if idx == len(s)-1 {
		fmt.Println("2")
		err = errNotSequent
		return
	}

	if idx == 0 {
		backstr = s[1:]
	} else {
		frontstr = s[:idx]
		backstr = s[idx+1:]
	}

	var parser func(string, bool) (*Node, error)
	if mode == O_PlainASCII || mode == O_ProofChecker {
		parser = ParseInfix
	} else {
		parser = ParseStrict
	}

	var n *Node

	n, err = parser(backstr, !allowGreekUpper)
	if err != nil {
		fmt.Println("3")
		err = errors.Join(errNotSequent, err)
		return
	}

	seq.back = n.String()

	ds := strings.Split(frontstr, ",")

	for _, d := range ds {

		if len(d) == 0 {
			continue
		}

		n, err = parser(d, true)
		if err != nil {
			fmt.Println(err)
			fmt.Println("4")
			err = errNotSequent
			return
		}

		seq.front = append(seq.front, n.String())
	}

	return
}

func (s sequent) String() string {

	w := new(strings.Builder)

	for i, d := range s.front {
		if d == "" {
			continue
		}
		w.WriteString(d)
		if i < len(s.front)-1 {
			w.WriteString(",")
		}
	}

	w.WriteString(":")

	w.WriteString(s.back)

	return w.String()
}

func (s sequent) isValid() bool {

	var permitGreekUpper = true

	if len(s.back) == 0 {
		return false
	}

	if _, err := ParseStrict(s.back, !permitGreekUpper); err != nil {
		return false
	}

	for _, d := range s.front {
		if _, err := ParseStrict(d, permitGreekUpper); err == nil {
			continue
		}
		return false
	}

	return true
}

func (s sequent) frontString() string {

	return strings.Split(s.String(), ":")[0]

}

func (s sequent) backString() string {

	return strings.Split(s.String(), ":")[1]

}

// check if s can be gotten by substitutions into template
// given substititution rules given in subst
func (s sequent) instanceOf(template sequent, subst map[string]string) bool {

	//make sure we are operating on copies of sequents
	target := s.mkCopy()
	templ := template.mkCopy()

	//	fmt.Println("replacements:", subst)

	for i := range templ.front {
		for k, v := range subst {
			templ.front[i] = strings.ReplaceAll(templ.front[i], k, v)
		}
	}
	for k, v := range subst {
		templ.back = strings.ReplaceAll(templ.back, k, v)
	}

	slices.Sort(templ.front)
	slices.Sort(target.front)

	return templ.String() == target.String()
}

func (s sequent) mkCopy() sequent {

	resp := sequent{
		back: s.back,
	}

	resp.front = append(resp.front, s.front...)

	return resp
}

func (s sequent) isReductionOf(s2 sequent) bool {

	if s.String() == s2.String() {
		return true
	}

	//any extra members of front formulas of s2 must be duplicates
	//of a member of front formulas of s:
	for i := range s.front {
		if !slices.Contains(s2.front, s.front[i]) {
			//			fmt.Println(s.front[i], " not in ", s2.front)
			return false
		}
	}
	for i := range s2.front {
		if !slices.Contains(s.front, s2.front[i]) {
			//			fmt.Println(s2.front[i], " not in ", s.front)
			return false
		}
	}
	return true
}

func (s sequent) stringF(mode PrintMode) string {

	w := new(strings.Builder)

	for i, d := range s.front {
		w.WriteString(StringF(Parse(d, allowGreekUpper), mode))
		if i < len(s.front)-1 {
			w.WriteString(`,`)
		}
	}
	frontString := w.String()

	w.Reset()

	backString := StringF(Parse(s.back, !allowGreekUpper), mode)

	if mode == O_Latex {
		w.WriteString(`\seq{`)
		w.WriteString(frontString)
		w.WriteString(`}{`)
		w.WriteString(backString)
		w.WriteString(`}`)
	}

	if mode == O_PlainText || mode == O_ProofChecker {
		w.WriteString(frontString)
		w.WriteString(`⊢`)
		w.WriteString(backString)
	}

	if mode == O_PlainASCII || mode == O_English || mode == O_Polish {
		w.WriteString(frontString)
		w.WriteString(":")
		w.WriteString(backString)
	}

	return w.String()
}
