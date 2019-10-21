---
title: strings.TrimSuffix
description: Returns a given string s without the provided trailing suffix string. If s doesn't end with suffix, s is returned unchanged.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["strings.TrimSuffix SUFFIX STRING"]
workson: []
hugoversion:
relatedfuncs: [strings.TrimPrefix]
deprecated: false
aliases: []
---

Given the string `"aabbaa"`, the specified suffix is only removed if `"aabbaa"` ends with it:

    {{ strings.TrimSuffix "a" "aabbaa" }} → "aabba"
    {{ strings.TrimSuffix "aa" "aabbaa" }} → "aabb"
    {{ strings.TrimSuffix "aaa" "aabbaa" }} → "aabbaa"