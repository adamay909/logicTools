package main

import (
	"errors"
	"strconv"
	"strings"
)

func text2data(text string) (resp []inputLine, err error) {

	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "[", "(")
	text = strings.ReplaceAll(text, "]", ")")
	text = strings.ReplaceAll(text, "{", "(")
	text = strings.ReplaceAll(text, "}", ")")

	lines := strings.Split(text, "\n")

	for i, l := range lines {
		if len(l) < 1 {
			break
		}
		var r []string
		//make sure line starts with number
		p := strings.Index(l, ".")
		if p == -1 {
			err = errors.New(strconv.Itoa(i+1) + ": no line number")
			return
		}
		_, err = strconv.Atoi(l[:p])
		if err != nil {
			return
		}
		if p+1 == len(l) {
			err = errors.New(strconv.Itoa(i+1) + ": nothing after line number")
			return
		}
		l = l[p+1:]
		for n := 0; n < len(l); {
			var t string
			var pos int
			t, pos, err = firstToken(l[n:])
			if err != nil {
				err = errors.New(strconv.Itoa(i+1) + ":" + l[n:] + ": not valid character")
				return
			}
			r = append(r, t)
			n = n + pos
		}
		resp = append(resp, r)
		r = nil
	}
	return
}

func firstToken(s string) (r string, l int, err error) {

	for _, e := range allBindings {
		if strings.HasPrefix(s, e[tktxt]) {
			r = e[tktex]
			l = len(e[tktxt])
			return
		}
	}
	err = errors.New("not found")
	return
}