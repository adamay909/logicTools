package main

func (d *console) deletion(m string) {

	switch m {
	case del:
		d.deleteChar()
	case backspace:
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
