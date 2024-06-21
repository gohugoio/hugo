---
title: Date
description: Returns the date of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/ExpiryDate
    - methods/page/LastMod
    - methods/page/PublishDate
  returnType: time.Time
  signatures: [PAGE.Date]
---

Set the date in front matter:

{{< code-toggle file=content/news/article-1.md fm=true >}}
title = 'Article 1'
date = 2023-10-19T00:40:04-07:00
{{< /code-toggle >}}

{{% note %}}
The date field in front matter is often considered to be the creation date, You can change its meaning, and its effect on your site, in the site configuration. See&nbsp;[details].

[details]: /getting-started/configuration/#configure-dates
{{% /note %}}

The date is a [time.Time] value. Format and localize the value with the [`time.Format`] function, or use it with any of the [time methods].

```go-html-template
{{ .Date | time.Format ":date_medium" }} → Oct 19, 2023
```

In the example above we explicitly set the date in front matter. With Hugo's default configuration, the `Date` method returns the front matter value. This behavior is configurable, allowing you to set fallback values if the date is not defined in front matter. See&nbsp;[details].

[`time.Format`]: /functions/time/format/
[details]: /getting-started/configuration/#configure-dates
[time methods]: /methods/time/
[time.Time]: https://pkg.go.dev/time#Time
