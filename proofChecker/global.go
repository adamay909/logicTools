package main

import (
	"errors"
)

type tkType int

const (
	tkraw tkType = 0
	tktex tkType = 1
	tktxt tkType = 2
)

var (
	keyBindings = [][3]string{
		[3]string{`A`, `A`, `A`},
		[3]string{`B`, `B`, `B`},
		[3]string{`C`, `C`, `C`},
		[3]string{`D`, `D`, `D`},
		[3]string{`E`, `E`, `E`},
		[3]string{`F`, `F`, `F`},
		[3]string{`G`, `G`, `G`},
		[3]string{`H`, `H`, `H`},
		[3]string{`I`, `I`, `I`},
		[3]string{`J`, `J`, `J`},
		[3]string{`K`, `K`, `K`},
		[3]string{`L`, `L`, `L`},
		[3]string{`M`, `M`, `M`},
		[3]string{`N`, `N`, `N`},
		[3]string{`O`, `O`, `O`},
		[3]string{`P`, `P`, `P`},
		[3]string{`Q`, `Q`, `Q`},
		[3]string{`R`, `R`, `R`},
		[3]string{`S`, `S`, `S`},
		[3]string{`T`, `T`, `T`},
		[3]string{`W`, `W`, `W`},
		[3]string{`Y`, `Y`, `Y`},
		[3]string{`Z`, `Z`, `Z`},
		[3]string{`a`, `a`, `a`},
		[3]string{`b`, `b`, `b`},
		[3]string{`c`, `c`, `c`},
		[3]string{`d`, `d`, `d`},
		[3]string{`e`, `e`, `e`},
		[3]string{`f`, `f`, `f`},
		[3]string{`g`, `g`, `g`},
		[3]string{`h`, `h`, `h`},
		[3]string{`i`, `i`, `i`},
		[3]string{`j`, `j`, `j`},
		[3]string{`k`, `k`, `k`},
		[3]string{`l`, `l`, `l`},
		[3]string{`m`, `m`, `m`},
		[3]string{`n`, `n`, `n`},
		[3]string{`o`, `o`, `o`},
		[3]string{`p`, `p`, `p`},
		[3]string{`q`, `q`, `q`},
		[3]string{`r`, `r`, `r`},
		[3]string{`s`, `s`, `s`},
		[3]string{`t`, `t`, `t`},
		[3]string{`u`, `u`, `u`},
		[3]string{`w`, `w`, `w`},
		[3]string{`x`, `x`, `x`},
		[3]string{`y`, `y`, `y`},
		[3]string{`z`, `z`, `z`},
		[3]string{`1`, `1`, `1`},
		[3]string{`2`, `2`, `2`},
		[3]string{`3`, `3`, `3`},
		[3]string{`4`, `4`, `4`},
		[3]string{`5`, `5`, `5`},
		[3]string{`6`, `6`, `6`},
		[3]string{`7`, `7`, `7`},
		[3]string{`8`, `8`, `8`},
		[3]string{`9`, `9`, `9`},
		[3]string{`0`, `0`, `0`},
	}

	punctBindings = [][3]string{
		[3]string{`(`, `(`, `(`},
		[3]string{`)`, `)`, `)`},
		[3]string{`,`, `,`, `,`},
	}

	connBindings = [][3]string{
		[3]string{`V`, `\vee`, "\u2228"},
		[3]string{`v`, `\vee`, "\u2228"},
		[3]string{`-`, `\neg`, "\u00ac"},
		[3]string{`^`, `\wedge`, "\u2227"},
		[3]string{`>`, `\supset`, "\u2283"},
	}

	plBindings = [][3]string{
		[3]string{`U`, `\forall`, "\u2200"},
		[3]string{`X`, `\exists`, "\u2203"},
		[3]string{"=", "=", "="},
		[3]string{`\=`, "≠", "≠"},
	}

	mlBindings = [][3]string{
		[3]string{`[`, `\lnec`, "\u25a1"},
		[3]string{`<`, `\lpos`, "\u25c7"},
	}

	turnstileBindings = [][3]string{
		[3]string{`|-`, `\vdash`, "⊢"},
	}

	dotsBindings = [][3]string{
		[3]string{`.`, `\ldots`, `...`},
	}

	greekBindings = [][3]string{
		[3]string{`\G`, `\Gamma`, "Γ"},
		[3]string{`\D`, `\Delta`, "Δ"},
		[3]string{`\T`, `\Theta`, "\u0398"},
		[3]string{`\L`, `\Lambda`, "\u039b"},
		[3]string{`\X`, `\Xi`, "\u039e"},
		[3]string{`\P`, `\Pi`, "\u03a0"},
		[3]string{`\R`, `\Rho`, "\u03a1"},
		[3]string{`\S`, `\Sigma`, "\u03a3"},
		[3]string{`\U`, `\Upsilon`, "\u03a5"},
		[3]string{`\F`, `\Phi`, "Φ"},
		[3]string{`\Q`, `\Psi`, "\u03a8"},
		[3]string{`\W`, `\Omega`, "\u03a9"},
		[3]string{`\a`, `\alpha`, "\u03b1"},
		[3]string{`\b`, `\beta`, "\u03b2"},
		[3]string{`\g`, `\gamma`, "\u03b3"},
		[3]string{`\d`, `\delta`, "\u03b4"},
		[3]string{`\e`, `\epsilon`, "\u03b5"},
		[3]string{`\z`, `\zeta`, "\u03b6"},
		[3]string{`\h`, `\eta`, "\u03b7"},
		[3]string{`\t`, `\theta`, "\u03b8"},
		[3]string{`\i`, `\iota`, "\u03b9"},
		[3]string{`\k`, `\kappa`, "\u03ba"},
		[3]string{`\l`, `\lambda`, "\u03bb"},
		[3]string{`\m`, `\mu`, "\u03bc"},
		[3]string{`\n`, `\nu`, "\u03bd"},
		[3]string{`\x`, `\xi`, "\u03be"},
		[3]string{`\o`, `\omicron`, "\u03bf"},
		[3]string{`\p`, `\pi`, "\u03c0"},
		[3]string{`\r`, `\rho`, "\u03c1"},
		[3]string{`\s`, `\sigma`, "\u03c3"},
		[3]string{`\y`, `\tau`, "\u03c4"},
		[3]string{`\u`, `\upsilon`, "\u03c5"},
		[3]string{`\f`, `\varphi`, "\u03c6"},
		[3]string{`\c`, `\chi`, "\u03c7"},
		[3]string{`\q`, `\psi`, "\u03c8"},
		[3]string{`\w`, `\omega`, "\u03c9"},
		[3]string{`\0`, `\emptyset`, "\u2300"},
	}

	extraBindings = [][3]string{
		[3]string{` `, ` `, ` `},
		[3]string{`.`, `.`, `.`},
		[3]string{`|X`, `X`, `X`},
		[3]string{`|U`, `U`, `U`},
		[3]string{`|v`, `v`, `v`},
		[3]string{`|V`, `V`, `V`},
		[3]string{`?`, `?`, `?`},
		[3]string{`!`, `!`, `!`},
		[3]string{`:`, `:`, `:`},
		[3]string{`;`, `;`, `;`},
	}
)

var allBindings [][3]string

func init() {
	allBindings = combineBindings(keyBindings, punctBindings, connBindings, plBindings, mlBindings, turnstileBindings, dotsBindings, greekBindings)
}

func tkOf(s string, srctype, dsttype tkType, b [][3]string) (tk string, err error) {

	for _, e := range b {
		if s == e[srctype] {
			tk = e[dsttype]
			return
		}
	}
	err = errors.New("Not Found")
	return
}

func textOf(s string) string {

	r, err := tkOf(s, tkraw, tktxt, allBindings)

	if err != nil {
		return s
	}

	return r
}

func combineBindings(b ...[][3]string) [][3]string {

	var r [][3]string

	for _, e := range b {
		r = append(r, e...)
	}

	return r
}
