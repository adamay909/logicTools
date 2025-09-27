package editor

var charMap map[string]string
var revMap map[string]string
var extraMap map[string]string

var permittedChars []string

const turnstile = "\u22a2"
const ldots = "\u2026"

func init() {
	s := `0123456789()^>,-\=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ` + "\n"

	for _, r := range s {
		permittedChars = append(permittedChars, string(r))
	}

	bindings := [][2]string{
		[2]string{`\G`, "\u0393"},
		[2]string{`\D`, "\u0394"},
		[2]string{`\T`, "\u0398"},
		[2]string{`\L`, "\u039b"},
		[2]string{`\X`, "\u039e"},
		[2]string{`\P`, "\u03a0"},
		[2]string{`\R`, "\u03a1"},
		[2]string{`\S`, "\u03a3"},
		[2]string{`\U`, "\u03a5"},
		[2]string{`\F`, "\u03a6"},
		[2]string{`\Q`, "\u03a8"},
		[2]string{`\W`, "\u03a9"},
		[2]string{`\a`, "\u03b1"},
		[2]string{`\b`, "\u03b2"},
		[2]string{`\g`, "\u03b3"},
		[2]string{`\d`, "\u03b4"},
		[2]string{`\e`, "\u03b5"},
		[2]string{`\z`, "\u03b6"},
		[2]string{`\h`, "\u03b7"},
		[2]string{`\t`, "\u03b8"},
		[2]string{`\i`, "\u03b9"},
		[2]string{`\k`, "\u03ba"},
		[2]string{`\l`, "\u03bb"},
		[2]string{`\m`, "\u03bc"},
		[2]string{`\n`, "\u03bd"},
		[2]string{`\x`, "\u03be"},
		[2]string{`\o`, "\u03bf"},
		[2]string{`\p`, "\u03c0"},
		[2]string{`\r`, "\u03c1"},
		[2]string{`\s`, "\u03c3"},
		[2]string{`\y`, "\u03c4"},
		[2]string{`\u`, "\u03c5"},
		[2]string{`\f`, "\u03c6"},
		[2]string{`\c`, "\u03c7"},
		[2]string{`\q`, "\u03c8"},
		[2]string{`\w`, "\u03c9"},
		[2]string{`\0`, "\u2300"},
		[2]string{`v`, "\u2228"},
		[2]string{`-`, "\u00ac"},
		[2]string{`^`, "\u2227"},
		[2]string{`>`, "\u2283"},
		[2]string{`U`, "\u2200"},
		[2]string{`X`, "\u2203"},
		[2]string{"=", "="},
		[2]string{`\=`, "≠"},
		[2]string{`,`, ",\u00a0"},
	}
	charMap = make(map[string]string)
	revMap = make(map[string]string)
	for _, g := range bindings {
		charMap[g[0]] = g[1]
		revMap[g[1]] = g[0]
	}

	extrabindings := [][2]string{
		[2]string{",", ","},
		[2]string{" ", "\u00a0"},
		[2]string{"|U", "U"},
		[2]string{"|X", "X"},
		[2]string{"|v", "v"},
		[2]string{`|-`, "\u22a2"},
		[2]string{`|>`, ">"},
		[2]string{`||`, "|"},
		[2]string{`\\`, `\`},
	}
	extraMap = make(map[string]string)
	for _, g := range extrabindings {
		extraMap[g[0]] = g[1]
	}
}
