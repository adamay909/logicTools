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

	if dsp.ypos > dsp.viewTop {
		return
	}

	e := dom.GetWindow().Document().GetElementByID("display").(dom.Node).Underlying()

	pxh := e.Get("clientHeight").Int()
	lh := pxh / 20
	sd := -lh

	e.Call("scrollBy", 0, sd)
	dsp.viewTop = dsp.viewTop - 1
	dsp.viewBottom = dsp.viewBottom - 1

}

func scrollDown() {

	if dsp.ypos < dsp.viewBottom {
		return
	}
	e := dom.GetWindow().Document().GetElementByID("display").(dom.Node).Underlying()

	pxh := e.Get("clientHeight").Int()
	lh := pxh / 20
	sd := lh

	e.Call("scrollBy", 0, sd)
	dsp.viewTop = dsp.viewTop + 1
	dsp.viewBottom = dsp.viewBottom + 1
}

func cursorUp() bool {

	return dsp.ypos < dsp.yprev

}

func cursorDown() bool {

	return dsp.ypos > dsp.yprev

}
