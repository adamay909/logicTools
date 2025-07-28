package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	pl         = flag.Bool("p", false, "turn on predicate logic processing")
	c          = flag.Bool("c", true, "use conditional")
	f          = flag.Bool("f", false, "print Latex of string (deprecated. use -l instead")
	l          = flag.Bool("l", true, "output formatted in Latex (default)")
	ascii      = flag.Bool("ascii", false, "output formatted as you would type into proofchecker")
	txt        = flag.Bool("txt", false, "output formatted as infix with unicode logical symbols")
	e          = flag.Bool("e", false, "print formula as English (for LaTeX)")
	seq        = flag.Bool("seq", false, "interpret input as a sequent; turnstile represented by colon")
	s          = flag.Bool("s", false, "interpret input as a sequent (deprecated. use -seq instead")
	tt         = flag.Bool("tt", false, "print truth table")
	ttn        = flag.Bool("ttn", false, "print narrow truth table")
	stf        = flag.Bool("stf", false, "print full syntax tree")
	sts        = flag.Bool("sts", false, "print simple syntax tree")
	cmpl       = flag.Bool("compile", false, "compile latex fragment")
	ans        = flag.Bool("answer", false, "compile with answers")
	standalone = flag.Bool("standalone", false, "make fragment into a standalone LaTeX document")
	m          = flag.Bool("m", true, "enclose output in math mode")
	letter     = flag.Bool("letterpaper", false, "use letter size paper for compile")
	pretty     = flag.Bool("pretty", true, "prettify brackets")
	fontsize   = flag.String("fontsize", "10pt", "fontsize in pt")
	resize     = flag.Bool("resize", false, "resize to fit page")
	polish     = flag.Bool("polish", false, "print Polish notation of string")
	infix      = flag.Bool("infix", false, "input is in infix notation (formula typed as in the proofchecker UI")
	normalize  = flag.Bool("normalize", false, "normalize output string")
	tableau    = flag.Bool("tableau", false, "print semantic tableaux (output is LaTeX code)")
	proof      = flag.Bool("proof", false, "print outline of proof")
	dest       = flag.String("dest", "", "write output to named file (defaults to stdout)")
)

func init() {

	flag.Usage = func() {

		fmt.Println()

		fmt.Println(os.Args[0], "[commands/options]", "formula/file_name")

		fmt.Println()

		flag.PrintDefaults()
	}

	flag.Parse()
}
