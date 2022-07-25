package main

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
}

func (d *console) arrowDown() {
	if len(d.Input) == 0 {
		return
	}
	if d.ypos == len(d.Input)-1 {
		return
	}
	d.ypos++
	d.home()
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
