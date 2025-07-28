/*
Package fileops offers a few helpers for file operation to make life easier
*/
//go:generate pkgdoc2readme
package fileops

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func confirmOverwrite(name string) bool {

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s already exists. Overwrite? (y/N) ", name)

	response, _ := reader.ReadString('\n')

	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"

}

func CreateFile(name string) *os.File {

	if FileExists(name) {

		if ok := confirmOverwrite(name); !ok {
			fmt.Println("Write openration aborted.")
			os.Exit(1)
		}
	}

	dir := filepath.Dir(name)

	_, err := os.Open(dir)

	if os.IsNotExist(err) {

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return file
}

func FileExists(name string) bool {

	_, err := os.Stat(name)

	return err == nil

}
