## Overview

Package gentzen provides some tools for parsing, serializing, and checking proofs (both porpositional logic and first order predicate logic with identity), printing truth tables, syntax trees (propositional logic), semantic tableaux, and a few other things.

You need the commands defined in logic\_commands.sty (included with the source files for the textbook available at [https://github.com/adamay909/logicbook](https://github.com/adamay909/logicbook).

The main entry points expect inputs as plain strings in the Polish notation. For historical reasons, the default for the logical constants is non-standard:

	negation: -
	conjunction: ^
	disjunction: v
	conditional: >
	universal quantifier: U
	existential quantifier X

You can switch to a more standard Polish notation with

	SetStandardPolish(true)

which will switch the notation to:

	negation: N
	conjunction: K
	disjunction: A
	conditional: C
	universal quantifier: U
	existential quantifier X

Either way, some letters are reserved for logical constants and are therefore not allowed for use as sentence or predicate letters.

If you need a bit more information on the Polish notation, see the [SEP entry on the Polish notation](https://plato.stanford.edu/entries/lukasiewicz/polish-notation.html). Once you get used to it, your quality of life will improve massively when it comes to typing material on logic.

For the language of first-order logic, terms must be lower case Roman alphabet letters. Lower case Greek letters can be used to stand for formulas. Both

	Ux/f

	Ux/f(x)

are accepted as input. /f stands for lower case phi. More generally, Greek letters can be entered using:

  - /G for Gamma
  - /D for Delta
  - /T for Theta
  - /L for Lambda
  - /X for Xi
  - /P for Pi
  - /R for Rho
  - /S for Sigma
  - /U for Upsilon
  - /F for Phi
  - /Q for Psi
  - /W for Omega
  - /a for alpha
  - /b for beta
  - /g for gamma
  - /d for delta
  - /e for epsilon
  - /z for zeta
  - /h for eta
  - /t for theta
  - /i for iota
  - /k for kappa
  - /l for lambda
  - /m for mu
  - /n for nu
  - /x for xi
  - /o for omicron
  - /p for pi
  - /r for rho
  - /s for sigma
  - /y for tau
  - /u for upsilon
  - /f for phi
  - /c for chi
  - /q for psi
  - /w for omega

Subscripts can be entered by prefixing with underscore (\_). Subscripts can be a string of Arabic numerals or else a single letter of the Roman or Greek alphabet. Subscripts cannot have further subscripts.

There is a parser for an infix notation. The input looks like the input for the proof checker except that the backslash should be forward slash.


