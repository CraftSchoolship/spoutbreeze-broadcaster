#!/bin/bash

set -euo pipefail

# Output files
COVERAGE_OUT="coverage.out"
COVERAGE_HTML="coverage.html"
JUNIT_REPORT="test-results.xml"
GINKGO_JSON_REPORT="report.json"
COVERAGE_HISTORY="coverage-history.csv"

echo "1. Run go tests with race and coverage"
go test -v -race -coverprofile=$COVERAGE_OUT -covermode=atomic ./...

echo "2. Generate HTML coverage report"
go tool cover -html=$COVERAGE_OUT -o $COVERAGE_HTML

echo "3. Get coverage percentage summary"
COVERAGE_PCT=$(go tool cover -func=$COVERAGE_OUT | grep total | awk '{print $3}')
echo "Total coverage: $COVERAGE_PCT"

echo "4. Run Ginkgo with detailed reporting and coverage"
ginkgo -r --randomize-all --randomize-suites --fail-on-pending --cover --coverprofile=$COVERAGE_OUT --race --trace -v

echo "5. Generate JUnit XML for CI integration"
ginkgo -r --randomize-all --randomize-suites --fail-on-pending --cover --junit-report=$JUNIT_REPORT

echo "6. Measure test execution time for Ginkgo"
time ginkgo -r --randomize-all --randomize-suites --fail-on-pending

echo "7. Count total test cases"
TOTAL_TESTS=$(grep -r "It(" . --include="*_test.go" | wc -l)
echo "Total test cases: $TOTAL_TESTS"

echo "8. Count integration tests (with 'Integration' in Describe)"
INTEGRATION_TESTS=$(grep -r "Integration" . --include="*_test.go" | wc -l)
echo "Integration tests: $INTEGRATION_TESTS"

echo "9. Generate coverage badge data"
echo "Coverage badge data:"
go tool cover -func=$COVERAGE_OUT | grep total | awk '{print "Coverage: " $3}'

echo "10. Package-wise coverage breakdown"
go tool cover -func=$COVERAGE_OUT | grep -v total | awk '{print $1 " " $3}' | sort -k2 -nr

# Optional: 11 requires gocov and gocov-html installed
if command -v gocov >/dev/null 2>&1 && command -v gocov-html >/dev/null 2>&1; then
  echo "11. Generate detailed coverage report with gocov"
  gocov test ./... | gocov-html > coverage-detailed.html
else
  echo "11. Skipping detailed coverage report (gocov or gocov-html not installed)"
fi

echo "12. Run benchmarks for performance"
go test -bench=. -benchmem ./...

echo "13. Profile test execution"
go test -cpuprofile cpu.prof -memprofile mem.prof -bench . ./...

echo "14. Generate Ginkgo JSON report"
ginkgo -r --randomize-all --randomize-suites --json-report=$GINKGO_JSON_REPORT

if command -v jq >/dev/null 2>&1; then
  echo "15. Extract specific metrics from JSON report"
  cat $GINKGO_JSON_REPORT | jq '.SuiteSuccesses, .SuiteFailures, .RunTime'
else
  echo "15. Skipping JSON metrics extraction (jq not installed)"
fi

echo "16. File-by-file coverage details"
go tool cover -func=$COVERAGE_OUT

echo "17. Append coverage trend data to $COVERAGE_HISTORY"
echo "$(date '+%Y-%m-%d %H:%M:%S'),$(go tool cover -func=$COVERAGE_OUT | grep total | awk '{print $3}')" >> $COVERAGE_HISTORY

echo "18. Coverage comparison between packages"
go tool cover -func=$COVERAGE_OUT | awk 'NF==3 {pkg=substr($1,1,index($1,":")-1); coverage[pkg]+=$3; count[pkg]++} END {for (p in coverage) printf "%s: %.1f%%\n", p, coverage[p]/count[p]}' | sort -t: -k2 -nr

echo "=== Automation Complete ==="
