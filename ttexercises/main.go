package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamay909/logicTools/fileops"
	"github.com/adamay909/logicTools/gentzen"
	"golang.org/x/term"
)

type ttable struct {
	gentzen.TruthTable
	hide    [][]bool
	mistake [][]bool
}

var (
	mistakes    = flag.Bool("errors", false, "insert errors")
	val         = flag.Bool("val", false, "include errors in interpretations/also hide interpretations")
	clearCol    = flag.Bool("delC", false, "clear column; -n columns")
	clearRow    = flag.Bool("delR", false, "clear row -n rows")
	clearRandom = flag.Bool("delXY", false, "clear up to n rows or columns")
	clearCell   = flag.Bool("del", false, "clear cell -n number of cells")
	number      = flag.Int("n", 0, "number")
	empty       = flag.Bool("empty", false, "generate truth table with just columns and valuations")
	correct     = flag.Bool("correct", false, "generate full truth table")
	dest        = flag.String("dest", "", "output destination (default is Stdout)")
	inputFile   = flag.String("inputfile", "", "read input from the named file")
	useNarrow   = flag.Bool("narrow", false, "use narrow format truth table")
)

func main() {

	flag.Parse()

	gentzen.SetStandardPolish(true)

	gentzen.SetPL(false)

	var outputDest *os.File

	if *dest == "" {

		outputDest = os.Stdout

	} else {

		*dest = filepath.Clean(*dest)

		outputDest = fileops.CreateFile(*dest)
	}

	defer outputDest.Close()

	input := flag.Arg(0)

	//Case 1: formula specified on command line
	if input != "" {

		if *inputFile != "" {
			os.Stderr.WriteString("You specified an input file but also supplied a formula. You can only do one of them. \n Exiting.\n")
			os.Exit(1)
		}

		output, err := processFormula(input, *useNarrow)

		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		outputDest.WriteString(output)
		return
	}

	//Case 2: inputfile containing formulas specfied
	if *inputFile != "" {

		*inputFile = filepath.Clean(*inputFile)

		file, err := os.Open(*inputFile)

		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {

			output, err := processFormula(scanner.Text(), *useNarrow)

			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}

			outputDest.WriteString(output)
			outputDest.WriteString("\n")

		}
		return

	}

	//Case 3: take input from standard in. We allow pipes
	terminal := term.IsTerminal(int(os.Stdin.Fd()))

	scanner := bufio.NewScanner(os.Stdin)

	if terminal {
		fmt.Print("Enter formula: ")
	}

	for waitForInput := scanner.Scan(); waitForInput; waitForInput = scanner.Scan() {

		input = scanner.Text()

		if len(input) == 0 {
			break
		}

		output, err := processFormula(input, *useNarrow)

		if err != nil {
			outputDest.WriteString(err.Error())
			if outputDest != os.Stdout {
				os.Stderr.WriteString(err.Error())
			}
		} else {

			outputDest.WriteString(output)

			outputDest.WriteString("\n")
		}

		if terminal {
			fmt.Print("Enter formula: ")
		}
	}

}
