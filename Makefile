test:
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --trace ./...

test-verbose:
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --trace -v ./...

test-controllers:
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --trace ./controllers/tests

test-coverage:
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --cover --coverprofile=coverage.out --trace ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated at: coverage.html"

test-watch:
	@ginkgo watch --randomize-all --fail-on-pending --trace -v ./...

test-ci:
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --cover --race --trace --junit-report=test-results.xml ./...
