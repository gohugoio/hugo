---
title: Math
description: Hugo provides nine mathematical operators in templates.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [math, operators]
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
signature: []
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

| Function       | Description                                                                   | Example                          |
|----------------|-------------------------------------------------------------------------------|----------------------------------|
| `add`          | Adds two integers.                                                            | `{{add 1 2}}` &rarr; 3           |
| `div`          | Divides two integers.                                                         | `{{div 6 3}}` &rarr; 2           |
| `mod`          | Modulus of two integers.                                                      | `{{mod 15 3}}` &rarr; 0          |
| `modBool`      | Boolean of modulus of two integers. Evaluates to `true` if result equals 0.   | `{{modBool 15 3}}` &rarr; true   |
| `mul`          | Multiplies two integers.                                                      | `{{mul 2 3}}` &rarr; 6           |
| `sub`          | Subtracts two integers.                                                       | `{{sub 3 2}}` &rarr; 1           |
| `math.Ceil`    | Returns the least integer value greater than or equal to the given number.    | `{{math.Ceil 2.1}}` &rarr; 3     |
| `math.Floor`   | Returns the greatest integer value less than or equal to the given number.    | `{{math.Floor 1.9}}` &rarr; 1    |
| `math.Round`   | Returns the nearest integer, rounding half away from zero.                    | `{{math.Round 1.5}}` &rarr; 2    |
