---
title: math
linktitle: Math
description: Hugo provides six mathematical operators in templates.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [math, operators]
categories: [functions]
toc:
ns:
signature: []
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

There are 6 basic mathematical operators that can be used in Hugo templates:

| Function | Description              | Example                       |
| -------- | ------------------------ | ----------------------------- |
| `add`    | Adds two integers.       | `{{add 1 2}}` &rarr; 3        |
| `div`    | Divides two integers.    | `{{div 6 3}}` &rarr; 2        |
| `mod`    | Modulus of two integers. | `{{mod 15 3}}` &rarr; 0       |
| `modBool`| Boolean of modulus of two integers. Evaluates to `true` if = 0. | `{{modBool 15 3}}` &rarr; true |
| `mul`    | Multiplies two integers. | `{{mul 2 3}}` &rarr; 6        |
| `sub`    | Subtracts two integers.  | `{{sub 3 2}}` &rarr; 1        |
