module github.com/adamay909/logicTools

go 1.18

replace github.com/adamay909/logicTools/gentzen => ./gentzen

<<<<<<< HEAD
require (
	github.com/adamay909/logicTools/gentzen v0.0.0-00010101000000-000000000000
	honnef.co/go/js/dom/v2 v2.0.0-20210725211120-f030747120f2
)
=======
require github.com/adamay909/logicTools/gentzen v0.0.0-00010101000000-000000000000 // indirect

replace github.com/adamay909/logicTools/gentzen => ./gentzen
>>>>>>> devel
