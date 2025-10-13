package main

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func vimhmain() {

	var outputDest *os.File
	outputDest = os.Stdout

	scanner := bufio.NewScanner(os.Stdin)

	for waitForInput := scanner.Scan(); waitForInput; waitForInput = scanner.Scan() {

		input := scanner.Text()
		output, err := vimhprocessString(input)

		if err != nil {
			outputDest.WriteString("")
		} else {
			outputDest.WriteString(output)
		}

		outputDest.WriteString("\n")
	}

}

func vimhprocessString(s string) (resp string, err error) {

	if ok, sep := isSequent(s); ok {
		return vimhprocessSequent(s, sep)
	}

	for _, parser = range []func(string, bool) (*gentzen.Node, error){gentzen.ParseStrict, gentzen.ParseInfix} {
		for _, pl := range []bool{false, true} {
			gentzen.SetPL(pl)
			n, err := parser(s, true)
			if err != nil {
				continue
			}
			resp = `\p{` + gentzen.StringF(n, gentzen.O_Latex) + `}`
			return resp, err
		}
	}
	err = errors.New("failed")
	return
}

func vimhprocessSequent(s string, sep string) (resp string, err error) {

	fail := false

	parts := strings.Split(s, sep)
	front := strings.Split(parts[0], ",")
	back := parts[1]

	var respback string

	w := new(strings.Builder)

	for _, parser = range []func(string, bool) (*gentzen.Node, error){gentzen.ParseStrict, gentzen.ParseInfix} {
		for _, pl := range []bool{false, true} {
			w.Reset()
			var respfront = make([]string, 0, len(front))
			fail = false
			gentzen.SetPL(pl)
			for i := range front {
				if len(front[i]) == 0 {
					continue
				}
				n, err := parser(front[i], true)
				if err != nil {
					fail = true
					break
				}
				respfront = append(respfront, gentzen.StringF(n, gentzen.O_Latex))
			}
			if fail {
				continue
			}
			n, err := parser(back, false)
			if err != nil {
				fail = true
				continue
			}
			respback = gentzen.StringF(n, gentzen.O_Latex)

			w.WriteString(`\p{`)
			w.WriteString(strings.Join(respfront, ","))
			switch sep {
			case ":", "|-":
				w.WriteString(`\lproves `)
			case "|=":
				w.WriteString(`\lentails `)
			}
			w.WriteString(respback)
			w.WriteString(`}`)
			return w.String(), nil
		}
	}
	return "", errors.New("fail")
}

func isSequent(s string) (bool, string) {
	var sep = []string{":", "|-", "|="}

	for i := range sep {
		if len(strings.Split(s, sep[i])) == 2 {
			return true, sep[i]
		}
	}
	return false, ""
}
