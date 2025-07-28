module github.com/adamay909/logicTools

go 1.23.0

toolchain go1.24.5

replace github.com/adamay909/logicTools/gentzen => ./gentzen

require (
	golang.org/x/term v0.33.0
	honnef.co/go/js/dom/v2 v2.0.0-20210725211120-f030747120f2
)

require golang.org/x/sys v0.34.0 // indirect
