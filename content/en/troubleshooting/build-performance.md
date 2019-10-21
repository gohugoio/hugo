---
title: Build Performance
linktitle: Build Performance
description: An overview of features used for diagnosing and improving performance issues in site builds.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-12
keywords: [performance, build]
categories: [troubleshooting]
menu:
  docs:
    parent: "troubleshooting"
weight: 3
slug:
aliases: []
toc: true
---

{{% note %}}
The example site used below is from https://github.com/gohugoio/hugo/tree/master/examples/blog
{{% /note %}}

## Template Metrics

Hugo is a very fast static site generator, but it is possible to write
inefficient templates.  Hugo's *template metrics* feature is extremely helpful
in pinpointing which templates are executed most often and how long those
executions take **in terms of CPU time**.

| Metric Name         | Description |
|---------------------|-------------|
| cumulative duration | The cumulative time spent executing a given template. |
| average duration    | The average time spent executing a given template. |
| maximum duration    | The maximum time a single execution took for a given template. |
| count               | The number of times a template was executed. |
| template            | The template name. |

```
▶ hugo --templateMetrics
Started building sites ...

Built site for language en:
0 draft content
0 future content
0 expired content
2 regular pages created
22 other pages created
0 non-page files copied
0 paginator pages created
4 tags created
3 categories created
total in 18 ms

Template Metrics:

     cumulative       average       maximum
       duration      duration      duration  count  template
     ----------      --------      --------  -----  --------
     6.419663ms     583.605µs     994.374µs     11  _internal/_default/rss.xml
     4.718511ms    1.572837ms    3.880742ms      3  indexes/category.html
     4.642666ms    2.321333ms    3.282842ms      2  posts/single.html
     4.364445ms     396.767µs    2.451372ms     11  partials/header.html
     2.346069ms     586.517µs     903.343µs      4  indexes/tag.html
     2.330919ms     211.901µs    2.281342ms     11  partials/header.includes.html
     1.238976ms     103.248µs     446.084µs     12  posts/li.html
       972.16µs      972.16µs      972.16µs      1  _internal/_default/sitemap.xml
      953.597µs     953.597µs     953.597µs      1  index.html
      822.263µs     822.263µs     822.263µs      1  indexes/post.html
      567.498µs       51.59µs     112.205µs     11  partials/navbar.html
       348.22µs      31.656µs      88.249µs     11  partials/meta.html
      346.782µs     173.391µs     276.176µs      2  posts/summary.html
      235.184µs       21.38µs     124.383µs     11  partials/footer.copyright.html
      132.003µs          12µs     117.999µs     11  partials/menu.html
       72.547µs       6.595µs      63.764µs     11  partials/footer.html
```

{{% note %}}
**A Note About Parallelism**

Hugo builds pages in parallel where multiple pages are generated
simultaneously. Because of this parallelism, the sum of "cumulative duration"
values is usually greater than the actual time it takes to build a site.
{{% /note %}}


## Cached Partials

Some `partial` templates such as sidebars or menus are executed many times
during a site build.  Depending on the content within the `partial` template and
the desired output, the template may benefit from caching to reduce the number
of executions.  The [`partialCached`][partialCached] template function provides
caching capabilities for `partial` templates.

{{% tip %}}
Note that you can create cached variants of each `partial` by passing additional
parameters to `partialCached` beyond the initial context.  See the
`partialCached` documentation for more details.
{{% /tip %}}


[partialCached]:{{< ref "/functions/partialCached.md" >}}
