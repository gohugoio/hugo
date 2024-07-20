#!/bin/bash

trap exit SIGINT

# I use "run tests on save" in my editor.
# Unfortunately, changes to text files does not trigger this. Hence this workaround.
while true; do find testscripts -type f -name "*.txt" | entr -pd touch main_test.go; done