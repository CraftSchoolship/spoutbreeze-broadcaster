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
	@echo "📊 Coverage Summary:"
	@go tool cover -func=coverage.out
	@rm -f coverage.out

test-coverage-core:
	@echo "🎯 Running core package tests with 100% coverage..."
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --cover --coverprofile=core-coverage.out --trace ./repositories ./models ./routes ./controllers/tests
	@echo "📊 Core Package Coverage Summary:"
	@go tool cover -func=core-coverage.out
	@rm -f core-coverage.out

test-watch:
	@ginkgo watch --randomize-all --fail-on-pending --trace -v ./...

test-ci:
	@go mod tidy
	@ginkgo run --randomize-all --randomize-suites --fail-on-pending --cover --race --trace --junit-report=test-results.xml ./...

clean:
	@echo "🧹 Cleaning up generated files..."
	@rm -f coverage.out core-coverage.out coverage.html core-coverage.html test-results.xml
	@echo "✅ Cleanup complete!"
