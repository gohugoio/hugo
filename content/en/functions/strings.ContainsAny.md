---
title: strings.ContainsAny
description: Reports whether a string contains any character from a given string.
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [string strings substring contains any]
signature: ["strings.ContainsAny STRING CHARACTERS"]
aliases: []
relatedfuncs: [strings.Contains]
---

    {{ strings.ContainsAny "Hugo" "gm" }} → true

The check is case sensitive: 

    {{ strings.ContainsAny "Hugo" "Gm" }} → false
