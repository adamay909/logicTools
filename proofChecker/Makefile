build:
	cp README.md assets/markdown/README.md &&\
	cd assets/markdown &&\
	for i in *.md; do md2html -o ../html/$${i%.md}.html $$i; done &&\
	cd ../.. &&\
	cd assets/html &&\
	i=`cat version` &&\
	let i++ &&\
	echo $$i>version &&\
	cd ../.. &&\
	cd assets &&\
	./genHTML.lua &&\
	cd .. &&\
	GOOS=js GOARCH=wasm go build  -o ../docs/assets/proofchecker.wasm &&\
	mv assets/html/index.html ../docs/ &&\
	cp /usr/lib/go/lib/wasm/wasm_exec.js ../docs/assets/


