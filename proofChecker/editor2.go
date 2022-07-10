package main

type inputLine []string

type console struct {
	input      []inputLine
	xpos, ypos int
	html       []string
	modifier   string
	overhang   bool
	offset     int
}

const (
	up         = "ArrowUp"
	down       = "ArrowDown"
	right      = "ArrowRight"
	left       = "ArrowLeft"
	home2      = "Home"
	end2       = "End"
	del        = "Delete"
	backspace2 = "Backspace"
	enter      = "Enter"
)

func (d *console) handleInput(key string) {

	if key == "Shift" {
		return
	}

	if isModifier(key) {
		d.modifier = key
		return
	}

	d.cursorMovement(key)
	d.deletion(key)
	d.insertion(key)
	d.checkOverhang()
}

func (d *console) checkOverhang() {
	if len(d.input) == 0 {
		d.overhang = true
		return
	}
	if d.xpos == len(d.input[d.ypos]) {
		d.overhang = true
	} else {
		d.overhang = false
	}
}

func (d *console) cursorMovement(m string) {

	switch m {
	case up:
		d.arrowUp()
	case down:
		d.arrowDown()
	case left:
		d.arrowLeft()
	case right:
		d.arrowRight()
	case home2:
		d.home()
	case end2:
		d.end()
	default:
		return
	}
}

func (d *console) arrowUp() {
	if d.ypos == 0 {
		return
	}
	d.ypos--
	if d.xpos > len(d.input[d.ypos]) {
		d.end()
		d.arrowRight()
	}
}

func (d *console) arrowDown() {
	if d.ypos == len(d.input)-1 {
		return
	}
	d.ypos++
	if d.xpos > len(d.input[d.ypos]) {
		d.end()
		d.arrowRight()
	}
}

func (d *console) arrowLeft() {
	if d.xpos == 0 {
		return
	}
	d.xpos--
}

func (d *console) arrowRight() {
	if d.xpos == len(d.currentLine()) {
		return
	}
	if d.xpos == len(d.currentLine())-1 {
		d.xpos++
		d.overhang = true
		return
	}
	d.xpos++
	d.overhang = false
}

func (d *console) home() {
	d.xpos = 0
}

func (d *console) end() {
	d.xpos = len(d.currentLine()) - 1
	d.arrowRight()
}

func (d *console) deletion(m string) {

	switch m {
	case del:
		d.deleteChar()
	case backspace2:
		d.backspace()
	default:
		return
	}
}

func (d *console) deleteChar() {
	var l1, l2 inputLine

	if d.xpos == 0 {
		if len(d.currentLine()) == 0 {
			d.deleteLine()
		} else {
			d.input[d.ypos] = d.currentLine()[1:]
		}
		return
	}
	l1 = d.currentLine()[:d.xpos]
	if d.xpos < len(d.currentLine())-1 {
		l2 = d.currentLine()[d.xpos+1:]
	}

	d.input[d.ypos] = append(l1, l2...)

}

func (d *console) backspace() {
	if d.xpos == 0 && d.overhang && d.ypos > 0 {
		d.deleteLine()
		d.arrowUp()
		d.end()
		d.arrowRight()
		return
	}
	if d.xpos == 0 {
		return
	}
	d.arrowLeft()
	d.deleteChar()
}

func (d *console) deleteLine() {

	if d.empty() {
		return
	}

	if d.ypos == len(d.input)-1 {
		d.input = d.input[:d.ypos]
		d.arrowUp()
		d.home()
		return
	}
	if d.ypos == 0 {
		d.input = d.input[1:]
		d.arrowUp()
		d.home()
		return
	}
	l1 := d.input[:d.ypos]
	l2 := d.input[d.ypos+1:]
	d.input = nil
	d.input = append(l1, l2...)

}

func (d *console) insertion(m string) {
	if d.modifier != "" {
		m = d.modifier + m
		d.modifier = ""
	}
	switch m {
	case enter:
		d.addNewline()
	default:
		if ok, target := among(m, keyBindings, punctBindings, logConstBindings, turnstileBindings, greekBindings); ok {
			d.addChar(target)
		}
	}
}

func (d *console) addNewline() {

	var frag1, frag2 []string

	frag1 = d.currentLine()[:d.xpos]

	if !d.overhang {
		frag2 = d.currentLine()[d.xpos:]
	}

	var newlines []inputLine

	for i := 0; i < d.ypos; i++ {
		newlines = append(newlines, d.input[i])
	}

	newlines = append(newlines, frag1)
	newlines = append(newlines, frag2)

	if d.ypos < len(d.input)-1 {
		for i := d.ypos + 1; i < len(d.input); i++ {
			newlines = append(newlines, d.input[i])
		}
	}

	d.input = nil
	d.input = newlines
	d.ypos++
	d.xpos = 0
	return
}

func (d *console) addChar(c string) {
	// 	if xpos == len(d.currentLine()) {
	//	d.input[d.ypos] = append(d.currentLine(), c)
	//	d.xpos++
	//	return
	//  }

	var n inputLine

	for i := 0; i < d.xpos; i++ {
		n = append(n, d.currentLine()[i])
	}
	n = append(n, c)
	for i := d.xpos; i < len(d.currentLine()); i++ {
		n = append(n, d.currentLine()[i])
	}
	d.input[d.ypos] = n

	d.xpos++
}

func (d *console) clear() {
	d.input = nil
	d.html = nil
	d.xpos = 0
	d.ypos = 0
	d.overhang = true
	d.modifier = ""
	d.offset = 1

}

func (d *console) currentLine() inputLine {
	if len(d.input) == 0 {
		d.input = append(d.input, make(inputLine, 0))
	}
	return d.input[d.ypos]
}

func (d *console) empty() bool {
	return len(d.input) == 0
}

func isModifier(k string) bool {

	return k == `\` || k == `|`

}

func (d *console) setOffset(n int) {
	d.offset = n
}

func among(s string, bindings ...[][3]string) (bool, string) {

	for _, b := range bindings {
		for _, e := range b {
			if s == e[0] {
				return true, e[1]
			}
		}
	}
	return false, ""
}
