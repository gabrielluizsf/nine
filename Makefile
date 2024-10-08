test:
	go test --race -cover

COVERAGE_DIR="coverage"
cover:
	@go test -coverprofile=./$(COVERAGE_DIR)/coverage.out ./... && mkdir -p $(COVERAGE_DIR) && go tool cover -html=./$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/index.html