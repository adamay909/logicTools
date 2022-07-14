package main

type inputLine []string

type console struct {
	Title           string
	Input           []inputLine
	xpos, ypos      int
	xprev, yprev    int
	xcuror, ycursor int
	html            []string
	modifier        string
	overhang        bool
	Offset          int
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

	d.xprev = d.xpos
	d.yprev = d.ypos

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
	if len(d.Input) == 0 {
		d.overhang = true
		return
	}
	if d.xpos == len(d.Input[d.ypos]) {
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
	d.home()
	/*
	   if d.xpos > len(d.Input[d.ypos]) {
	   		d.end()
	   		d.arrowRight()
	   	}
	*/
}

func (d *console) arrowDown() {
	if d.ypos == len(d.Input)-1 {
		return
	}
	d.ypos++
	d.home()
	/*
	   if d.xpos > len(d.Input[d.ypos]) {
	   		d.end()
	   		d.arrowRight()
	   			}
	*/
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
			d.Input[d.ypos] = d.currentLine()[1:]
		}
		return
	}
	l1 = d.currentLine()[:d.xpos]
	if d.xpos < len(d.currentLine())-1 {
		l2 = d.currentLine()[d.xpos+1:]
	}

	d.Input[d.ypos] = append(l1, l2...)

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

	if d.ypos == len(d.Input)-1 {
		d.Input = d.Input[:d.ypos]
		d.arrowUp()
		d.home()
		return
	}
	if d.ypos == 0 {
		d.Input = d.Input[1:]
		d.arrowUp()
		d.home()
		return
	}
	l1 := d.Input[:d.ypos]
	l2 := d.Input[d.ypos+1:]
	d.Input = nil
	d.Input = append(l1, l2...)

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
		tk, err := tkOf(m, tkraw, tktex, allBindings)
		if err == nil {
			d.addChar(tk)
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
		newlines = append(newlines, d.Input[i])
	}

	newlines = append(newlines, frag1)
	newlines = append(newlines, frag2)

	if d.ypos < len(d.Input)-1 {
		for i := d.ypos + 1; i < len(d.Input); i++ {
			newlines = append(newlines, d.Input[i])
		}
	}

	d.Input = nil
	d.Input = newlines
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
	d.Input[d.ypos] = n

	d.xpos++
}

func (d *console) clear() {
	d.Input = nil
	d.html = nil
	d.xpos = 0
	d.ypos = 0
	d.overhang = true
	d.modifier = ""
	d.Offset = 1
	d.Title = ""

}

func (d *console) currentLine() inputLine {
	if len(d.Input) == 0 {
		d.Input = append(d.Input, make(inputLine, 0))
	}
	return d.Input[d.ypos]
}

func (d *console) empty() bool {
	return len(d.Input) == 0
}

func isModifier(k string) bool {

	return k == `\` || k == `|`

}

func (d *console) setOffset(n int) {
	d.Offset = n
}
