---
title: strings.TrimPrefix
description: Returns a given string s without the provided leading prefix string. If s doesn't start with prefix, s is returned unchanged.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["strings.TrimPrefix PREFIX STRING"]
relatedfuncs: [strings.TrimSuffix]
---

Given the string `"aabbaa"`, the specified prefix is only removed if `"aabbaa"` starts with it:

    {{ strings.TrimPrefix "a" "aabbaa" }} → "abbaa"
    {{ strings.TrimPrefix "aa" "aabbaa" }} → "bbaa"
    {{ strings.TrimPrefix "aaa" "aabbaa" }} → "aabbaa"
