---
title: substr
# linktitle:
description: Extracts parts of a string from a specified character's position and returns the specified number of characters.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
aliases: []
signature: ["substr STRING START [LENGTH]"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

It normally takes two parameters: `start` and `length`. It can also take one parameter: `start`, i.e. `length` is omitted, in which case the substring starting from start until the end of the string will be returned.

To extract characters from the end of the string, use a negative start number.

In addition, borrowing from the extended behavior described at https://php.net substr, if `length` is given and is negative, that number of characters will be omitted from the end of string.

```
{{substr "BatMan" 0 -3}} → "Bat"
{{substr "BatMan" 3 3}} → "Man"
```
