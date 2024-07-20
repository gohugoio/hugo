---
title: time.Now
description: Returns the current local time.
categories: []
keywords: []
action:
  aliases: [now]
  related:
    - functions/time/AsTime
    - functions/time/Duration
    - functions/time/Format
    - functions/time/ParseDuration
  returnType: time.Time
  signatures: [time.Now]
aliases: [/functions/now]
---

For example, when building a site on October 15, 2023 in the America/Los_Angeles time zone:

```go-html-template
{{ time.Now }}
```

This produces a `time.Time` value, with a string representation such as:

```text
2023-10-15 12:59:28.337140706 -0700 PDT m=+0.041752605
```

To format and [localize] the value, pass it through the [`time.Format`] function:

```go-html-template
{{ time.Now | time.Format "Jan 2006" }} → Oct 2023
```

The `time.Now` function returns a `time.Time` value, so you can chain any of the [time methods] to the resulting value. For example:


```go-html-template
{{ time.Now.Year }} → 2023 (int)
{{ time.Now.Weekday.String }} → Sunday
{{ time.Now.Month.String }} → October
{{ time.Now.Unix }} → 1697400955 (int64)
```

[`time.Format`]: /functions/time/format/
[localize]: /getting-started/glossary/#localization
[time methods]: /methods/time/
