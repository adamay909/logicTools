package main

import (
	"encoding/json"
	"strconv"

	"github.com/adamay909/logicTools/gentzen"
)

var tautologies = []string{`>PP`}
var exercises []string
var expos int

func nextProblem() {

	if len(exercises) == 0 {
		genNewExercise()
		exercises = append(exercises, marshalJson())
		displayExercise(exercises[expos])
		return
	}

	exercises[expos] = string(marshalJson())
	expos++

	if expos == len(exercises) {
		genNewExercise()
		exercises = append(exercises, marshalJson())
	}

	displayExercise(exercises[expos])

}

func prevProblem() {
	exercises[expos] = marshalJson()
	expos--
	if expos < 0 {
		expos = 0
	}

	displayExercise(exercises[expos])
}

func displayExercise(s string) {

	saveHistory()
	json.Unmarshal([]byte(s), dsp)
	dsp.xpos, dsp.ypos = 0, 0
	dsp.overhang = false
	dsp.modifier = ""

	if dsp.SystemPL != oPL {
		togglePL()
	}

	if dsp.Theorems != oTHM {
		toggleTheorems()
	}

	backToNormal()

	show("controls2")

	setTextByID("controls2", `<button class="controls" id="prevExercise">Previous</button><button class="controls" id="nextExercise">Next</button>`)

	hide("messages")

	stopInput()
}

func endRandomExercise() {

	oExercises = false

	dsp.clear()

	display()

	hide("controls2")

	setTextByID("controls2", "")
}

func genNewExercise() {

	if oPL {
		togglePL()
	}

	oExercises = true
	var p1, p2, c *gentzen.Node
	s := gentzen.Parse(genRandomTautology())

	p1 = gentzen.Parse(s.Child1Must().String())
	if p1.IsConjunction() {
		p2 = gentzen.Parse(p1.Child2Must().String())
		p1 = gentzen.Parse(p1.Child1Must().String())
	}
	c = gentzen.Parse(s.Child2Must().String())

	var seq1, seq2, seq3 string

	seq1 = textOf(`\G`) + " " + textOf(`|-`) + " " + p1.StringPlain()

	seq3 = textOf(`\G`) + " " + textOf(`|-`) + " " + c.StringPlain()

	if p2 != nil {

		seq2 = textOf(`\D`) + " " + textOf(`|-`) + " " + p2.StringPlain()

		seq3 = textOf(`\G`) + ", " + textOf(`\D`) + " " + textOf(`|-`) + "  " + c.StringPlain()

	}

	problem := `Q` + strconv.Itoa(expos+1) + `. Derive from &ensp;` + seq1 + ` to &ensp;` + seq3

	if p2 != nil {

		problem = `Q` + strconv.Itoa(expos+1) + `. Derive from ` + seq1 + `&ensp; and &ensp; ` + seq2 + `&ensp; to &ensp; ` + seq3

	}

	dsp.clear()

	dsp.setTitle(problem)

	display()

}

func genRandomTautology() string {

	var s string

	for s = gentzen.RandomSentence(3, 8); ; s = gentzen.RandomSentence(3, 8) {

		if !gentzen.IsTautology(s) {
			continue
		}

		sn := gentzen.Parse(s)

		if !sn.IsConditional() {
			continue
		}

		if gentzen.IsTautology(sn.Child2Must().String()) {
			continue
		}

		if sn.Child1Must().String() == sn.Child2Must().String() {
			continue
		}

		if contains(tautologies, s) {
			continue
		}

		break
	}

	tautologies = append(tautologies, s)
	return s
}

func contains(s []string, e string) bool {

	return indexOf(s, e) > -1
}

func indexOf(s []string, e string) int {
	for i := range s {
		if s[i] == e {
			return i
		}
	}
	return -1
}
