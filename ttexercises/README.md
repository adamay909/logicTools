## Overview

Ttexercises is a tool for generating truth table exercises to go with my [logic text book](https://github.com/adamay909/logicbook). There are two types of exercises: fill in the blanks, and detect errors.

For fill-in-the-blanks type exercises, you can choose to hide the values of some rows, columns, or cells.

E.g, to hide up to 3 rows:

	$ttexercises -delR -n 3 CApqKpq

will hide up to 3 rows (exactly how many and which is chosen using a random number generator) of the truth table for CApqKpq ((p∨q)⊃(p∧q)). You can hide up to 3 columns with:

	$ttexercises -delC -n 3 CApqKpq

And hide up to 6 cells with:

	$ttexercises -del -n 6 CApqKpq

You can get a truth table that only has the interpretations by:

	$ttexercises -empty CApqKpq

You can make the table completely empty with:

	$ttexercises -empty -val CApqKpq

You can introduce random errors with:

	$ttexercises -errors -n 4 CApqKpq

By default, the output is Stdout but you probably want to direct the output to a file. You can do that with the -dest flag:

	$ttexercises -dest exercise.tex -del -n 8 CApqKpq

The output is always in LaTeX with commands defined in the style files for the text book. It's coded so that you can use [lbhelper](https://github.com/adamay909/logicTools/tree/main/lbhelper) to compile the exercise version as well as the answer key version from a single file. For instance, assuming you saved the output to exericise.tex

	$lbhelper -compile exercise.tex

will produce exercise.pdf which has the table with blanks.

	$lbhelper -compile -answer exercise.tex

will produce exercise\_answers.pdf which has the table with blanks filled in and highlighted. You do need to make sure you have the relevant styles installed in the appropriate places. See the documentation for [lbhelper](https://github.com/adamay909/logicTools/tree/main/lbhelper) for more information on that.

You probably also don't want to input the relevant formulas by hand. You can read from a file using the flag -inputfile. E.g.:

	$ttexercises -inputfile list -dest exercise.tex -errors -n 8

You could generate the list of sentences using [logicFormulas](https://github.com/adamay909/logicTools/tree/main/logicFormulas). The resulting LaTeX might want some extra editing like inserting question numbers but that's trivial work that any decent text editor can accomplish in just a few steps (each individual table starts with a line that begins with %%% followed by a single space and then the formula in the Polish notation; it's LaTeX so you can use the enumerate environment). This way, you can easily generate hundreds of thousands of truth table exercises in no time.

The generated truth tables are 'wide' ones that I prefer and use in the text book. You can switch to 'narrow' truth tables with the -narrow flag.

### Complete list of flags

	$ ttexercises -h
	Usage of ttexercises:
	  -correct
	    	generate full truth table
	  -del
	    	clear cell -n number of cells
	  -delC
	    	clear column; -n columns
	  -delR
	    	clear row -n rows
	  -delXY
	    	clear up to n rows or columns
	  -dest string
	    	output destination (default is Stdout)
	  -empty
	    	generate truth table with just columns and valuations
	  -errors
	    	insert errors
	  -inputfile string
	    	read input from the named file
	  -n int
	    	number
	  -narrow
	    	use narrow format truth table
	  -val
	    	include errors in interpretations/also hide interpretations


