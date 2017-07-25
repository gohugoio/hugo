#!/usr/bin/env bash

# allow user to override go executable by running as GOEXE=xxx make ...
GOEXE="${GOEXE-go}"

# Convenience script to
# - For a given branch
# - Run  benchmark tests for a given package
# - Do the same for master
# - then compare the two runs with benchcmp

benchFilter=".*"

if (( $# < 2 ));
  then
    echo "USAGE: ./bench.sh <git-branch> <package-to-bench> (and <benchmark filter> (regexp, optional))"
    exit 1
fi



if [ $# -eq 3 ]; then
  benchFilter=$3
fi


BRANCH=$1
PACKAGE=$2

git checkout $BRANCH
"${GOEXE}" test -test.run=NONE -bench="$benchFilter" -test.benchmem=true ./$PACKAGE > /tmp/bench-$PACKAGE-$BRANCH.txt

git checkout master
"${GOEXE}" test -test.run=NONE -bench="$benchFilter" -test.benchmem=true ./$PACKAGE > /tmp/bench-$PACKAGE-master.txt


benchcmp /tmp/bench-$PACKAGE-master.txt /tmp/bench-$PACKAGE-$BRANCH.txt
