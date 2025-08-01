package gentzen

/**************************************
The variables defined here should be treated as constants except by functins defined in this
file.
*************************************/

// LogicalConstant represents the logical constants.
type LogicalConstant int

var (
	lneg   = "-"
	lconj  = "^"
	ldisj  = "V"
	lcond  = ">"
	luni   = "U"
	lex    = "X"
	lnec   = "["
	lpos   = "<"
	lident = "="
)

//go:generate stringer -type LogicalConstant
const (
	None  LogicalConstant = iota //no connective
	Neg                          //negation
	Conj                         //conjunction
	Disj                         //disjunction
	Cond                         //conditional
	Uni                          //universal quantifier
	Ex                           //existential quantifier
	Ident                        //identity
	Nec                          //necessity
	Pos                          //possibility
)

//go:generate stringer -type=PrintMode
const (
	O_Polish       PrintMode = iota //Polish notation
	O_Latex                         //LaTeX code
	O_PlainLatex                    //ditto
	O_PlainText                     //plain text
	O_Simple                        //unused
	O_English                       //plain English
	O_PlainASCII                    //plain text restricted to ASCII
	O_ProofChecker                  //for proof checker
)

var brackets = [][2]string{
	{"", ""},
	{`\big(`, `\big)`},
	{`\big[`, `\big]`},
	{`\big\{`, `\big\}`},
	{`\Big(`, `\Big)`},
	{`\Big[`, `\Big]`},
	{`\Big\{`, `\Big\}`},
	{`\bigg(`, `\bigg)`},
	{`\bigg[`, `\bigg]`},
	{`\bigg\{`, `\bigg\}`},
	{`\Bigg(`, `\Bigg)`},
	{`\Bigg[`, `\Bigg]`},
	{`\Bigg\{`, `\Bigg\}`},
	{`\Bigg\langle`, `\Bigg\rangle`},
}

var textBrackets = [][2]string{
	{"", ""},
	{`$\langle$`, `$\rangle$`},
	{`$\big\langle$`, `$\big\rangle$`},
	{`$\Big\langle$`, `$\Big\rangle$`},
	{`$\bigg\langle$`, `$\bigg\rangle$`},
	{`$\Bigg\langle$`, `$\Bigg\rangle$`},
}

var proofCheckerBrackets = [][2]string{
	{"", ""},
	{`<span class="big">(</span>`, `<span class="big">)</span>`},
	{`<span class="big">[</span>`, `<span class="big">]</span>`},
	{`<span class="big">{</span>`, `<span class="big">}</span>`},
	{`<span class="bigg">(</span>`, `<span class="bigg">)</span>`},
	{`<span class="bigg">[</span>`, `<span class="bigg">]</span>`},
	{`<span class="bigg">{</span>`, `<span class="bigg">}</span>`},
	{`<span class="biggg">(</span>`, `<span class="biggg">)</span>`},
	{`<span class="biggg">[</span>`, `<span class="biggg">]</span>`},
	{`<span class="biggg">{</span>`, `<span class="biggg">}</span>`},
}

var simpleBrackets = [][2]string{
	{"", ""},
	{"(", ")"},
	{"[", "]"},
	{`{`, `}`},
}

var plainBrackets = simpleBrackets[:2]

var connectivesSL, connectivesPL, connectivesML [][7]string

var infRules = [][]string{
	{`a`, `A`, `A`, `A`},
	{`m`, `M`, `M`, `M`},
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
	{mli, `S5\necI`, `S5\necI`, "S5\u25a1I"},
	{pli, `S4\necI`, `S4\necI`, "S4\u25a1I"},
	{tli, `T\necI`, `T\necI`, "T\u25a1I"},
	{le, `\necE`, `\necE`, "\u25a1E"},
	{mi, `\posI`, `\posI`, "\u25c7I"},
	{me, `\posE`, `\posE`, "\u25c7E"},
	{mme, `S5\posE`, `S5\posE`, "S5\u25c7E"},
	{sc, "SC", "SC", "SC"},
	{sl, "logic", "logic", "logic"},
	{"rewrite", "", "", ""},
	{`\lposDR`, `\lposDR`, `\lposDR`, "\u25c7DR"},
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
	//	[3]string{`/W`, `\Omega`, "\u03a9"},
	//	[3]string{`/W`, `\Omega`, "\u03a9"},
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
	[3]string{`/f`, `\phi`, "\u03c6"},
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

var connectives [][7]string

// SetStandardPolish sets whether to use more standard notations for the
// logical constants.
func SetStandardPolish(v bool) {

	oSYMB = !v

	setupConnectives()

	return
}

func setupConnectives() {

	if !oSYMB {

		lneg = "N"
		lconj = "K"
		ldisj = "A"
		lcond = "C"
		luni = "U"
		lex = "X"
		lnec = "[" //don't use L and M
		lpos = "<"
		lident = "="
	} else {

		lneg = "-"
		lconj = "^"
		ldisj = "V"
		lcond = ">"
		luni = "U"
		lex = "X"
		lnec = "["
		lpos = "<"
		lident = "="
	}
	/*
		Neg = LogicalConstant(lneg)
		Conj = LogicalConstant(lconj)
		Disj = LogicalConstant(ldisj)
		Cond = LogicalConstant(lcond)
		Uni = LogicalConstant(luni)
		Ex = LogicalConstant(lex)
		Nec = LogicalConstant(lnec)
		Pos = LogicalConstant(lpos)
		Ident = LogicalConstant("=")
		None = LogicalConstant(`*`)
	*/

	connectivesSL = nil
	connectivesPL = nil
	connectivesML = nil

	connectivesSL = [][7]string{
		{lneg, `\lnot `, `\neg `, "\u00ac", "\u00ac", " it is not the case that ", "-"},
		{lconj, `\land `, `\wedge `, "\u2227", "\u2227", " and ", "^"},
		{ldisj, `\lor `, `\vee `, "\u2228", "\u2228", " or ", "v"},
	}

	connectivesPL = [][7]string{
		{luni, `\lforall `, `\forall `, "\u2200", "\u2200", " for all ", "U"},
		{lex, `\lthereis `, `\exists `, "\u2203", "\u2203", " there is a ", "X"},
		{lident, `\mathbin{=}`, `\mathbin{=}`, `=`, `=`, " equals ", "="},
	}

	connectivesML = [][7]string{
		{lpos, `\lpos `, `\Diamond `, "\u25c7", "\u25c7", " possibly ",
			"<"},
		{lnec, `\lnec `, `\Box `, "\u25fb", "\u25fb", " necessarily ", "["},
	}

	if oCOND {
		connectivesSL = append(connectivesSL, [7]string{lcond, `\limplies `, `\supset `, "\u2283", "\u2283", " if , then ", ">"})
	}

	connectives = append(connectivesSL, connectivesPL...)
	connectives = append(connectives, connectivesML...)

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

// SetPrettyBrackets sets whether to use variations in brackets styles (form and possibly size) for improved readability.
func SetPrettyBrackets(v bool) {
	prettifyBrackets = v
}

func codeOf(l LogicalConstant) string {

	switch l {
	case Neg:
		return lneg
	case Conj:
		return lconj
	case Disj:
		return ldisj
	case Cond:
		return lcond
	case Uni:
		return luni
	case Ex:
		return lex
	case Nec:
		return lnec
	case Pos:
		return lpos
	default:
		return ""
	}

}
