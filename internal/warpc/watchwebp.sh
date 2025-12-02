#!/bin/bash

trap exit SIGINT


while true; do find genwebp -type f -name "*.c" | entr -pd ./buildwebp.sh; done

