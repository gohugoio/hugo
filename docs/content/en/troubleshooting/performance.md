---
title: Performance
description: Tools and suggestions for evaluating and improving performance.
categories: []
keywords: []
aliases: [/troubleshooting/build-performance/]
---

## Virus scanning

Virus scanners are an essential component of system protection, but the performance impact can be severe for applications like Hugo that frequently read and write to disk. For example, with Microsoft Defender Antivirus, build times for some sites may increase by 400% or more.

Before building a site, your virus scanner has already evaluated the files in your project directory. Scanning them again while building the site is superfluous. To improve performance, add Hugo's executable to your virus scanner's process exclusion list.

For example, with Microsoft Defender Antivirus:

**Start**&nbsp;> **Settings**&nbsp;> **Privacy&nbsp;&&nbsp;security**&nbsp;> **Windows&nbsp;Security**&nbsp;> **Open&nbsp;Windows&nbsp;Security**&nbsp;> **Virus&nbsp;&&nbsp;threat&nbsp;protection**&nbsp;> **Manage&nbsp;settings**&nbsp;> **Add&nbsp;or&nbsp;remove&nbsp;exclusions**&nbsp;> **Add&nbsp;an&nbsp;exclusion**&nbsp;> **Process**

Then type `hugo.exe` and press the **Add** button.

> [!NOTE]
> Virus scanning exclusions are common, but use caution when changing these settings. See the [Microsoft Defender Antivirus documentation][] for details.

Other virus scanners have similar exclusion mechanisms. See their respective documentation.

## Template metrics

Hugo is fast, but inefficient templates impede performance. Enable template metrics to determine which templates take the most time, and to identify caching opportunities:

```sh
hugo build --templateMetrics --templateMetricsHints
```

The result will look something like this:

<!--
Run hugo build --templateMetrics --templateMetricsHints in this repo and paste
the results below. You must remove any line containing "__hdeferred" to prevent
this:

panic: deferred execution with id "__h" not found
-->

```text {details=true}
Template Metrics:

       cumulative       average       maximum      cache  percent  cached  total  
         duration      duration      duration  potential   cached   count  count  template
       ----------      --------      --------  ---------  -------  ------  -----  --------
         14.95  s      20.54 ms     278.74 ms          0        0       0    728  single.html
          3.20  s       4.03 ms      33.23 ms         99        0       0    793  _partials/layouts/header/header.html
          2.91  s       3.67 ms      29.58 ms         34        0       0    793  _partials/layouts/head/head.html
          1.82  s       2.30 ms      23.97 ms        100        0       0    793  _partials/layouts/head/head-js.html
          1.54  s     231.79 µs      23.49 ms          0        0       0   6632  _markup/render-link.html
          1.37  s       1.73 ms      19.99 ms         91        0       0    793  _partials/layouts/header/qr.html
        898.03 ms     338.24 µs      28.84 ms          0        0       0   2655  _markup/render-codeblock.html
        828.14 ms      13.15 ms      42.49 ms          0        0       0     63  list.html
        648.23 ms     817.44 µs      18.71 ms        100        0       0    793  _partials/layouts/footer.html
        624.87 ms     787.98 µs      18.88 ms         23        0       0    793  _partials/opengraph/opengraph.html
        515.61 ms       1.44 ms      23.17 ms          0        0       0    359  _shortcodes/code-toggle.html
        431.74 ms     544.44 µs      16.66 ms        100        0       0    793  _partials/layouts/search/input.html
        361.07 ms     227.66 µs       3.56 ms         58        0       0   1586  _partials/_inline/qr
        353.60 ms     445.90 µs      17.17 ms         18        0       0    793  _partials/schema.html
        292.99 ms     369.00 µs      18.63 ms         99        0       0    794  _partials/layouts/home/sponsors.html
        282.03 ms     355.64 µs      15.65 ms         28        0       0    793  _partials/twitter_cards.html
        279.74 ms     279.74 ms     279.74 ms        100        0       0      1  _partials/helpers/linkcss.html
        230.48 ms     642.01 µs       9.04 ms          0        0       0    359  _markup/render-blockquote.html
        227.22 ms     287.25 µs      12.19 ms         38        0       0    791  _partials/layouts/docsheader.html
        206.41 ms     260.29 µs      10.04 ms        100        0       0    793  _partials/layouts/header/mobilemenu.html
        202.75 ms     127.83 µs      15.57 ms         75        0       0   1586  _partials/helpers/linkjs.html
        152.33 ms     192.09 µs      16.63 ms        100        0       0    793  _partials/layouts/search/results.html
        137.02 ms     172.78 µs      10.88 ms        100        0       0    793  _partials/layouts/header/githubstars.html
        121.06 ms     152.66 µs       1.55 ms          0        0       0    793  _partials/opengraph/get-featured-image.html
        120.00 ms     165.07 µs       2.71 ms         15        0       0    727  _partials/layouts/in-this-section.html
        116.10 ms     607.86 µs      11.93 ms          0        0       0    191  _shortcodes/new-in.html
         84.77 ms      29.47 µs       5.92 ms        100        0       0   2876  _partials/inline/h-rh-l/validate-fragment.html
         71.97 ms      80.51 µs       2.81 ms          0        0       0    894  _partials/inline/h-rh-l/get-glossary-link-attributes.html
         70.27 ms      88.84 µs       5.66 ms         42        0       0    791  _partials/layouts/breadcrumbs.html
         53.72 ms       1.68 ms       8.50 ms          0        0       0     32  _shortcodes/img.html
         50.23 ms      63.34 µs       4.77 ms          0        0       0    793  _partials/layouts/hooks/body-main-start.html
         48.01 ms      48.01 ms      48.01 ms          0        0       0      1  home.html
         42.23 ms      42.23 ms      42.23 ms          0        0       0      1  404.html
         41.13 ms      56.50 µs       2.49 ms          5        0       0    728  _partials/layouts/toc.html
         39.26 ms      24.75 µs       1.28 ms         99        0       0   1586  _partials/_funcs/get-page-images.html
         24.70 ms      31.15 µs       1.09 ms        100        0       0    793  _partials/helpers/gtag.html
         23.87 ms       2.39 ms      10.01 ms          0        0       0     10  _markup/render-codeblock-goat.html
         22.62 ms      96.24 µs       2.12 ms         80        0       0    235  _partials/layouts/blocks/feature-state.html
         21.58 ms       9.07 µs       3.63 ms        100      100    2378   2380  _partials/helpers/funcs/get-github-info.html
         21.19 ms      29.10 µs     902.28 µs         97        0       0    728  _partials/layouts/page-edit.html
         20.03 ms      90.24 µs       1.49 ms          0        0       0    222  _shortcodes/include.html
         20.02 ms      27.50 µs     593.31 µs         97        0       0    728  _partials/layouts/related.html
         19.91 ms      19.91 ms      19.91 ms          0        0       0      1  /news/_content.gotmpl
         19.83 ms     450.58 µs       4.69 ms          0        0       0     44  _shortcodes/deprecated-in.html
         18.70 ms      18.70 ms      18.70 ms        100        0       0      1  _partials/helpers/funcs/get-remote-data.html
         17.03 ms      74.37 µs       1.80 ms          0        0       0    229  _markup/render-table.html
         10.61 ms      13.38 µs       2.29 ms         90        0       0    793  _partials/layouts/blocks/modal.html
         10.14 ms      12.79 µs       1.17 ms        100        0       0    793  _partials/layouts/header/mastodon.html
          9.57 ms       9.57 ms       9.57 ms          0        0       0      1  sitemap.xml
          9.35 ms      24.68 µs     637.73 µs         50        0       0    379  _partials/layouts/blocks/alert.html
          6.96 ms       8.77 µs     935.10 µs         99        0       0    794  _partials/layouts/search/button.html
          5.38 ms       2.69 ms       3.20 ms          0        0       0      2  _shortcodes/quick-reference.html
          5.32 ms       9.92 µs       1.03 ms         75        0       0    536  _partials/docs/functions-signatures.html
          5.13 ms       6.47 µs       1.36 ms        100        0       0    793  _partials/layouts/hooks/body-end.html
          5.02 ms     152.09 µs       1.22 ms          0        0       0     33  _shortcodes/glossary-term.html
          4.65 ms     132.78 µs       1.71 ms          2        0       0     35  _partials/inline/get-resource.html
          3.98 ms       1.99 ms       3.56 ms          0        0       0      2  _shortcodes/datatable.html
          3.89 ms     388.92 µs     929.41 µs          0        0       0     10  _shortcodes/render-list-of-pages-in-section.html
          3.80 ms       2.39 µs     963.41 µs        100        0       0   1586  _partials/layouts/header/theme.html
          3.44 ms       6.42 µs     795.02 µs         19        0       0    536  _partials/docs/functions-aliases.html
          3.10 ms       3.10 ms       3.10 ms          0        0       0      1  _shortcodes/glossary.html
          2.57 ms       4.80 µs     256.63 µs         88        0       0    536  _partials/docs/functions-return-type.html
          2.54 ms       2.54 ms       2.54 ms        100        0       0      1  _partials/layouts/home/opensource.html
          1.97 ms       2.48 µs     501.94 µs        100        0       0    793  _partials/layouts/icons.html
          1.85 ms      32.98 µs     178.18 µs          0        0       0     56  _shortcodes/eturl.html
          1.74 ms       1.74 ms       1.74 ms          0        0       0      1  home.redir
          1.70 ms     340.42 µs     596.87 µs          0        0       0      5  _shortcodes/qr.html
          1.65 ms     235.82 µs       1.03 ms          0        0       0      7  _markup/render-passthrough.html
          1.53 ms       1.53 ms       1.53 ms          0        0       0      1  _shortcodes/syntax-highlighting-styles.html
          1.51 ms       1.51 ms       1.51 ms          0        0       0      1  _shortcodes/root-configuration-keys.html
          1.45 ms       1.82 µs     839.08 µs        100        0       0    793  _partials/layouts/templates.html
          1.34 ms       1.34 ms       1.34 ms        100        0       0      1  _partials/helpers/validation/validate-keywords.html
          1.30 ms       1.30 ms       1.30 ms          0        0       0      1  _shortcodes/per-lang-config-keys.html
          1.18 ms      21.89 µs      76.24 µs          0        0       0     54  _markup/render-image.html
          1.15 ms       1.46 µs      36.50 µs        100        0       0    793  _partials/layouts/search/algolialogo.html
        775.57 µs     775.57 µs     775.57 µs        100        0       0      1  _partials/helpers/picture.html
        748.41 µs     374.21 µs     402.25 µs          0        0       0      2  list.rss.xml
        718.91 µs     718.91 µs     718.91 µs          0        0       0      1  _shortcodes/figure.html
        696.14 µs     232.04 µs     263.44 µs          0        0       0      3  _shortcodes/newtemplatesystem.html
        676.21 µs     852.00 ns      16.23 µs        100        0       0    793  _partials/layouts/head/speculationrules.html
        618.07 µs     618.07 µs     618.07 µs          0        0       0      1  _shortcodes/chroma-lexers.html
        605.82 µs     763.00 ns      29.52 µs        100        0       0    793  _partials/layouts/hooks/body-start.html
        444.33 µs     222.17 µs     342.78 µs          0        0       0      2  _shortcodes/highlight.html
        398.33 µs     398.33 µs     398.33 µs          0        0       0      1  _shortcodes/render-table-of-pages-in-section.html
        243.64 µs     121.82 µs     125.32 µs          0        0       0      2  _shortcodes/youtube.html
        235.46 µs      39.24 µs      66.03 µs          0        0       0      6  _shortcodes/current-go-version.html
        229.98 µs      12.78 µs      39.15 µs          0        0       0     18  _shortcodes/get-page-desc.html
        193.19 µs       8.05 µs      34.85 µs         52        0       0     24  _partials/layouts/date.html
        184.80 µs     184.80 µs     184.80 µs          0        0       0      1  _shortcodes/x.html
        125.61 µs     125.61 µs     125.61 µs          0        0       0      1  _shortcodes/details.html
        106.29 µs     106.29 µs     106.29 µs          0        0       0      1  _shortcodes/hl.html
         81.00 µs      81.00 µs      81.00 µs          0        0       0      1  _shortcodes/syntax-highlighting-styles-light-with-counterpart.html
         47.32 µs      23.66 µs      40.66 µs          0        0       0      2  _shortcodes/param.html
         42.67 µs      42.67 µs      42.67 µs        100        0       0      1  _partials/layouts/home/features.html
         37.00 µs      37.00 µs      37.00 µs          0        0       0      1  _shortcodes/instagram.html
         35.17 µs      35.17 µs      35.17 µs          0        0       0      1  _shortcodes/vimeo.html
          1.95 µs       1.95 µs       1.95 µs          0        0       0      1  home.headers
          1.81 µs     258.00 ns     977.00 ns          0        0       0      7  _shortcodes/module-mounts-note.html
```

From left to right, the columns represent:

cumulative duration
: The cumulative time spent executing the template.

average duration
: The average time spent executing the template.

maximum duration
: The maximum time spent executing the template.

cache potential
: Displayed as a percentage, any _partial_ template with a 100% cache potential should be called with the [`partialCached`][] function instead of the [`partial`][] function. See the [caching](#caching) section below.

  > [!WARNING]
  > A 100% cache potential is calculated by comparing each execution's rendered output to the first execution's output, so it cannot detect side effects such as calls to [`warnf`][] or [`errorf`][], or execution order dependencies such as a conditional based on [`IsHome`][]. A _partial_ template that produces no visible output but performs validation, logging, or conditional logic based on context can show 100% cache potential even though caching it would suppress that behavior for all but the first invocation. Review the template's logic before switching it to `partialCached`.

percent cached
: The number of times the rendered templated was cached divided by the number of times the template was executed.

cached count
: The number of times the rendered templated was cached.

total count
: The number of times the template was executed.

template
: The path to the template, relative to the `layouts` directory.

> [!NOTE]
> Hugo builds pages in parallel where multiple pages are generated simultaneously. Because of this parallelism, the sum of "cumulative duration" values is usually greater than the actual time it takes to build a site.

## Caching

Some _partial_ templates such as sidebars or menus are executed many times during a site build. Depending on the content within the _partial_ template and the desired output, the template may benefit from caching to reduce the number of executions. The [`partialCached`][] function provides caching capabilities for _partial_ templates.

> [!NOTE]
> Note that you can create cached variants of each _partial_ template by passing additional arguments to `partialCached` beyond the initial context. See the `partialCached` documentation for more details.

## Timers

Use the [`debug.Timer`][] function to determine execution time for a block of code, useful for finding performance bottlenecks in templates.

[Microsoft Defender Antivirus documentation]: https://support.microsoft.com/en-us/topic/how-to-add-a-file-type-or-process-exclusion-to-windows-security-e524cbc2-3975-63c2-f9d1-7c2eb5331e53
[`debug.Timer`]: /functions/debug/timer/
[`errorf`]: /functions/fmt/errorf/
[`IsHome`]: /methods/page/ishome/
[`partialCached`]: /functions/partials/includecached/
[`partial`]: /functions/partials/include/
[`warnf`]: /functions/fmt/warnf/
