package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamay909/logicTools/fileops"
	"github.com/adamay909/logicTools/gentzen"
	"golang.org/x/term"
)

var parser func(string, bool) (*gentzen.Node, error)

func main() {

	if *cmpl || *standalone {

		processFile()

		return

	}

	var outputDest *os.File

	if *dest == "" {

		outputDest = os.Stdout

	} else {

		*dest = filepath.Clean(*dest)

		outputDest = fileops.CreateFile(*dest)
	}

	defer outputDest.Close()

	input := flag.Arg(0)

	terminalOutput := term.IsTerminal(int(os.Stdout.Fd()))

	if input != "" {

		output, err := processString(input)

		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		outputDest.WriteString(output)
		if terminalOutput {
			outputDest.WriteString("\n")
		}

		return

	}

	terminalInput := term.IsTerminal(int(os.Stdin.Fd()))

	scanner := bufio.NewScanner(os.Stdin)

	if terminalInput {
		fmt.Print("Enter formula: ")
	}

	for waitForInput := scanner.Scan(); waitForInput; waitForInput = scanner.Scan() {

		input = scanner.Text()

		if len(input) == 0 {
			break
		}

		output, err := processString(input)

		if err != nil {
			outputDest.WriteString(err.Error() + "\n")
			if outputDest != os.Stdout {
				os.Stderr.WriteString(err.Error() + "\n")
			}
		} else {

			outputDest.WriteString(output)
			outputDest.WriteString("\n")
		}

		if terminalInput {
			fmt.Print("Enter formula: ")
		}
	}

	return

}
