#!/bin/bash

# Send in a regexp mathing the benchmarks you want to run, i.e. './benchSite.sh "YAML"'. 
# Note the quotes, which will be needed for more complex expressions.
# The above will run all variations, but only for front matter YAML.

echo "Running with BenchmarkSiteBuilding/${1}"

go test -run="NONE" -bench="BenchmarkSiteBuilding/${1}$" -test.benchmem=true ./hugolib -memprofile mem.prof -cpuprofile cpu.prof