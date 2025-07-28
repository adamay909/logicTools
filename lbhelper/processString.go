package main

import (
	"errors"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func processString(input string) (output string, err error) {

	mode := gentzen.O_Polish

	if *s {
		*seq = true
	}

	switch {

	case *ascii:
		mode = gentzen.O_PlainASCII
	case *txt:
		mode = gentzen.O_PlainText
	case *e:
		mode = gentzen.O_English
	case *polish:
		mode = gentzen.O_Polish
	case *l, *f:
		mode = gentzen.O_Latex
	}

	gentzen.SetPrettyBrackets(*pretty)
	gentzen.SetPL(*pl)
	gentzen.SetConditional(*c)

	//	gentzen.SetML(true)

	n := new(gentzen.Node)

	parser = gentzen.ParseStrict

	gu := false

	if *infix {
		parser = gentzen.ParseInfix
	}

	if !*seq {
		if n, err = parser(input, gu); err != nil {
			err = errors.New("ERROR: " + err.Error() + "\n")
			return
		}
	}

	if *normalize && !*seq {
		input = gentzen.Normalize(input)
	}

	switch {

	case *tt:
		if *pl {
			err = errors.New("ERROR: truth table can only produced for sentential logic")
			return "", err
		}

		n, _ = parser(input, gu)
		input = n.String()

		tt, err := gentzen.GenerateTruthTable(input)
		if err != nil {
			err = errors.New("ERROR: " + err.Error() + "\n")
			return "", err
		}

		output = tt.PrintTruthTable(mode, false)

	case *ttn:
		if *pl {
			err = errors.New("ERROR: truth table can only produced for sentential logic")
			return "", err
		}
		n, _ = parser(input, gu)
		input = n.String()

		tt, err := gentzen.GenerateTruthTableNarrow(input)
		if err != nil {
			err = errors.New("ERROR: " + err.Error() + "\n")
			return "", err
		}

		output = tt.PrintTruthTable(mode, false)

	case *stf:
		n, _ = parser(input, gu)
		output = n.SyntaxTree()

	case *sts:
		n, _ = parser(input, gu)
		output = n.SyntaxTreeSimple()

	case *seq:
		output, err = sequentString(input, mode)
		if err != nil {
			err = errors.New("ERROR: " + err.Error() + "\n")
			return
		}
		output = encloseMath(output, *m)

	case *tableau:
		if *pl {
			err = errors.New("ERROR: tableau not supported for predicate logic" + "\n")
			return
		}
		output = gentzen.PrintSemanticTableau(input)

	case *proof:
		if *pl {
			err = errors.New("ERROR: proof outline not supported for predicate logic" + "\n")
			return "", err
		}
		output = gentzen.PrintProofOutline(input, mode)
		if !gentzen.IsTautology(input) {
			output = output + "\n" + input + " is not a theorem.\n"
		}

	default:
		n, _ = parser(input, gu)
		if mode == gentzen.O_Latex {
			output = encloseMath(n.StringF(mode), *m)
		} else {
			output = n.StringF(mode)
		}

	}

	return

}

func encloseMath(s string, enclose bool) string {
	if !enclose {
		return s
	}
	return `\p{` + s + `}`
}

func sequentString(s string, mode gentzen.PrintMode) (r string, err error) {

	gu := true

	parts := strings.Split(s, ":")

	if len(parts) != 2 {
		err = errors.New("BAD Sequent")
		return
	}

	r = `\seq{`
	datum := parts[0]
	succedent := parts[1]
	d := strings.Split(datum, `,`)

	for _, e := range d {
		if len(e) == 0 {
			continue
		}

		var n *gentzen.Node

		n, err = parser(e, gu)
		if err != nil {
			r = "**BAD Datum"
			return
		}
		r = r + n.StringF(mode) + `,`
	}
	r = strings.TrimRight(r, `,`)

	r = r + `}{`
	n, err := parser(succedent, !gu)
	if err != nil {
		r = "**BAD Succedent"
		return
	}
	r = r + n.StringF(mode) + `}`

	return
}

func convertGreek(s string) string {

	var o string

	r := []rune(s)

	for len(r) > 0 {
		if string(r[0]) != "/" {
			o = o + string(r[0])
			if len(r) == 1 {
				break
			}
			r = r[1:]
			continue
		}
		if len(r) < 2 {
			return "ERROR: malformed1"
		}

		found := false
		for _, g := range greekLetters {
			if string(r[:2]) == g[0] {
				o = o + g[2]
				found = true
				break
			}
		}
		if !found {
			return "ERROR: malformed2"
		}
		if len(r) == 2 {
			break
		}
		r = r[2:]
	}
	return o
}

var greekLetters = [][3]string{
	[3]string{`/G`, `\Gamma`, "\u0393"},
	[3]string{`/D`, `\Delta`, "\u0394"},
	[3]string{`/T`, `\Theta`, "\u0398"},
	[3]string{`/L`, `\Lambda`, "\u039b"},
	[3]string{`/X`, `\Xi`, "\u039e"},
	[3]string{`/P`, `\Pi`, "\u03a0"},
	[3]string{`/R`, `\Rho`, "\u03a1"},
	[3]string{`/S`, `\Sigma`, "\u03a3"},
	[3]string{`/U`, `\Upsilon`, "\u03a5"},
	[3]string{`/F`, `\Phi`, "\u03a6"},
	[3]string{`/Q`, `\Psi`, "\u03a8"},
	[3]string{`/W`, `\Omega`, "\u03a9"},
	[3]string{`/a`, `\alpha`, "\u03b1"},
	[3]string{`/b`, `\beta`, "\u03b2"},
	[3]string{`/g`, `\gamma`, "\u03b3"},
	[3]string{`/d`, `\delta`, "\u03b4"},
	[3]string{`/e`, `\epsilon`, "\u03b5"},
	[3]string{`/z`, `\zeta`, "\u03b6"},
	[3]string{`/h`, `\eta`, "\u03b7"},
	[3]string{`/t`, `\theta`, "\u03b8"},
	[3]string{`/i`, `\iota`, "\u03b9"},
	[3]string{`/k`, `\kappa`, "\u03ba"},
	[3]string{`/l`, `\lambda`, "\u03bb"},
	[3]string{`/m`, `\mu`, "\u03bc"},
	[3]string{`/n`, `\nu`, "\u03bd"},
	[3]string{`/x`, `\xi`, "\u03be"},
	[3]string{`/o`, `\omicron`, "\u03bf"},
	[3]string{`/p`, `\pi`, "\u03c0"},
	[3]string{`/r`, `\rho`, "\u03c1"},
	[3]string{`/s`, `\sigma`, "\u03c3"},
	[3]string{`/y`, `\tau`, "\u03c4"},
	[3]string{`/u`, `\upsilon`, "\u03c5"},
	[3]string{`/f`, `\phi`, "\u03c6"},
	[3]string{`/c`, `\chi`, "\u03c7"},
	[3]string{`/q`, `\psi`, "\u03c8"},
	[3]string{`/w`, `\omega`, "\u03c9"},
}
