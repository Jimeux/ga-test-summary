test:
	go test ./tests -json | go run ./test_parser.go

gen-fixtures:
	go test ./tests -json > ./files/fixture_fail.json || true
	cat ./files/fixture_fail.json | go run test_parser.go > ./files/golden_fail.md
	go test ./tests/tests_4_test.go -json > ./files/fixture_pass.json
	cat ./files/fixture_pass.json | go run test_parser.go > ./files/golden_pass.md
