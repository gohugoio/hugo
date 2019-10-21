---
title: .AddDate
description: Returns the time corresponding to adding the given number of years, months, and days passed to the function.
godocref: https://golang.org/pkg/time/#Time.AddDate
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dates,time]
signature: [".AddDate YEARS MONTHS DAYS"]
workson: [times]
hugoversion:
relatedfuncs: [now]
deprecated: false
aliases: []
---


The `AddDate` function takes three arguments in logical order of `years`, `months`, and `days`.

## Example: Randomized Tweets from the Last 2 Years

Let's assume you have a file at `data/tweets.toml` that contains a list of Tweets to display on your site's homepage. The file is filled with `[[tweet]]` blocks; e.g.---

```
[[tweet]]
name = "Steve Francia"
twitter_handle = "@spf13"
quote = "I'm creator of Hugo. #metadocreference"
link = "https://twitter.com/spf13"
date = "2017-01-07T00:00:00Z"
```

Let's assume you want to grab Tweets from the last two years and present them in a random order. In conjunction with the [`where`](/functions/where/) and [`now`](/functions/now/) functions, you can limit our range to the last two years via `now.AddDate -2 0 0`, which represents a point in time 2 years, 0 days, and 0 hours before the time of your last site build.

{{< code file="partials/templates/random-tweets.html" download="tweets.html" >}}
{{ range where $.Site.Data.tweets.tweet "date" "ge" (now.AddDate -2 0 0) | shuffle }}
    <div class="item">
        <blockquote>
            <p>
            {{ .quote | safeHTML }}
            </p>
            &mdash; {{ .name }} ({{ .twitter_handle }}) <a href="{{ .link }}">
                {{ dateFormat "January 2, 2006" .date }}
            </a>
        </blockquote>
    </div>
{{ end }}
{{< /code >}}
