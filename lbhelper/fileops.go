package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findFile(shortname string, exts []string) (realname string, exists bool) {

	exists = false
	for ext := range exts {
		realname := shortname + `.` + exts[ext]
		if fileExists(realname) {
			exists = true
			return realname, exists
		}
	}
	return shortname, exists
}

//look for files of a given name and candidate extensions recursively.
//Can look for a specific file name by giving full name as shortname and
//leaving exts empty. The returned string is path relative to CWD.
func findGuess(shortname string, exts []string) (string, bool) {
	depth := 4
	for ext := range exts {
		fullname := filepath.Base(shortname) + `.` + exts[ext]
		for i := 0; i <= depth; i++ {

			match, _ := filepath.Glob(fullname)
			if match != nil {
				return match[0], true
			}

			fullname = `*/` + fullname
		}
	}
	return shortname, false
}

func copyFiles(in, out string) {

	indata, _ := ioutil.ReadFile(in)

	err := ioutil.WriteFile(out, indata, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func fileExists(filename string) bool {

	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func fileWrite(filename string, data string) {
	if len(filename) == 0 {
		return
	}
	if len(data) == 0 {
		return
	}
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func fileRead(filename string) (out string) {
	if len(filename) == 0 {
		return
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	if data == nil {
		return
	}
	return string(data)
}

func setupLogger(file string) *log.Logger {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	//	defer f.Close()
	logger := log.New(f, "", 0)
	logger.Println(strings.Repeat("=", 80))
	return logger
}

func applyPatch(orig, patch string) {
	cmd := exec.Command("patch", "-f", orig, patch)
	cmd.Run()
}

func fileDelete(p string) {
	files, _ := filepath.Glob(p)
	for _, f := range files {
		os.Remove(f)
	}
	return
}
