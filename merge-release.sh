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

# If $v contains -alpha, -beta or -rc, skip the remaining steps.
if [[ "${v}" == *"-alpha"* ||  "${v}" == *"-beta"* ||  "${v}" == *"-rc"* ]]; then
  echo "Pre-release version detected; skipping stable and docs update."
  exit 0
fi

git checkout stable || die;
git reset --hard "v${v}" || die;
git push -f || die;

git checkout master || die;

 git subtree push --prefix=docs/ docs-local "tempv${v}";


