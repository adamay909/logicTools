package main

type inputLine []string

type console struct {
	Title               string
	Input               []inputLine
	SystemPL            bool
	Theorems            bool
	xpos, ypos          int
	xprev, yprev        int
	xcursor, ycursor    int
	html                []string
	modifier            string
	overhang            bool
	Offset              int
	viewTop, viewBottom int
	fontSize            int
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

func (d *console) clear() {
	d.Input = nil
	d.html = nil
	d.xpos = 0
	d.ypos = 0
	d.xprev, d.yprev = 0, 0
	d.xcursor, d.ycursor = 0, 0
	d.html = nil
	d.viewTop = 0
	d.viewBottom = 25
	d.overhang = true
	d.modifier = ""
	d.Offset = 1
	if !oExercises {
		d.Title = ""
	}

}

func (d *console) reset() {
	d.clear()
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
