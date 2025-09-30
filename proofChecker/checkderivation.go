package main

import (
	"fmt"
	"strings"

	"github.com/adamay909/logicTools/gentzen"
)

func checkDerivation() {
	cleanupEditorWindow()
	arglines := dsp.editor.GetArglines()
	gentzen.ClearLog()
	makeVisible("messages")
	for i := range arglines {
		fmt.Println(arglines[i])
	}
	showSuccess(gentzen.CheckDerivation(arglines, gentzen.O_ProofChecker, oOffset))
	return
}

func printMessage(s string, cleanup bool) {

	if cleanup {
		s = strings.ReplaceAll(s, `/`, `\`)

		for _, e := range allBindings {
			s = strings.ReplaceAll(s, e[0], e[2])
		}
	}

	l := strings.Split(s, "\n")
	for i := range l {
		l[i] = `<p class="message">` + l[i] + "</p>\n"
	}
	s = strings.Join(l, "\n")

	setTextByID("messages", s)
}

func showSuccess(err error) {
	s := domDocument.Call("querySelector", "#editorWindow").Get("classList")
	if err != nil {
		printMessage(err.Error(), true)
		s.Call("add", "fail")
		return
	}
	s.Call("add", "success")
	printMessage("OK", true)
}
