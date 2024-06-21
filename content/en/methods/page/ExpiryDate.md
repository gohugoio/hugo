---
title: ExpiryDate
description: Returns the expiry date of the given page. 
categories: []
keywords: []
action:
  related:
    - methods/page/Date
    - methods/page/LastMod
    - methods/page/PublishDate
  returnType: time.Time
  signatures: [PAGE.ExpiryDate]
---

By default, Hugo excludes expired pages when building your site. To include expired pages, use the `--buildExpired` command line flag.

Set the expiry date in front matter:

{{< code-toggle file=content/news/article-1.md fm=true >}}
title = 'Article 1'
expiryDate = 2024-10-19T00:32:13-07:00
{{< /code-toggle >}}

The expiry date is a [time.Time] value. Format and localize the value with the [`time.Format`] function, or use it with any of the [time methods].

```go-html-template
{{ .ExpiryDate | time.Format ":date_medium" }} → Oct 19, 2024
```

In the example above we explicitly set the expiry date in front matter. With Hugo's default configuration, the `ExpiryDate` method returns the front matter value. This behavior is configurable, allowing you to set fallback values if the expiry date is not defined in front matter. See&nbsp;[details].

[`time.Format`]: /functions/time/format/
[details]: /getting-started/configuration/#configure-dates
[time methods]: /methods/time/
[time.Time]: https://pkg.go.dev/time#Time
