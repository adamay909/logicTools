package main

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

func scrollTo(elem js.Value) {

	elem.Call("scrollIntoView")

	return

}

// addEventListener adds an event listener. The function f is called
// with the event object as its first argument and the arguments given by params.
func addEventListener(elem js.Value, eventType string, f func(event js.Value, args ...any), params ...any) {

	elem.Call("addEventListener", eventType, js.FuncOf(func(this js.Value, margs []js.Value) any {

		f(margs[0], params...)
		return true
	}), true)

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

func jswrapFunc(f func()) func(js.Value, ...any) {
	jsf := func(js.Value, ...any) {
		f()
	}
	return jsf
}
