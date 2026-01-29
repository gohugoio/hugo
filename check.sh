#!/bin/bash

set -e

# Default to all packages if none specified
PACKAGES="${1:-./...}"

echo "==> Checking packages: $PACKAGES"

# Timing arrays
declare -a STEP_NAMES
declare -a STEP_TIMES

time_step() {
    local name="$1"
    shift
    local start=$(date +%s.%N)
    "$@"
    local end=$(date +%s.%N)
    local elapsed=$(echo "$end - $start" | bc)
    STEP_NAMES+=("$name")
    STEP_TIMES+=("$elapsed")
}

# Check gofmt
run_gofmt() {
    echo "==> Running gofmt..."
    # Convert package pattern to path (e.g., ./hugolib/... -> ./hugolib)
    local path="${PACKAGES%/...}"
    GOFMT_OUTPUT=$(gofmt -l "$path" 2>&1) || true
    if [ -n "$GOFMT_OUTPUT" ]; then
        echo "gofmt found issues in:"
        echo "$GOFMT_OUTPUT"
        exit 1
    fi
    echo "    OK"
}

# Run staticcheck
run_staticcheck() {
    # Check if staticcheck is installed, install if not
    if ! command -v staticcheck &> /dev/null; then
        echo "==> Installing staticcheck..."
        go install honnef.co/go/tools/cmd/staticcheck@latest
    fi
    echo "==> Running staticcheck..."
    staticcheck $PACKAGES
    echo "    OK"
}

# Run tests
run_tests() {
    echo "==> Running tests..."
    local output
    if ! output=$(go test $PACKAGES 2>&1); then
        echo "$output"
        exit 1
    fi
    echo "    OK"
}

# Run all steps with timing
TOTAL_START=$(date +%s.%N)

time_step "gofmt" run_gofmt
time_step "staticcheck" run_staticcheck
time_step "tests" run_tests

TOTAL_END=$(date +%s.%N)
TOTAL_ELAPSED=$(echo "$TOTAL_END - $TOTAL_START" | bc)

# Print timing summary
echo ""
echo "==> All checks passed!"
echo ""
echo "Timing summary:"
echo "---------------"
for i in "${!STEP_NAMES[@]}"; do
    printf "  %-15s %6.2fs\n" "${STEP_NAMES[$i]}" "${STEP_TIMES[$i]}"
done
echo "---------------"
printf "  %-15s %6.2fs\n" "Total" "$TOTAL_ELAPSED"
