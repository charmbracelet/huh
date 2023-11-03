$(V).SILENT:
test:
	go test ./...

taco:
	cd examples/taco && go run .

theme:
	cd examples/theme && go run .

gh:
	cd examples/gh && go run .
