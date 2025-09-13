package gentzen

/**************************************
The variables defined here should be treated as constants except by functions defined in this
file.
*************************************/

// logicalConstant represents the logical constants.
type logicalConstant int

var (
	lneg   = "N"
	lconj  = "K"
	ldisj  = "A"
	lcond  = "C"
	luni   = "U"
	lex    = "X"
	lident = "="
)

//go:generate stringer -type logicalConstant
const (
	noConstant logicalConstant = iota //no connective
	neg                               //negation
	conj                              //conjunction
	disj                              //disjunction
	cond                              //conditional
	uni                               //universal quantifier
	ex                                //existential quantifier
	ident                             //identity
	nec                               //necessity
	pos                               //possibility
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
	O_ProofChecker                  //for proof checker [../proofChecker]
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

func setupConnectives() {

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

func codeOf(l logicalConstant) string {

	switch l {
	case neg:
		return lneg
	case conj:
		return lconj
	case disj:
		return ldisj
	case cond:
		return lcond
	case uni:
		return luni
	case ex:
		return lex
	default:
		return ""
	}

}

/*
SetSpecialConn directs gentzen to use non-standard connectives for Polish notation to free some characters (if v is true):

negation: -
conjunction: ^
disjunction: v
conditional: >

This also affects the processing of inference rules so that you cannot use the standard rule declared through [SetStandardInferenceRules].
You will have to redefine the inference rules.
*/
func SetSpecialConn(v bool) {

	if v {
		lneg = "-"
		lconj = "^"
		ldisj = "v"
		lcond = ">"
		luni = "U"
		lex = "X"
		lident = "="
	} else {
		lneg = "N"
		lconj = "K"
		ldisj = "A"
		lcond = "C"
		luni = "U"
		lex = "X"
		lident = "="
	}
	setupConnectives()
}
