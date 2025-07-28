package main

import (
	"math/rand"

	"github.com/adamay909/logicTools/gentzen"
)

func generateRandomSentence(n int, maxh, maxatomic int, uniqueStructure bool, normalize bool) {

	var s, t string
	var ok bool

	present := make(map[string]bool, n)

	for i := 0; i < n; {

		s = gentzen.RandomSentence(maxh, rand.Intn(maxatomic)+1)

		if !uniqueStructure {
			if _, ok = present[s]; ok {
				continue
			}
			present[s] = true

		} else {

			t = gentzen.Normalize(s)
			if _, ok = present[t]; ok {
				continue
			}
			present[t] = true
		}

		i++

		if !normalize {
			outputDest.WriteString(s)
		} else {
			outputDest.WriteString(t)
		}

		outputDest.WriteString("\n")
	}
}
