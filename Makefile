test:
	go test ./tests -json | go run ./test_parser.go

test-file:
	$(shell go test ./tests -json > output.json)
	go run ./test_parser.go -file 'output.json'
	@rm output.json

gen-fixtures:
	$(shell go test ./tests -json > ./files/fixture_fail.json)
	$(shell cat ./files/fixture_fail.json | go run test_parser.go > ./files/golden_fail.md)
	$(shell go test ./tests/tests_4_test.go -json > ./files/fixture_pass.json)
	$(shell cat ./files/fixture_pass.json | go run test_parser.go > ./files/golden_pass.md)
