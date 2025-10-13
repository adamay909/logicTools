vim9script noclear
# Filetype plugin for helping writing logic formulas in  a LaTeX file
# By Masahiro Yamada

if exists("g:lbh_donotload") 
 finish
endif

if exists("loaded")
 finish
endif

var loaded = true

if !exists("g:lbh_program")
 g:lbh_program = "lbhelper -vim"
endif

if !exists("g:lbh_delimstart")
 g:lbh_delimstart = ";;"
endif

if !exists("g:lbh_delimend")
 g:lbh_delimend = ";"
endif

var searchpat = g:lbh_delimstart .. '\(.\{-}\)' .. g:lbh_delimend

var converter = job_start(g:lbh_program)

var convchan = job_getchannel(converter)
 
def ExpandFunc()
  var mlist = matchbufline("%", searchpat, line("."), line("."))
 	for m in mlist  
	 var ms = m['text']
  	 var startpos = m['byteidx']
   	 var rawrepl = ch_evalraw(convchan, substitute(m['text'], searchpat, '\1', "") .. "\n")
	 if rawrepl != ""
      setline('.', substitute(getline('.'), m['text'], escape(rawrepl, '\'), ""))	
      cursor('.', m['byteidx'] + len(rawrepl) + 1)
     endif
	endfor
enddef

def SetupBuffer() 
	if &filetype != "tex" 
	 return
	endif	
 augroup InlineExpand
  autocmd TextChangedI,TextChanged <buffer> call ExpandFunc()
augroup END
enddef
 
augroup InlineExpand
  autocmd!
  autocmd BufEnter,FileType * call SetupBuffer() 
  #autocmd TextChangedI,TextChanged <buffer> call ExpandFunc()
augroup END



