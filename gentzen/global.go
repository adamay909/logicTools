package gentzen

type logicalConstant string

var (
	lneg  = "-"
	lconj = "^"
	ldisj = "V"
	lcond = ">"
	luni  = "U"
	lex   = "X"
	lnec  = "["
	lpos  = "<"
)
var (
	neg   = logicalConstant(lneg)
	conj  = logicalConstant(lconj)
	disj  = logicalConstant(ldisj)
	cond  = logicalConstant(lcond)
	uni   = logicalConstant(luni)
	ex    = logicalConstant(lex)
	ident = logicalConstant("=")
	nec   = logicalConstant(lnec)
	pos   = logicalConstant(lpos)
	none  = logicalConstant(`*`)
)

const (
	mPolish     printMode = 0
	mLatex      printMode = 1
	mPlainLatex printMode = 2
	mPlainText  printMode = 3
	mSimple     printMode = 4
	mEnglish    printMode = 5
)

var brackets = [][2]string{
	{"", ""},
	{"(", ")"},
	{"[", "]"},
	{`\{`, `\}`},
	//	{`\big(`, `\big)`},
	//	{`\big[`, `\big]`},
	//	{`\big\{`, `\big\}`},
	{`\Big(`, `\Big)`},
	{`\Big[`, `\Big]`},
	{`\Big\{`, `\Big\}`},
	//	{`\bigg(`, `\bigg)`},
	//	{`\bigg[`, `\bigg]`},
	//	{`\bigg\{`, `\bigg\}`},
	{`\Bigg(`, `\Bigg)`},
	{`\Bigg[`, `\Bigg]`},
	{`\Bigg\{`, `\Bigg\}`},
	{`\Bigg\langle`, `\Bigg\rangle`},
}

var textBrackets = [][2]string{
	{`$\langle$`, `$\rangle$`},
	{`$\big\langle$`, `$\big\rangle$`},
	{`$\Big\langle$`, `$\Big\rangle$`},
	{`$\bigg\langle$`, `$\bigg\rangle$`},
	{`$\Bigg\langle$`, `$\Bigg\rangle$`},
}

/*
	var connectivesSL = [][]string{
		{string(neg), `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that "},
		{string(conj), `\land `, `\wedge `, "\u2227", "\u2227", " and "},
		{string(disj), `\lor `, `\vee `, "\u2228", "\u2228", " or "},
		{string(cond), `\limplies `, `\supset `, "\u2283", "\u2283", " if, then "},
	}

	var connectivesPL = [][]string{
		{string(uni), `\lforall `, `\forall `, "\u2200", "\u2200", " for all "},
		{string(ex), `\lthereis `, `\exists `, "\u2203", "\u2203", " there is a "},
		{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`, `=`, " equals "},
	}
*/
var connectivesSL, connectivesPL, connectivesML [][6]string

var infRules = [][]string{
	{`a`, `A`, `A`, `A`},
	{`ki`, `\conjI`, `\conjI`, "\u2227I"},
	{`ke`, `\conjE`, `\conjE`, "\u2227E"},
	{`di`, `\disjI`, `\disjI`, "\u2228I"},
	{`de`, `\disjE`, `\disjE`, "\u2228E"},
	{`ni`, `\negI`, `\negI`, "\u00acI"},
	{`ne`, `\negE`, `\negE`, "\u00acE"},
	{`ci`, `\condI`, `\condI`, "\u2283I"},
	{`ce`, `\condE`, `\condE`, "\u2283E"},
	{`ui`, `\uniI`, `\uniI`, "\u2200I"},
	{`ue`, `\uniE`, `\uniE`, "\u2200E"},
	{`ei`, `\exI`, `\exI`, "\u2203I"},
	{`ee`, `\exE`, `\exE`, "\u2203E"},
	{`=i`, `\iI`, `\iI`, `=I`},
	{`=e`, `\iE`, `\iE`, `=E`},
	{li, `\necI`, `\necI`, "\u25a1I"},
	{mli, `M\necI`, `M\necI`, "M\u25a1I"},
	{pli, `P\necI`, `P\necI`, "P\u25a1I"},
	{le, `\necE`, `\necE`, "\u25a1E"},
	{mi, `\posI`, `\posI`, "\u25c7I"},
	{me, `\posE`, `\posE`, "\u25c7E"},
	{mme, `M\posE`, `M\posE`, "M\u25c7E"},
	{sc, "SC", "SC", "SC"},
	{sl, "SL", "SL", "SL"},
}

var greekUCBindings = [][3]string{
	[3]string{`/G`, `\Gamma`, "\u0393"},
	[3]string{`/D`, `\Delta`, "\u0394"},
	[3]string{`/T`, `\Theta`, "\u0398"},
	[3]string{`/L`, `\Lambda`, "\u039b"},
	[3]string{`/X`, `\Xi`, "\u039e"},
	[3]string{`/P`, `\Pi`, "\u03a0"},
	[3]string{`/R`, `\Rho`, "\u03a1"},
	[3]string{`/S`, `\Sigma`, "\u03a3"},
	[3]string{`/U`, `\Upsilon`, "\u03a5"},
	[3]string{`/F`, `\Phi`, "\u03a6"},
	[3]string{`/Q`, `\Psi`, "\u03a8"},
	[3]string{`/W`, `\Omega`, "\u03a9"},
	[3]string{`/W`, `\Omega`, "\u03a9"},
	[3]string{`/W`, `\Omega`, "\u03a9"},
	[3]string{`\0`, `\emptyset`, "\u2300"},
}

var greekLCBindings = [][3]string{
	[3]string{`/a`, `\alpha`, "\u03b1"},
	[3]string{`/b`, `\beta`, "\u03b2"},
	[3]string{`/g`, `\gamma`, "\u03b3"},
	[3]string{`/d`, `\delta`, "\u03b4"},
	[3]string{`/e`, `\epsilon`, "\u03b5"},
	[3]string{`/z`, `\zeta`, "\u03b6"},
	[3]string{`/h`, `\eta`, "\u03b7"},
	[3]string{`/t`, `\theta`, "\u03b8"},
	[3]string{`/i`, `\iota`, "\u03b9"},
	[3]string{`/k`, `\kappa`, "\u03ba"},
	[3]string{`/l`, `\lambda`, "\u03bb"},
	[3]string{`/m`, `\mu`, "\u03bc"},
	[3]string{`/n`, `\nu`, "\u03bd"},
	[3]string{`/x`, `\xi`, "\u03be"},
	[3]string{`/o`, `\omicron`, "\u03bf"},
	[3]string{`/p`, `\pi`, "\u03c0"},
	[3]string{`/r`, `\rho`, "\u03c1"},
	[3]string{`/s`, `\sigma`, "\u03c3"},
	[3]string{`/y`, `\tau`, "\u03c4"},
	[3]string{`/u`, `\upsilon`, "\u03c5"},
	[3]string{`/f`, `\varphi`, "\u03c6"},
	[3]string{`/c`, `\chi`, "\u03c7"},
	[3]string{`/q`, `\psi`, "\u03c8"},
	[3]string{`/w`, `\omega`, "\u03c9"},
}

var greekUpperCaseLetters = []string{
	`\Gamma`,
	`\Delta`,
	`\Theta`,
	`\Lambda`,
	`\Xi`,
	`\Pi`,
	`\Rho`,
	`\Sigma`,
	`\Upsilon`,
	`\Phi`,
	`\Psi`,
	`\Omega`,
}

var greekLowerCaseLetters = []string{
	`\alpha`,
	`\beta`,
	`\gamma`,
	`\delta`,
	`\epsilon`,
	`\zeta`,
	`\eta`,
	`\theta`,
	`\iota`,
	`\kappa`,
	`\lambda`,
	`\mu`,
	`\nu`,
	`\xi`,
	`\omicron`,
	`\pi`,
	`\rho`,
	`\sigma`,
	`\tau`,
	`\upsilon`,
	`\varphi`,
	`\chi`,
	`\psi`,
	`\omega`,
}

var connectives [][6]string

// SetStandardPolish sets whether to use more standard notations for the
// logical constants.
func SetStandardPolish(v bool) {

	if v {

		lneg = "N"
		lconj = "K"
		ldisj = "A"
		lcond = "C"
		luni = "U"
		lex = "X"
		lnec = "L"
		lpos = "M"
	} else {

		lneg = "-"
		lconj = "^"
		ldisj = "V"
		lcond = ">"
		luni = "U"
		lex = "X"
		lnec = "["
		lpos = "<"
	}

	neg = logicalConstant(lneg)
	conj = logicalConstant(lconj)
	disj = logicalConstant(ldisj)
	cond = logicalConstant(lcond)
	uni = logicalConstant(luni)
	ex = logicalConstant(lex)
	nec = logicalConstant(lnec)
	pos = logicalConstant(lpos)
	ident = logicalConstant("=")
	none = logicalConstant(`*`)

	connectivesSL = [][6]string{
		{string(neg), `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that "},
		{string(conj), `\land `, `\wedge `, "\u2227", "\u2227", " and "},
		{string(disj), `\lor `, `\vee `, "\u2228", "\u2228", " or "},
		{string(cond), `\limplies `, `\supset `, "\u2283", "\u2283", " if , then "},
	}

	connectivesPL = [][6]string{
		{string(uni), `\lforall `, `\forall `, "\u2200", "\u2200", " for all "},
		{string(ex), `\lthereis `, `\exists `, "\u2203", "\u2203", " there is a "},
		{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`, `=`, " equals "},
	}

	connectivesML = [][6]string{
		{string(pos), `\lpos `, `\Diamond `, "\u25c7", "\u25c7", " possibly "},
		{string(nec), `\lnec `, `\Box `, "\u25a1", "\u25a1", " necessarily "},
	}

	connectives = append(connectivesSL, connectivesPL...)
	connectives = append(connectives, connectivesML...)

	return
}

// SetStandardPolish sets whether to use more standard notations for the
// logical constants.
func _SetStandardPolish(v bool) {

	if v {

		lneg = "N"
		lconj = "K"
		ldisj = "A"
		lcond = "C"
		luni = "U"
		lex = "X"
		lnec = "L"
		lpos = "M"

		neg = logicalConstant(lneg)
		conj = logicalConstant(lconj)
		disj = logicalConstant(ldisj)
		cond = logicalConstant(lcond)
		uni = logicalConstant(luni)
		ex = logicalConstant(lex)
		nec = logicalConstant(lnec)
		pos = logicalConstant(lpos)
		ident = logicalConstant("=")
		none = logicalConstant(`*`)

		connectivesSL = [][6]string{
			{string(neg), `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that "},
			{string(conj), `\land `, `\wedge `, "\u2227", "\u2227", " and "},
			{string(disj), `\lor `, `\vee `, "\u2228", "\u2228", " or "},
			{string(cond), `\limplies `, `\supset `, "\u2283", "\u2283", " if , then "},
		}

		connectivesPL = [][6]string{
			{string(uni), `\lforall `, `\forall `, "\u2200", "\u2200", " for all "},
			{string(ex), `\lthereis `, `\exists `, "\u2203", "\u2203", " there is a "},
			{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`, `=`, " equals "},
		}
		connectives = append(connectivesSL, connectivesPL...)

		return
	}

	lneg = "-"
	lconj = "^"
	ldisj = "V"
	lcond = ">"
	luni = "U"
	lex = "X"
	lnec = "["
	lpos = "<"

	neg = logicalConstant(lneg)
	conj = logicalConstant(lconj)
	disj = logicalConstant(ldisj)
	cond = logicalConstant(lcond)
	uni = logicalConstant(luni)
	ex = logicalConstant(lex)
	ident = logicalConstant("=")
	none = logicalConstant(`*`)

	connectivesSL = [][6]string{
		{string(neg), `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that "},
		{string(conj), `\land `, `\wedge `, "\u2227", "\u2227", " and "},
		{string(disj), `\lor `, `\vee `, "\u2228", "\u2228", " or "},
		{string(cond), `\limplies `, `\supset `, "\u2283", "\u2283", " if , then "},
	}

	connectivesPL = [][6]string{
		{string(uni), `\lforall `, `\forall `, "\u2200", "\u2200", " for all "},
		{string(ex), `\lthereis `, `\exists `, "\u2203", "\u2203", " there is a "},
		{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`, `=`, " equals "},
	}

	connectives = append(connectivesSL, connectivesPL...)
	return
}

func greekCharOf(s string) string {

	for _, e := range greekLCBindings {
		if e[1] == s {
			return e[2]
		}
	}
	for _, e := range greekUCBindings {
		if e[1] == s {
			return e[2]
		}
	}
	return s
}

var prettifyBrackets = true

func SetPrettyBrackets(v bool) {
	prettifyBrackets = v
}
