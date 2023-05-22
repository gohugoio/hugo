---
title: strings.TrimSuffix
description: Returns a given string s without the provided trailing suffix string. If s doesn't end with suffix, s is returned unchanged.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["strings.TrimSuffix SUFFIX STRING"]
relatedfuncs: [strings.TrimPrefix]
---

Given the string `"aabbaa"`, the specified suffix is only removed if `"aabbaa"` ends with it:

    {{ strings.TrimSuffix "a" "aabbaa" }} → "aabba"
    {{ strings.TrimSuffix "aa" "aabbaa" }} → "aabb"
    {{ strings.TrimSuffix "aaa" "aabbaa" }} → "aabbaa"
