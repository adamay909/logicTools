package main

import (
	"math/rand"

	"github.com/adamay909/logicTools/gentzen"
)

func generateMalform(n int, maxh, maxatomic int) {

	var s string

	for i := 0; i < n; {

		s = gentzen.RandomSentence(maxh, rand.Intn(maxatomic)+1)

		if rand.Intn(2) == 1 {

			outputDest.WriteString(gentzen.Parse(s, false).StringF(gentzen.O_Latex) + "\n")

			i++

			continue
		}

		outputDest.WriteString(gentzen.Malform(s) + "\n")

		i++

	}
}
