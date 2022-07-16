package main

import "github.com/adamay909/logicTools/gentzen"

func generateRandomExercise() {

	var p1, p2, c *gentzen.Node
	s := gentzen.Parse(randomTautology())

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

	problem := `Derive from &ensp;` + seq1 + ` to &ensp;` + seq3

	if p2 != nil {

		problem = `Derive from ` + seq1 + `&ensp; and &ensp; ` + seq2 + `&ensp; to &ensp; ` + seq3

	}

	dsp.clear()

	dsp.setTitle(problem)

	display()

	backToNormal()

	stopInput()

}

var tautologies = []string{`>PP`}

func randomTautology() string {

	var s string

	for s = gentzen.RandomSentence(3, 8); ; s = gentzen.RandomSentence(3, 8) {

		if !gentzen.IsTautology(s) {
			continue
		}

		if !gentzen.Parse(s).IsConditional() {
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
