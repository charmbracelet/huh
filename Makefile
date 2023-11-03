.PHONY: spinner

$(V).SILENT:
test:
	go test ./...

spinner:
	cd spinner/examples/loading && go run .

taco:
	cd examples/taco && go run .

theme:
	cd examples/theme && go run .

gh:
	cd examples/gh && go run .
