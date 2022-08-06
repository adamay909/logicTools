package main

import (
	"strconv"
	"strings"
)

func (d *console) typeset() string {

	var html string

	html = `<h3 id="extitle"></h3>`
	if d.Title != "" {
		html = `<h3 id="extitle">` + prettyGreek(d.Title) + `</h3>`
	}
	html = html + `<div id="deriv">`
	html = html + strings.Join(d.html, "\n")
	return html + `</div>`
}

func (d *console) format() {
	cursor := true
	d.formatHTML(cursor)
}

func (d *console) formatDerivation() {
	cursor := false
	d.formatHTML(cursor)
}

func (d *console) formatHTML(cursor bool) {

	d.html = nil

	for n, l := range d.Input {
		var dat, succ, annot []string

		if cursor && n == d.ypos {
			l = setCursor(d, d.Input[n])
		}
		raw := true
		dat, tstl, succ, dots, annot := parseNsplit(l, !raw)

		ln := strconv.Itoa(n+d.Offset) + `.&emsp;`

		r := `<div class="ln">#ln#</div><div class="ddat">#dat#</div><div class="dtstl">#tstl#</div><div class="succ">#succ#</div><div class="dsep">#dot#</div><div class="dannot">#annot#</div>`

		r = strings.Replace(r, `#ln#`, ln, 1)
		r = strings.Replace(r, `#dat#`, stringOf(dat), 1)
		r = strings.Replace(r, `#tstl#`, stringOf(tstl), 1)
		r = strings.Replace(r, `#succ#`, stringOf(succ), 1)
		r = strings.Replace(r, `#dot#`, stringOf(dots), 1)
		r = strings.Replace(r, `#annot#`, stringOf(annot), 1)

		d.html = append(d.html, r)
	}

	if len(d.html) == 0 {
		ln := strconv.Itoa(0+d.Offset) + `.&emsp;`
		r := `<div class="ln">` + ln + `</div><div class="ddat"><div id="cursor">&thinsp;</div></div><div class="dtstl"></div><div class="succ"></div><div class="dsep"></div><div class="dannot"></div>`
		d.html = append(d.html, r)
	}
}

func (d *console) setTitle(t string) {
	d.Title = t
}

func setCursor(d *console, l []string) []string {

	rv := insertCursor(l, d.xpos)
	return rv
}

func unsetCursor(l []string) []string {
	var cursor = `<div id="cursor">&thinsp;</div>`

	if dsp.modifier != "" {
		cursor = strings.ReplaceAll(cursor, `&thinsp;`, dsp.modifier)
	}

	var rv []string

	for _, e := range l {
		if e == cursor {
			continue
		}
		rv = append(rv, e)
	}
	return rv
}

func insertCursor(l []string, p int) []string {

	var m []string

	var cursor = `<div id="cursor">&thinsp;</div>`

	if dsp.modifier != "" {
		cursor = strings.ReplaceAll(cursor, `&thinsp;`, dsp.modifier)
	}

	m = append(m, l[:p]...)
	m = append(m, cursor)
	m = append(m, l[p:]...)
	return m
}

func parseNsplit(l inputLine, raw bool) (dat, tstl, succ, dots, annot []string) {

	tst := getTstIdx(l)
	dot := getDotIdx(l)

	for i, e := range l {
		var text string
		if !raw {
			text = plainHTML(e)
		} else {
			text = plainText(e)
		}

		if i < tst {
			dat = append(dat, text)
			continue
		}

		if i == tst {
			tstl = append(tstl, text)
			continue
		}

		if i < dot {
			succ = append(succ, text)
			continue
		}

		if i == dot {
			dots = append(dots, text)
			continue
		}

		annot = append(annot, text)
	}
	return
}

func getTstIdx(l []string) int {

	tst := index(l, `\vdash`)
	dot := index(l, `\ldots`)

	if tst == -1 {
		return len(l)
	}

	if dot > -1 && dot < tst {
		return len(l)
	}

	return tst
}

func getDotIdx(l []string) int {

	if getTstIdx(l) == len(l) {
		return -1
	}

	idx := index(l, `\ldots`)
	if idx == -1 {
		idx = len(l)
	}
	return idx
}

func stringOf(src []string) string {
	return strings.Join(src, "")
}

func spaceyStringOf(src []string) string {
	return strings.Join(src, " ")
}
