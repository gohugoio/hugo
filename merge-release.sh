#!/usr/bin/env bash

if (( $# < 1 ));
  then
    echo "USAGE: ./merge-release.sh 0.76.0"
    exit 1
fi

die() { echo "$*" 1>&2 ; exit 1; }

v=$1
git merge "release-${v}" || die;
git push || die;

git checkout stable || die;
git reset --hard "v${v}" || die;
git push -f || die;

git checkout master || die;

 git subtree push --prefix=docs/ docs-local "tempv${v}";


