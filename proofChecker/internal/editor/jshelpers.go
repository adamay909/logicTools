package editor

import (
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"syscall/js"
)

var (
	uint8ArrayConstructor,
	fileConstructor,
	arrayConstructor,
	blobConstructor,
	domWindow,
	domDocument,
	domHTML,
	domBody js.Value
)

type caretLocation struct {
	caretNode   js.Value
	caretOffset js.Value
}

func init() {

	uint8ArrayConstructor = js.Global().Get("Uint8Array")

	fileConstructor = js.Global().Get("File")

	arrayConstructor = js.Global().Get("Array")

	blobConstructor = js.Global().Get("Blob")

	domWindow = js.Global()

	domDocument = domWindow.Get("document")

	domHTML = domDocument.Get("documentElement")

	domBody = domDocument.Get("body")
}

func uint8arrayOf(data []byte) js.Value {

	jsdata := uint8ArrayConstructor.New(len(data))

	js.CopyBytesToJS(jsdata, data)

	return jsdata

}

func saveFile(file js.Value) {

	href := domWindow.Get("URL").Call("createObjectURL", file)

	link := domDocument.Call("createElement", "a")

	link.Set("href", href)

	link.Set("download", file.Get("name"))

	domBody.Call("appendChild", link)

	link.Call("click")

	domBody.Call("removeChild", link)

	domWindow.Get("URL").Call("revokeObjectURL", href)

	return

}

// Thanks to https://javascript.plainenglish.io/javascript-create-file-c36f8bccb3be for how to do this
func createJSFile(data []byte, name string) js.Value {

	jsdata := uint8arrayOf(data)

	blob := blobConstructor.New(arrayConstructor.New(jsdata), map[string]any{"type": mime.TypeByExtension(filepath.Ext(name))})

	return fileConstructor.New(arrayConstructor.New(blob), filepath.Base(name))

}

func getURL() string {

	return domWindow.Get("location").Get("href").String()

}

func getHost() string {

	return domWindow.Get("location").Get("host").String()

}

func getHash() string {

	return domWindow.Get("location").Get("hash").String()

}

func setHash(h string) {

	domWindow.Get("location").Set("hash", h)

	return

}

func replaceBody(p string) {

	domBody.Set("innerHTML", p)

	return

}

func fetchData(path string) (data []byte) {

	//	loc := path.String()

	log.Println("fetching", path)

	r, err := http.Get(path)

	if err != nil {
		log.Println(err)
		return
	}

	data, _ = io.ReadAll(r.Body)

	return data

}

func getElementByID(id string) (elem js.Value, err error) {

	elem = domDocument.Call("getElementById", id)

	if elem.IsNull() {
		err = errors.New("no element found with ID " + id)
		log.Println(err)
	}

	return
}

func getElementsByClassName(elem js.Value, class string) (list []js.Value) {

	htmlCollection := elem.Call("getElementsByClassName", class)

	l := htmlCollection.Get("length").Int()

	for i := 0; i < l; i++ {

		list = append(list, htmlCollection.Call("item", i))
	}

	return list

}

func childOfClass(e js.Value, class string) js.Value {

	htmlCollection := e.Get("children")
	l := htmlCollection.Get("length").Int()
	for i := 0; i < l; i++ {
		if htmlCollection.Call("item", i).Get("className").String() == class {
			return htmlCollection.Call("item", i)
		}
	}
	return js.ValueOf(nil)
}

func scrollTo(elem js.Value) {

	elem.Call("scrollIntoView")

	return

}

// addEventListener adds an event listener. The function f is called
// with the event object as its first argument and the arguments given by params.
func addEventListener(elem js.Value, eventType string, f func(event js.Value, args ...any), params ...any) {

	elem.Call("addEventListener", eventType, js.FuncOf(func(this js.Value, margs []js.Value) any {
		//margs[0] is the event that gets passed to the function.
		f(margs[0], params...)
		return true
	}), false)

	return
}

func jsWrapper(f func(event js.Value, args ...any), params ...any) js.Func {

	return js.FuncOf(func(this js.Value, margs []js.Value) any {
		f(margs[0], params...)
		return true
	})

}

func removeEventListener(elem js.Value, eventType string, f func(event js.Value, args ...any), params ...any) {

	elem.Call("removeEventListener", eventType, js.FuncOf(func(this js.Value, margs []js.Value) any {
		f(margs[0], params...)
		return true
	}), true)

	return

}

func inactivateElement(elem js.Value) {

	newEl := createElement("div", "")

	newEl.Call("setAttribute", "class", "cover")

	elem.Call("append", newEl)

}

func reactivateElement(elem js.Value) {

	cover := getElementsByClassName(elem, "cover")[0]

	cover.Call("remove")

	return
}

func createElement(tag string, innerHTML string) js.Value {

	newEl := domDocument.Call("createElement", tag)

	newEl.Set("innerHTML", innerHTML)

	return newEl

}

func coverElement(elem js.Value, opacity int) {

	inactivateElement(elem)

	cover := getElementsByClassName(elem, "cover")[0]

	cover.Call("setAttribute", "style", "opacity: "+strconv.Itoa(opacity)+"%;")

}

func uncoverElement(elem js.Value) {

	covers := getElementsByClassName(domBody, "cover")

	//just in case we covered elem multiple times

	for _, c := range covers {

		c.Call("remove")

	}

}

func coverScreen(opacity int) {

	coverElement(domBody, opacity)

}

func uncoverScreen() {

	uncoverElement(domBody)

}

func coverAndWait(elem js.Value, opacity int) {

	coverElement(elem, opacity)

	cover := getElementsByClassName(elem, "cover")[0]

	addStyle(cover, "cursor: wait;")

}

func addStyle(elem js.Value, css string) {

	css1 := elem.Call("getAttribute", "style").String()

	elem.Call("setAttribute", "style", css1+css)

}

func writeConsoleLog(msg ...string) {

	logmsg := js.Global().Get("Date").New().Call("toISOString").String() + strings.Join(msg, " ")

	domWindow.Get("console").Call("log", logmsg)

	return
}

func elementAtCaret() js.Value {

	selection := domWindow.Call("getSelection")
	srange := selection.Call("getRangeAt", 0)
	an := srange.Get("commonAncestorContainer")
	if an.Get("nodeType").Int() == 1 {
		return an
	}
	return an.Get("parentNode")

	//	return domWindow.Call("getSelection").Get("focusNode").Get("parentElement")

}

func nodeAtCaret() js.Value {

	selection := domWindow.Call("getSelection")
	srange := selection.Call("getRangeAt", 0)
	return srange.Get("startContainer")

}

/*
Find the next element of specified class. If found, it returns the element and the index irelative to current element given by curElem. -1 indicates not found. 0 indicates the current element is of the relevant class.
*/
func nextElementOfClass(curElem js.Value, name string) (js.Value, int) {

	index := -1
	for e := curElem; !e.IsNull(); e = e.Get("nextElementSibling") {
		index++
		if e.Call("getAttribute", "class").String() == name {
			return e, index
		}
	}
	return js.ValueOf(nil), -1
}

/*
Find the previous element of specified class. If found, it returns the element and the index irelative to current element given by curElem. -1 indicates not found. 0 indicates the current element is of the relevant class.
*/
func previousElementOfClass(curElem js.Value, name string) (js.Value, int) {

	index := -1
	for e := curElem; !e.IsNull(); e = previousElement(e) {
		index++
		if e.Call("getAttribute", "class").String() == name {
			return e, index
		}
	}
	return js.ValueOf(nil), -1
}

func nextElement(e js.Value) js.Value {
	if isElement(e) {
		return e.Get("nextElementSibling")
	}
	for n := e.Get("nextSibling"); !n.IsNull(); n = n.Get("nextSibling") {
		if isElement(n) {
			return n
		}
	}
	return js.ValueOf(nil)
}

func previousElement(e js.Value) js.Value {
	if isElement(e) {
		return e.Get("previousElementSibling")
	}
	for n := e.Get("previousSibling"); !n.IsNull(); n = n.Get("previousSibling") {
		if isElement(n) {
			return n
		}
	}
	return js.ValueOf(nil)
}

func rejectSelection() {
	selection := domWindow.Call("getSelection")
	srange := selection.Call("getRangeAt", 0)
	if selection.Get("direction").String() == "forward" {
		srange.Call("collapse")
		return
	}
	srange.Call("collapse", "toStart")
}

func caretOffset(editor js.Value) int {
	selection := domWindow.Call("getSelection")
	srange := selection.Call("getRangeAt", 0)
	crange := srange.Call("cloneRange")
	e := getElementsByClassName(editor, "derivation")[0]
	crange.Call("setStart", e, 0)
	return len(crange.Call("toString").String())
	//Get("startOffset").Int()
}

func classOf(e js.Value) string {
	return e.Call("getAttribute", "class").String()
}

func valOf(e js.Value, k string) string {
	if e.IsNull() {
		return "<undefined>"
	}
	return e.Call("getAttribute", k).String()
}

func isElement(e js.Value) bool {
	return e.Get("nodeType").Int() == 1
}

func getTextOfLine(c js.Value) string {
	start, i := previousElementOfClass(c, "ls")
	if i == -1 {
		panic("line start not found")
	}
	end, i := nextElementOfClass(c, "le")
	if i == -1 {
		panic("line end not found")
	}
	srange := domDocument.Call("createRange")
	srange.Call("setStartAfter", start, 0)
	srange.Call("setEndBefore", end)
	return srange.Call("toString").String()
}

func setCaretAt(e js.Value) {
	selection := domWindow.Call("getSelection")
	srange := domDocument.Call("createRange")
	e.Set("innerHTML", "<div></div>")
	srange.Call("setStart", e, 1)
	srange.Call("collapse", "true")
	selection.Call("addRange", srange)
}

func currentCaretLocation() *caretLocation {
	srange := domWindow.Call("getSelection").Call("getRangeAt", 0)
	l := new(caretLocation)
	l.caretNode = srange.Get("startContainer")
	l.caretOffset = srange.Get("startOffset")
	return l
}

func (l *caretLocation) setToCurrent() {
	srange := domWindow.Call("getSelection").Call("getRangeAt", 0)
	l.caretNode = srange.Get("startContainer")
	l.caretOffset = srange.Get("startOffset")
}

func (l *caretLocation) set(n js.Value, offset int) {
	l.caretNode = n
	l.caretOffset = js.ValueOf(offset)
}

func (l *caretLocation) void() {
	l.caretNode = js.ValueOf(nil)
	l.caretOffset = js.ValueOf(nil)
}

func (l *caretLocation) get() (caretNode, caretOffset js.Value) {
	return l.caretNode, l.caretOffset
}

func (l *caretLocation) elementContainer() js.Value {
	if l.caretNode.Get("nodeType").Int() == 1 {
		return l.caretNode
	}
	return l.caretNode.Get("parentElement")
}

func (l *caretLocation) String() string {
	if l.caretNode.Get("nodeType").Int() == 1 {
		return l.caretNode.Get("outerHTML").String()
	}
	return "content:" + l.caretNode.Get("textContent").String() + " offset:" + strconv.Itoa(l.caretOffset.Int())
}

func (l *caretLocation) _moveTo() {
	srange := domWindow.Call("getSelection").Call("getRangeAt", 0)
	srange.Call("insertNode", l.caretNode)
	srange.Call("setStart", l.caretNode, l.caretOffset)
}

func moveCaretToNextElementStart(current js.Value) {
	sel := domWindow.Call("getSelection")
	s := nextElement(current)
	r := domDocument.Call("createRange")
	t := domDocument.Call("createTextNode", "!")
	s.Call("prepend", t)
	r.Call("setStart", t, 0)
	sel.Call("removeAllRanges")
	sel.Call("addRange", r)
	t.Call("remove")

}

func moveCaretToPreviousElementEnd(current js.Value) {
	sel := domWindow.Call("getSelection")
	s := previousElement(current)
	r := domDocument.Call("createRange")
	t := domDocument.Call("createTextNode", "!")
	s.Call("append", t)
	r.Call("setStart", t, 0)
	sel.Call("removeAllRanges")
	sel.Call("addRange", r)
	t.Call("remove")

}

func moveCaretTo(target js.Value, offset int) {
	sel := domWindow.Call("getSelection")
	r := domDocument.Call("createRange")
	r.Call("setStart", target, 0)
	sel.Call("removeAllRanges")
	sel.Call("addRange", r)
}

func moveCaretRight() {
	sel := domWindow.Call("getSelection")
	sel.Call("modify", "move", "forward", "character")
	//	r.Call("setStart", r.Get("startContainer"), r.Get("startOffset").Int()+1)
	r := sel.Call("getRangeAt", 0)
	if r.Get("startOffset").Int() == len(r.Get("startContainer").Get("textContent").String()) {
		moveCaretToNextElementStart(elementAtCaret())
	}
}

func moveCaretLeft() {
	sel := domWindow.Call("getSelection")
	sel.Call("move", "backward", "character")
}
