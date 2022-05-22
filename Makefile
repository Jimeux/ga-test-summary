test:
	go test ./... -json | go run ./.github/test_parser.go

test-file:
	$(shell go test ./... -json > output.json)
	go run ./.github/test_parser.go -file 'output.json'
	@rm output.json

gen-fixture:
	$(shell go test -v ./... -json > .github/fixture.json)
	$(shell cat .github/fixture.json | go run .github/test_parser.go > .github/golden.md)
