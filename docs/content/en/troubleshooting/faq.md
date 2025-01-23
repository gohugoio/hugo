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
---

Hugo’s [forum] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

These are just a few of the questions most frequently asked by new users.

###### An error message indicates that a feature is not available. Why? {#feature-not-available}

{{% include "installation/_common/01-editions.md" %}}

When you attempt to use a feature that is not available in the edition that you installed, Hugo throws this error:

```go-html-template
this feature is not available in this edition of Hugo
```

To resolve, install a different edition based on the feature table above. See the [installation] section for details.

###### Why do I see "Page Not Found" when visiting the home page?

In the `content/_index.md` file:

  - Is `draft` set to `true`?
  - Is the `date` in the future?
  - Is the `publishDate` in the future?
  - Is the `expiryDate` in the past?

If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

###### Why is a given page not published?

In the `content/section/page.md` file, or in the `content/section/page/index.md` file:

- Is `draft` set to `true`?
- Is the `date` in the future?
- Is the `publishDate` in the future?
- Is the `expiryDate` in the past?

If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

###### Why can't I see any of a page's descendants?

You may have an `index.md` file instead of an `_index.md` file. See&nbsp;[details](/content-management/page-bundles/).

###### What is the difference between an index.md file and an _index.md file?

A directory with an `index.md file` is a [leaf bundle](g). A directory with an `_index.md` file is a [branch bundle](g). See&nbsp;[details](/content-management/page-bundles/).

###### Why is my partial template not rendered as expected?

You may have neglected to pass the required [context](g) when calling the partial. For example:

```go-html-template
{{/* incorrect */}}
{{ partial "_internal/pagination.html" }}

{{/* correct */}}
{{ partial "_internal/pagination.html" . }}
```

###### In a template, what's the difference between `:=` and `=` when assigning values to variables?

Use `:=` to initialize a variable, and use `=` to assign a value to a variable that has been previously initialized. See&nbsp;[details](https://pkg.go.dev/text/template#hdr-Variables).

###### When I paginate a list page, why is the page collection not filtered as specified?

You are probably invoking the [`Paginate`] or [`Paginator`] method more than once on the same page. See&nbsp;[details](/templates/pagination/).

###### Why are there two ways to call a shortcode?

Use the `{{%/* shortcode */%}}` notation if the shortcode template, or the content between the opening and closing shortcode tags, contains Markdown. Otherwise use the\
`{{</* shortcode */>}}` notation. See&nbsp;[details](/content-management/shortcodes/).

###### Can I use environment variables to control configuration?

Yes. See&nbsp;[details](/getting-started/configuration/#configure-with-environment-variables).

###### Why am I seeing inconsistent output from one build to the next?

The most common causes are page collisions (publishing two pages to the same path) and the effects of concurrency. Use the `--printPathWarnings` command line flag to check for page collisions, and create a topic on the [forum] if you suspect concurrency problems.

###### Why isn't Hugo's development server detecting file changes?

In its default configuration, Hugo's file watcher may not be able detect file changes when:

- Running Hugo within Windows Subsystem for Linux (WSL/WSL2) with project files on a Windows partition
- Running Hugo locally with project files on a removable drive
- Running Hugo locally with project files on a storage server accessed via the NFS, SMB, or CIFS protocols

In these cases, instead of monitoring native file system events, use the `--poll` command line flag. For example, to poll the project files every 700 milliseconds, use `--poll 700ms`.

###### Why is my page Scratch or Store missing a value?

The [`Scratch`] and [`Store`] methods on a `Page` object allow you to create a [scratch pad](g) on the given page to store and manipulate data. Values are often set within a shortcode, a partial template called by a shortcode, or by a Markdown render hook. In all three cases, the scratch pad values are not determinate until Hugo renders the page content.

If you need to access a scratch pad value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop](g) variable:

```go-html-template
{{ $noop := .Content }}
{{ .Store.Get "mykey" }}
```

You can trigger content rendering with other methods as well. See next FAQ.

[`Scratch`]: /methods/page/scratch
[`Store`]: /methods/page/store

###### Which page methods trigger content rendering?

The following methods on a `Page` object trigger content rendering: `Content`, `ContentWithoutSummary`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount`.

{{% note %}}
For other questions please visit the [forum]. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

[forum]: https://discourse.gohugo.io
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
{{% /note %}}

[`Paginate`]: /methods/page/paginate/
[`Paginator`]: /methods/page/paginator/
[forum]: https://discourse.gohugo.io
[installation]: /installation/
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
