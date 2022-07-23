package main

import (
	"honnef.co/go/js/dom/v2"
)

func scrollDisplay() {

	if cursorUp() {
		scrollUp()
	}

	if cursorDown() {
		scrollDown()
	}

	return
}

func scrollUp() {

	e := dom.GetWindow().Document().GetElementByID("display").(dom.Node).Underlying()

	sh := e.Get("scrollHeight").Int()
	ch := e.Get("clientHeight").Int()

	ln := len(dsp.Input)
	if ln == 0 {
		ln = 1
	}

	lh := float64(sh) / float64(ln)

	maxl := int(float64(ch) / lh)

	dsp.viewBottom = dsp.viewTop + maxl

	if dsp.ypos > dsp.viewTop {
		return
	}

	e.Call("scrollBy", 0, -lh)
	dsp.viewTop = dsp.viewTop - 1
	dsp.viewBottom = dsp.viewBottom - 1

}

func scrollDown() {

	e := dom.GetWindow().Document().GetElementByID("display").(dom.Node).Underlying()

	sh := e.Get("scrollHeight").Int()
	ch := e.Get("clientHeight").Int()

	ln := len(dsp.Input)
	if ln == 0 {
		ln = 1
	}

	lh := float64(sh) / float64(ln)

	maxl := int(float64(ch) / lh)

	dsp.viewBottom = dsp.viewTop + maxl

	if dsp.ypos < dsp.viewBottom {
		return
	}

	e.Call("scrollBy", 0, lh+5)
	dsp.viewTop = dsp.viewTop + 1
	dsp.viewBottom = dsp.viewBottom + 1
}

func cursorUp() bool {

	return dsp.ypos < dsp.yprev

}

func cursorDown() bool {

	return dsp.ypos > dsp.yprev

}
