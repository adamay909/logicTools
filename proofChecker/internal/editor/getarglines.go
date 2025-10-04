package editor

import (
	"fmt"
	"strings"
)

func (e *Editor) GetArglines() (lines []string) {

	if e.editorType != derivationeditor && e.editorType != axiomaticeditor {
		return
	}

	rows := e.elem.Call("querySelector", ".derivation").Get("children")

	for i := 0; i < rows.Get("length").Int(); i++ {
		r := rows.Call("item", i)
		parts := r.Get("children")
		l := ""
		for j := 0; j < parts.Get("length").Int(); j++ {
			l = l + parts.Call("item", j).Get("innerHTML").String()
		}
		l = strings.ReplaceAll(l, "&nbsp;", "")
		l = strings.ReplaceAll(l, turnstile, ":")
		l = strings.ReplaceAll(l, ldots, ".")
		l = strings.ReplaceAll(l, `<sub>`, "_")
		l = strings.ReplaceAll(l, `</sub>`, "")
		l = toascii(l)
		l = strings.ReplaceAll(l, `\`, `/`)
		dummy := domDocument.Call("createElement", "div")
		dummy.Set("innerHTML", l)
		for c := dummy.Get("children"); c.Get("length").Int() != 0; c = dummy.Get("children") {
			c.Call("item", 0).Call("remove")
		}
		l = dummy.Get("textContent").String()
		if e.editorType == axiomaticeditor {
			l = ":" + l
		}
		if l == ":." {
			l = ""
		}
		dummy.Set("innerHTML", l)
		lines = append(lines, StripWhiteSpaceFromHTML(dummy.Get("textContent").String()))
		fmt.Println(l)
	}
	return
}

func toascii(s string) string {
	cl := []rune(s)
	var r string
	for _, c := range cl {
		if a, ok := revMap[string(c)]; ok {
			r = r + a
		} else {
			r = r + string(c)
		}
	}
	return r
}
