package main

import (
	_ "embed" //embed
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adamay909/logicTools/fileops"
)

func processFile() {

	input := flag.Arg(0)

	if input == "" {
		os.Stderr.WriteString("You need to supply file name.\n")
		os.Exit(1)
	}

	if !fileops.FileExists(input) {
		os.Stderr.WriteString(input + " does not exist.\n")
		os.Exit(1)
	}

	if *cmpl {
		compileFile(input, *ans, *letter, *fontsize, *resize)
		return
	}

	if *standalone {
		mkLaTeX(input, *ans, *letter, *fontsize, *resize)
		return
	}

}

func compileFile(input string, ans bool, letter bool, fontsize string, resize bool) {

	f := input
	if filepath.Ext(f) == "" {
		f = f + ".tex"
	}
	compileFragment(f, ans, letter, fontsize, resize)
	return
}

//go:embed assets/latexSkel.tex
var latexSkel string

func mkLaTeX(f string, ans bool, letter bool, fontsize string, resize bool) {

	nameBase := stripExt(filepath.Base(f)) + "FULL"

	name := nameBase + ".tex"

	data := fixLatexBody(f, ans, letter, fontsize, resize)

	writeFile(name, data)
	return
}

func fixLatexBody(f string, ans bool, letter bool, fontsize string, resize bool) string {

	text := fileRead(f)

	lines := strings.Split(text, "\n")

	for count, l := range lines {

		if !strings.HasPrefix(l, "%") {
			break
		}

		if count > 3 {
			break
		}

		if strings.HasPrefix(l, "%lbhelper:") {

			cmds := strings.Split(strings.TrimPrefix(l, "%lbhelper:"), ",")

			for _, e := range cmds {

				if e == "letterpaper" {
					letter = true
				}

				if strings.HasPrefix(e, "fontsize=") {
					fontsize = strings.TrimPrefix(e, "fontsize=")
				}

				if e == "resize" {
					resize = true
				}
			}
		}

		if !ans {
			continue
		}

		if strings.HasPrefix(l, "%lbhelper-answer:") {

			cmds := strings.Split(strings.TrimPrefix(l, "%lbhelper-answer:"), ",")

			for _, e := range cmds {

				if e == "letterpaper" {
					letter = true
				}

				if strings.HasPrefix(e, "fontsize=") {
					fontsize = strings.TrimPrefix(e, "fontsize=")
				}

				if e == "resize" {
					resize = true
				}
			}
		}
	}

	if resize {
		text = `\resizebox{\txw}{!}{` + text + `}`
	}
	data := strings.Replace(latexSkel, "#DUMMY#", text, 1)

	data = strings.Replace(data, "#FONTSIZE#", fontsize, 1)
	if ans {
		data = strings.Replace(data, "#ANSWERS#", `\usepackage{ex-ans}`, 1)
	} else {
		data = strings.Replace(data, "#ANSWERS#", "", 1)
	}
	if letter {
		data = strings.Replace(data, "#PAPERSIZE#", "letterpaper", 1)
	} else {
		data = strings.Replace(data, "#PAPERSIZE#", "paperwidth=15cm, paperheight=20cm", 1)
	}

	return data
}

func compileFragment(f string, ans bool, letter bool, fontsize string, resize bool) {

	nameBase := stripExt(filepath.Base(f)) + "TMP"
	if ans {
		nameBase = stripExt(filepath.Base(f)) + "_answers_TMP"
	}
	name := nameBase + ".tex"

	data := fixLatexBody(f, ans, letter, fontsize, resize)

	writeFile(name, data)

	compileLatex(name)

	if ans {
		os.Rename(nameBase+".pdf", stripExt(filepath.Base(f))+"_answers.pdf")
	} else {
		os.Rename(nameBase+".pdf", stripExt(filepath.Base(f))+".pdf")
	}

	fileDelete(nameBase + "*")
	return
}

func compileLatex(f string) {
	name := stripExt(filepath.Base(f))
	cmd := exec.Command("latexmk", "-xelatex", name)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	return
}

func stripExt(s string) string {

	return strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))

}

func writeFile(name string, data string) {

	file, err := os.Create(name)

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	defer file.Close()

	file.WriteString(data)

}
