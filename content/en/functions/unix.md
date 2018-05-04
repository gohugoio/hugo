---
title: .Unix
draft: false
description: .Unix returns the local Time corresponding to the given Unix time, sec seconds and nsec nanoseconds since January 1, 1970 UTC.
godocref: https://golang.org/search?q=Unix#Functions
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [dates,time]
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: [".Unix"]
workson: [times]
hugoversion:
relatedfuncs: [Format,dateFormat,now,time]
deprecated: false
aliases: []
---

## Example: Time Passed Since Last Modification

This very simple one-liner uses `now.Unix` to calculate the amount of time that has passed between the `.LastMod` for the current page and the last build of the current page.

{{< code file="time-passed.html" >}}
{{ div (sub now.Unix .Lastmod.Unix) 86400 }}
{{< /code >}}

Since both values are integers, they can be subtracted and then divided by the number of seconds in a day (i.e., `60 * 60 * 24 == 86400`).

{{% note %}}
Hugo's output is *static*. For the example above to be realistic, the site needs to be built every day.
 {{% /note %}}



[partial template]: /templates/partials/
