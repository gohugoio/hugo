#!/bin/bash

# Se https://github.com/gohugoio/hugo/issues/8955
objdump -T dist/hugo_extended_linux_linux_amd64/hugo | grep -E -q 'GLIBC_2.2[0-9]'
RESULT=$?
if [ $RESULT -eq 0 ]; then
    echo "Found  GLIBC_2.2x in Linux binary, this will not work in older Vercel/Netlify images.";
    exit -1;
fi
