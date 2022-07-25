package main

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
