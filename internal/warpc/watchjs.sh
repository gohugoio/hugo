#!/bin/bash

trap exit SIGINT

while true; do find . -type f -name "*.js" | entr -pd ./buildjs.sh; done
