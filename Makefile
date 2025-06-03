test: test_e2e
	go test ./... --race -cover

test_e2e:
	go test ./__test__

COVERAGE_DIR="coverage"
cover:
	@go test -coverprofile=./$(COVERAGE_DIR)/coverage.out ./... && mkdir -p $(COVERAGE_DIR) && go tool cover -html=./$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/index.html

new-version:
	@go run ./version