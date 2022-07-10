package main

import (
	"strconv"
	"strings"
)

func (d *console) typeset() string {

	html := `<div id="deriv">`
	if len(d.html) == 0 {
		html = html + "start typing"
	} else {
		html = html + strings.Join(d.html, "\n")
	}
	return html + `</div>`
}

func (d *console) format() {
	d.html = nil

	const dummyC = `<div id="cursor">&emsp;</div>`

	for n := range d.input {
		var dat, succ, annot []string
		var tstl, dots string

		tst := getTstIdx(d.input[n])
		dot := getDotIdx(d.input[n])

		ln := strconv.Itoa(n+d.offset) + `.&emsp;`

		r := `<div class="ln">#ln#</div><div class="ddat">#dat#</div><div class="dtstl">#tstl#</div><div class="succ">#succ#</div><div class="dsep">#dot#</div><div class="dannot">#annot#</div>`

		for i, e := range d.input[n] {
			var text string
			text = plainText(e)
			if n == d.ypos {
				switch {

				case i == d.xpos && d.modifier != "":
					text = `<div id="cursor">` + d.modifier + `</div>` + text

				case d.overhang && i == len(d.input[n])-1:
					text = text + dummyC

				case i == d.xpos:
					text = `<div id="cursor">` + text + `</div>`

				default:
				}
			}

			if i < tst {
				dat = append(dat, text)
				continue
			}

			if i == tst {
				tstl = text
				continue
			}

			if i < dot {
				succ = append(succ, text)
				continue
			}

			if i == dot {
				dots = text
				continue
			}

			annot = append(annot, text)
		}
		if n == d.ypos && len(d.input[n]) == 0 {
			dat = append(dat, dummyC)
		}

		r = strings.Replace(r, `#ln#`, ln, 1)
		r = strings.Replace(r, `#dat#`, stringOf(dat), 1)
		r = strings.Replace(r, `#tstl#`, tstl, 1)
		r = strings.Replace(r, `#succ#`, stringOf(succ), 1)
		r = strings.Replace(r, `#dot#`, dots, 1)
		r = strings.Replace(r, `#annot#`, stringOf(annot), 1)

		d.html = append(d.html, r)
	}
	if len(d.html) == 0 {
		ln := strconv.Itoa(0+d.offset) + `.&emsp;`
		r := `<div class="ln">` + ln + `</div><div class="ddat"><div id="cursor">&emsp;</div></div><div class="dtstl"></div><div class="succ"></div><div class="dsep"></div><div class="dannot"></div>`
		d.html = append(d.html, r)
	}
}

func (d *console) formatDerivation() {
	d.html = nil
	tstl := plainText(`\vdash`)
	dots := plainText(`\ldots`)

	for n := range d.input {
		dat, succ, annot, err := parseLineDisplay(d.input[n])

		if err != nil {
			continue
		}

		ln := strconv.Itoa(n+d.offset) + `.&emsp;`

		r := `<div class="ln">#ln#</div><div class="ddat">#dat#</div><div class="dtstl">#tstl#</div><div class="succ">#succ#</div><div class="dsep">#dot#</div><div class="dannot">#annot#</div>`
		r = strings.Replace(r, `#ln#`, ln, 1)
		r = strings.Replace(r, `#dat#`, dat, 1)
		r = strings.Replace(r, `#tstl#`, tstl, 1)
		r = strings.Replace(r, `#succ#`, succ, 1)
		r = strings.Replace(r, `#dot#`, dots, 1)
		r = strings.Replace(r, `#annot#`, annot, 1)

		d.html = append(d.html, r)
	}
	if len(d.html) == 0 {
		ln := strconv.Itoa(0+d.offset) + `.&emsp;`
		r := `<div class="ln">` + ln + `</div><div class="ddat"><div id="cursor">&emsp;</div></div><div class="dtstl"></div><div class="succ"></div><div class="dsep"></div><div class="dannot"></div>`
		d.html = append(d.html, r)
	}
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
