package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/adamay909/logicTools/fileops"
	"github.com/adamay909/logicTools/gentzen"
)

var outputDest *os.File

func main() {

	var (
		dest = flag.String("dest", "", "output destination; default is stdout")

		random = flag.Bool("random", false, "generate n random sentences with m atomic sentences and at most d connectives")

		n = flag.Int("n", 1, "number of formulas to generate")

		maxClass = flag.Int("maxClass", -1, "max syntax tree height (class) of generated sentence. Nagative number means no limit")

		maxAtomic = flag.Int("maxAtomic", 10, "max number of distinct atomic sentences to use")

		tautology = flag.Bool("tautology", false, "extract tautologies from f")

		malform = flag.Bool("malform", false, "like random but about half are malformed. Output is LaTeX.")

		withConditional = flag.Bool("withC", true, "with conditional")

		unique = flag.Bool("uniqueS", false, "treat sentences with same structure as the same")

		normalize = flag.Bool("normalize", false, "normalize sentence letters, etc.")
	)

	flag.Parse()
	gentzen.SetPL(false)
	gentzen.SetStandardPolish(true)
	gentzen.SetConditional(*withConditional)

	if *dest == "" {

		outputDest = os.Stdout

	} else {

		*dest = filepath.Clean(*dest)

		outputDest = fileops.CreateFile(*dest)
	}

	defer outputDest.Close()

	switch {
	case *random:
		generateRandomSentence(*n, *maxClass, *maxAtomic, *unique, *normalize)

	case *tautology:
		generateTautology(*n, *maxClass, *maxAtomic, *unique, *normalize)

	case *malform:
		generateMalform(*n, *maxClass, *maxAtomic)

	default:
		flag.PrintDefaults()
	}
	return
}
