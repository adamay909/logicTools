package gentzen

import (
	"math/rand"
	"strconv"
	"strings"
)

var (
	atomicE = []string{
		"P",
		"Q",
		"R",
		"S",
		"W",
		"M",
		"L",
		"G",
		"H",
		"J",
	}

	logConn = []string{
		"N",
		"K",
		"A",
		"C",
		"U",
		"X",
	}

	predicateLetters = []string{
		"F",
		"G",
		"H",
		"J",
		"E",
		"P",
		"Q",
		"S",
		"R",
		"Y",
		"W",
		"Z",
		"M",
		"L",
	}

	constantLetters = []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"h",
		"m",
		"n",
	}

	variableLetters = []string{
		"x",
		"y",
		"z",
		"t",
		"p",
	}

	pletters = strings.Join(predicateLetters, "")
)

// RandomSentence returns a randomly generated sentence of
// sentential logic with maximum class of maxh and
// at most maxatomic atomic sentences. Use -1 if you don't
// want any caps on the class or number of atomic sentences.
func RandomSentence(maxh, maxatomic int) string {

	return genRand(maxh, maxatomic)

}

// maxh is maximum height of the syntax tree of generated sentence
// corresponds to class in the text
// maxatomic is the maximum number of atomic sentence to use (if less than 11, the sentences are distinct upper case letters; if more, the sentences are lowercase s followed by an arabic numeral subscript)
func genRand(maxh, maxatomic int) string {

	if oPL {
		return "Predicate Logic Not Implemented"
	}

	generateAtomicSentences := false

	if maxatomic > len(atomicE) {
		generateAtomicSentences = true
	}

	s := make([]string, 1000)

	lastConn := 4

	if !oCOND {
		lastConn = 3
	}

	newNode := logConn[rand.Intn(lastConn)]

	s = append(s, newNode)

	constantCount := 1

	openNode := 0

	switch newNode {

	case "K", "A", "C":
		openNode = 2

	default:
		openNode = 1
	}

	height := 1

	openRight := height

	for openNode > 0 {

		if height != maxh && rand.Intn(2) == 1 {

			newNode = logConn[rand.Intn(lastConn)]

			constantCount++

			height++

			if newNode != "N" {

				openNode++

				openRight++

			}
		} else {

			if !generateAtomicSentences {
				newNode = atomicE[rand.Intn(maxatomic)]
			} else {
				newNode = "s_" + strconv.Itoa(rand.Intn(maxatomic))
			}
			openNode--

			height = openRight

			openRight = height

		}

		s = append(s, newNode)

	}

	for openNode > 0 {

		if !generateAtomicSentences {
			newNode = atomicE[rand.Intn(maxatomic)]
		} else {
			newNode = "s_" + strconv.Itoa(rand.Intn(maxatomic))
		}
		newNode = atomicE[rand.Intn(maxatomic)]

		openNode--

		s = append(s, newNode)

	}

	return strings.Join(s, "")

}

// maxh is maximum height of the syntax tree of generated sentence
// corresponds to class in the text
func genRandPL() string {

	return ""

}
