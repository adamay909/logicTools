## lbhelpervim

This is a vim filetype plugin for helping with editing LaTeX for logic related material like my [logic textbook](https://github.com/adamay909/logicbook) but it could be adapted for other purposes as it is just a plugin for replacing certain strings with other strings (see the options below).

### Installing the Plugin

The Makefile will install the plugin as a package in the folder `$HOME/.vim/pack/lbhelpervim`. If you are not on Linux, or have a different folder for vim to look for its stuff, you need to adjust accordingly. 

For the plugin to work with its default settings, you need to install lbhelper in your path. You need a version that accepts the "-vim" option. When in doubt, reinstall it.

### Default Behavior

The default working is as follows:
Given you type a string of the form `;;...;`, the `...` part is sent to lbhelper for parsing and if successful, `;;...;` is replaced with a LaTeX encoded version of `...` . For example, `;;Cpq;` is replaced with `\p{p\mc{\limplies }q}`. If conversion does not succeed, the `;;...;` is left alone. 

Sequents and entailments are also detected: a ":" or "|-" are interpreted as the turnstile and appropriate conversion will be attempted. A "|=" is interpreted as the double turnstile. 

You do not need to specify the format of the `...` part. The conversion tries conversion as sentential logic in polish format, sentential logic in infix format, predicate logic in polish format, predicate logic in infix format, in that order. As soon as a conversion succeeds, that is returned as output. 

### Options

You can set some options in your vimrc file through the following global variables:

    g:lbh_program 
	    This controls the program to do the actual conversion. Default 
	    is "lbhelper -vim". But you can set it to any program provided: 
	    it listens on stdin for input, and outputs its response to stdout. 
	    Each input and output are delimited by newlines. Conversion failure
        returns an empty string. The program is started only once for
        each buffer so it should not time out waiting for input.

    g:lbh_delimstart 
	    This controls the delimiter for the start of the string. Default 
        is ";;".

    g:lbh_delimend = ";"
	    This controls the delimiter for the end of the string. Default is 
        ";".

    g:lbh_donotload = true
	    This keeps the plugin from loading. Actually, just setting this 
        variable to anything will disable the plugin.

This plugin only looks at the current line to do its conversion, so if you let vim hard-wrap your lines you might need to adjust where vim is allowed to insert line breaks with the `breakat` option.
