package gentzen

type logicalConstant string

var (
	lneg  = "-"
	lconj = "^"
	ldisj = "V"
	lcond = ">"
	luni  = "U"
	lex   = "X"
)
var (
	neg   = logicalConstant(lneg)
	conj  = logicalConstant(lconj)
	disj  = logicalConstant(ldisj)
	cond  = logicalConstant(lcond)
	uni   = logicalConstant(luni)
	ex    = logicalConstant(lex)
	ident = logicalConstant("=")
	none  = logicalConstant(`*`)
)

const (
	mPolish     printMode = 0
	mLatex      printMode = 1
	mPlainLatex printMode = 2
	mPlainText  printMode = 3
	mSimple     printMode = 4
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

var connectivesSL = [][]string{
	{string(neg), `\lnot `, `\neg `, "\u00ac"},
	{string(conj), `\land `, `\wedge `, "\u2227"},
	{string(disj), `\lor `, `\vee `, "\u2228"},
	{string(cond), `\limplies `, `\supset `, "\u2283"},
}

var connectivesPL = [][]string{
	{string(uni), `\lforall `, `\forall `, "\u2200"},
	{string(ex), `\lthereis `, `\exists `, "\u2203"},
	{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`},
}

var connectives [][]string

//SetStandardPolish sets whether to use more standard notations for the
//logical constants.
func SetStandardPolish(v bool) {

	if v {

		lneg = "N"
		lconj = "K"
		ldisj = "A"
		lcond = "C"
		luni = "U"
		lex = "X"

		neg = logicalConstant(lneg)
		conj = logicalConstant(lconj)
		disj = logicalConstant(ldisj)
		cond = logicalConstant(lcond)
		uni = logicalConstant(luni)
		ex = logicalConstant(lex)
		ident = logicalConstant("=")
		none = logicalConstant(`*`)

		connectivesSL = [][]string{
			{string(neg), `\lnot `, `\neg `, "\u00ac"},
			{string(conj), `\land `, `\wedge `, "\u2227"},
			{string(disj), `\lor `, `\vee `, "\u2228"},
			{string(cond), `\limplies `, `\supset `, "\u2283"},
		}

		connectivesPL = [][]string{
			{string(uni), `\lforall `, `\forall `, "\u2200"},
			{string(ex), `\lthereis `, `\exists `, "\u2203"},
			{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`},
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

	neg = logicalConstant(lneg)
	conj = logicalConstant(lconj)
	disj = logicalConstant(ldisj)
	cond = logicalConstant(lcond)
	uni = logicalConstant(luni)
	ex = logicalConstant(lex)
	ident = logicalConstant("=")
	none = logicalConstant(`*`)

	connectivesSL = [][]string{
		{string(neg), `\lnot `, `\neg `, "\u00ac"},
		{string(conj), `\land `, `\wedge `, "\u2227"},
		{string(disj), `\lor `, `\vee `, "\u2228"},
		{string(cond), `\limplies `, `\supset `, "\u2283"},
	}

	connectivesPL = [][]string{
		{string(uni), `\lforall `, `\forall `, "\u2200"},
		{string(ex), `\lthereis `, `\exists `, "\u2203"},
		{string(ident), `\mathbin{=}`, `\mathbin{=}`, `=`},
	}

	connectives = append(connectivesSL, connectivesPL...)
	return
}
