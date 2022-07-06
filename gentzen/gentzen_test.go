package gentzen

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTree(t *testing.T) {
	var r string
	rand.Seed(int64(time.Now().Second()))

	SetPL(false)
	SetStandardPolish(true)
	SetConditional(false)
	for i := 0; i < 10; i++ {

		s := RandomSentence(3, 5)

		if s == "P" {
			i--
			continue
		}

		r = r + LatexTreeSimple(Parse(s)) + "\n\n"

	}
	FileWrite("results.tex", r)

}

func _TestRandom(t *testing.T) {

	rand.Seed(int64(time.Now().Second()))

	SetPL(false)
	SetStandardPolish(true)

	for i := 0; i < 10; i++ {

		s := RandomSentence(3, 5)

		fmt.Println(s)

	}
}
