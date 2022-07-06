build:
	cd assets/markdown&& \
	for i in *.md; do md2html $$i > ../html/$${i%.md}.html; done
	GOOS=js GOARCH=wasm go build -o ../docs/assets/proofchecker.wasm

