#!/bin/bash

# allow user to override go executable by running as GOEXE=xxx make ...
GOEXE="${GOEXE-go}"

# Send in a regexp matching the benchmarks you want to run, i.e. './benchSite.sh "YAML"'.
# Note the quotes, which will be needed for more complex expressions.
# The above will run all variations, but only for front matter YAML.

echo "Running with BenchmarkSiteBuilding/${1}"

"${GOEXE}" test -run="NONE" -bench="BenchmarkSiteBuilding/${1}" -test.benchmem=true ./hugolib -memprofile mem.prof -count 3 -cpuprofile cpu.prof
