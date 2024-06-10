---
title: Frequently asked questions
linkTitle: FAQs
description: These questions are frequently asked by new users.
categories: [troubleshooting]
keywords: [faq]
menu:
  docs:
    parent: troubleshooting
    weight: 70
weight: 70
# Use level 6 headings for each question.
---

Hugo’s [forum] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

These are just a few of the questions most frequently asked by new users.

###### An error message indicates that a feature is not available. Why?

Hugo is available in two editions: standard and extended. With the extended edition you can (a) encode to the WebP format when processing images, and (b) transpile Sass to CSS using the embedded LibSass transpiler. The extended edition is not required to use the Dart Sass transpiler.

When you attempt to perform either of the operations above with the standard edition, Hugo throws this error:

```go-html-template
Error: this feature is not available in your current Hugo version
```

To resolve, uninstall the standard edition, then install the extended edition. See the [installation] section for details.

###### Why do I see "Page Not Found" when visiting the home page?

In the content/_index.md file:

  - Is `draft` set to `true`?
  - Is the `date` in the future?
  - Is the `publishDate` in the future?
  - Is the `expiryDate` in the past?

If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

###### Why is a given section not published?

In the content/section/_index.md file:

  - Is `draft` set to `true`?
  - Is the `date` in the future?
  - Is the `publishDate` in the future?
  - Is the `expiryDate` in the past?

If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

###### Why is a given page not published?

In the content/section/page.md file, or in the content/section/page/index.md file:

  - Is `draft` set to `true`?
  - Is the `date` in the future?
  - Is the `publishDate` in the future?
  - Is the `expiryDate` in the past?

If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

###### Why can't I see any of a page's descendants?

You may have an index.md file instead of an _index.md file. See&nbsp;[details](/content-management/page-bundles/).

###### What is the difference between an index.md file and an _index.md file?

A directory with an index.md file is a [leaf bundle]. A directory with an _index.md file is a [branch bundle]. See&nbsp;[details](/content-management/page-bundles/).

[branch bundle]: /getting-started/glossary/#branch-bundle
[leaf bundle]: /getting-started/glossary/#leaf-bundle

###### Why is my partial template not rendered as expected? {#foo}

You may have neglected to pass the required [context] when calling the partial. For example:

```go-html-template
{{/* incorrect */}}
{{ partial "_internal/pagination.html" }}

{{/* correct */}}
{{ partial "_internal/pagination.html" . }}
```

###### In a template, what's the difference between `:=` and `=` when assigning values to variables?

Use `:=` to initialize a variable, and use `=` to assign a value to a variable that has been previously initialized. See&nbsp;[details](https://pkg.go.dev/text/template#hdr-Variables).

###### When I paginate a list page, why is the page collection not filtered as specified?

You are probably invoking the [`Paginate`] or [`Paginator`] method more than once on the same page. See&nbsp;[details](/templates/pagination/#list-paginator-pages).

###### Why are there two ways to call a shortcode?

Use the `{{%/* shortcode */%}}` notation if the shortcode template, or the content between the opening and closing shortcode tags, contains markdown. Otherwise use the\
`{{</* shortcode */>}}` notation. See&nbsp;[details](/content-management/shortcodes/).

###### Can I use environment variables to control configuration?

Yes. See&nbsp;[details](/getting-started/configuration/#configure-with-environment-variables).

###### Why am I seeing inconsistent output from one build to the next?

The most common causes are page collisions (publishing two pages to the same path) and the effects of concurrency. Use the `--printPathWarnings` command line flag to check for page collisions, and create a topic on the [forum] if you suspect concurrency problems.

###### Which page methods trigger content rendering?

The following methods on a `Page` object trigger content rendering: `Content`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount`.

{{% note %}}
For other questions please visit the [forum]. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

[forum]: https://discourse.gohugo.io
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
{{% /note %}}

[`Paginate`]: /methods/page/paginate
[`Paginator`]: /methods/page/paginator
[context]: /getting-started/glossary/#context
[forum]: https://discourse.gohugo.io
[installation]: /installation
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
