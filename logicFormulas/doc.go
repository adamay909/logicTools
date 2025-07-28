/*
LogicFormulas is a tool for generating random sentences of sentential logic. It is intended for use with my [logic text book] (to generate formulas for exercises and the like).

# Some examples.

Generate 100 random sentences:

	$logicFormulas -random -n 10

This will generate 10 wffs of sentential logic using at most 10 distinct atomic sentences and no restrictions on complexity. The output is formatted in the Polish notation. Here is a possible output:

	ANCNNNNANNKALNMRNCSSWP
	KMNP
	NNCNPR
	CSKKNQPP
	KMCSH
	CNKRCKGKKMMKWKLNSNACKPGCNMAKCKCPLCCPQMMCALCRPAKNLNRPNAAANGGLSCMCCAAQANNMNKMKNAPMLANKSCCSNRAACCSPASNRGGNKPKNLNKWSCALQSKGCCAAKACLNLANKLGSRPPKNANRNGNCNAWSMNKCMCCKKGNRNKRSCAKAKAAGMLRCLARKKAGKANPLQGMLQMNMNNAQLAKCRRGL
	AQP
	KSACWQAKCSCKAKCWNSNNKWNAANCCQWSWACACKSKKCRCSKNWKSNPSCKSNKSKWKWPPCCPWCNWQCNRCAPWKPAPAPKSKSCAAKPKNPWPRANKCCCPQKCSKPCWSQSNNNRNKNQCQCNNANNQSRPQQWPRS
	CRL
	AWS

You might notice that some outputs differ only by the sentence letters used. You can exclude such duplication with the -uniqueS flag:

	$logicFormulas -random -uniqueS -n 10

Generate 10 random sentences with maximum class of 3 (see text book for more details on what this means):

	$ logicFormulas -random -n 10 -maxClass 3
	NNL
	KKNSWP
	AQG
	NCNPP
	CNCPRW
	KKQRW
	NP
	KCRMCWQ
	KWANPCPP
	CKPKPPP

Restrict the number of atomic sentences to just 2:

	$ logicFormulas -random -n 10 -maxAtomic 2
	NAAKAPPAKCNKPPPPPPKAANKPAAPPPNPCPANPNAAPPKPPP
	ANPP
	KPCPP
	NP
	CNPKQQ
	ACCNAAPCPAQQAPPNPPAPQ
	CPP
	KNPNP
	CCNNPNPKPQ
	CCPKQCCPQQQ

You can also generate a bunch of tautologies:

	$ logicFormulas -tautology -n 10 -maxAtomic 2
	APNP
	CPP
	CPCANPPP
	AANANNKPCANPKPCPCNPPPKKCNAPPPPPPP
	AKPKAQCQQNQAKNPCPPAKPPP
	ANNPCPP
	CKQQQ
	CKKQCCCQCQKKQPPQPNQP
	CCNPPKKAPKPPPP
	ANPP

You can direct the output to a file:

	$ logicFormulas -random -n 100000 -dest list

This will write 100000 random sentences to the file called "list".

If for some reason you don't like the Polish notation, you can convert the output to a more standard infix notation by piping the output to [lbhelper]:

	$ logicFormulas -random -n 5 | lbhelper -txt
	¬{{{(Q∨P)∨{[(Q∨Q)∧P]⊃{[S⊃(P⊃P)]⊃P}}}∨R}∨(S∧P)}
	¬¬P
	¬{[¬H⊃(L∨M)]⊃W}
	{[(P∧P)∧P]⊃P}∨P
	G⊃¬L

Of course, you can convert to other formats, too, or do things like generate truth tables for random sentences. See the documentation for [lbhelper] for more details.

# Complete list of flags

	$ logicFormulas -h
	Usage of logicFormulas:
	  -dest string
	    	output destination; default is stdout
	  -malform
	    	like random but about half are malformed. Output is LaTeX.
	  -maxAtomic int
	    	max number of distinct atomic sentences to use (default 10)
	  -maxClass int
	    	max syntax tree height (class) of generated sentence. Nagative number means no limit (default -1)
	  -n int
	    	number of formulas to generate (default 1)
	  -normalize
	    	normalize sentence letters, etc.
	  -random
	    	generate n random sentences with m atomic sentences and at most d connectives
	  -tautology
	    	extract tautologies from f
	  -uniqueS
	    	treat sentences with same structure as the same
	  -withC
	    	with conditional (default true)



[logic text book]: https://github.com/adamay909/logicbook

[lbhelper]: https://github.com/adamay909/logicTools/tree/main/lbhelper

*/
//go:generate pkgdoc2readme
package main
