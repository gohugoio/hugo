---
title: Frequently asked questions
linkTitle: FAQs
description: These questions are frequently asked by new users.
categories: []
keywords: []
---

Hugo's [forum] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

These are just a few of the questions most frequently asked by new users.

An error message indicates that a feature is not available. Why?
: <!-- do not remove preceding space -->
  {{% include "/_common/installation/01-editions.md" %}}

  When you attempt to use a feature that is not available in the edition that you installed, Hugo throws this error:

  ```go-html-template
  this feature is not available in this edition of Hugo
  ```

  To resolve, install a different edition based on the feature table above. See the [installation] section for details.

Why do I see "Page Not Found" when visiting the home page?
: In the `content/_index.md` file:

  - Is `draft` set to `true`?
  - Is the `date` in the future?
  - Is the `publishDate` in the future?
  - Is the `expiryDate` in the past?

  If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

Why is a given page not published?
: In the `content/section/page.md` file, or in the `content/section/page/index.md` file:

  - Is `draft` set to `true`?
  - Is the `date` in the future?
  - Is the `publishDate` in the future?
  - Is the `expiryDate` in the past?

  If the answer to any of these questions is yes, either change the field values, or use one of these command line flags: `--buildDrafts`, `--buildFuture`, or `--buildExpired`.

Why can't I see any of a page's descendants?
: You may have an&nbsp;`index.md`&nbsp;file instead of an&nbsp;`_index.md`&nbsp;file. See&nbsp;[details](/content-management/page-bundles/).

What is the difference between an&nbsp;`index.md`&nbsp;file and an&nbsp;`_index.md`&nbsp;file?
: A directory with an `index.md file` is a [leaf bundle](g). A directory with an&nbsp;`_index.md`&nbsp;file is a [branch bundle](g). See&nbsp;[details](/content-management/page-bundles/).

Why is my _partial_ template not rendered as expected?
: You may have neglected to pass the required [context](g) when calling the partial. For example:

  ```go-html-template
  {{/* incorrect */}}
  {{ partial "pagination.html" }}

  {{/* correct */}}
  {{ partial "pagination.html" . }}
  ```

In a template, what's the difference between `:=` and `=` when assigning values to variables?
: Use `:=` to initialize a variable, and use `=` to assign a value to a variable that has been previously initialized. See&nbsp;[details](https://pkg.go.dev/text/template#hdr-Variables).

When I paginate a list page, why is the page collection not filtered as specified?
: You are probably invoking the [`Paginate`] or [`Paginator`] method more than once on the same page. See&nbsp;[details](/templates/pagination/).

Why are there two ways to call a shortcode?
: Use the `{{%/* shortcode */%}}` notation if the _shortcode_ template, or the content between the opening and closing shortcode tags, contains Markdown. Otherwise use the\
`{{</* shortcode */>}}` notation. See&nbsp;[details](/content-management/shortcodes/#notation).

Can I use environment variables to control configuration?
: Yes. See&nbsp;[details](/configuration/introduction/#environment-variables).

Why am I seeing inconsistent output from one build to the next?
: The most common causes are page collisions (publishing two pages to the same path) and the effects of concurrency. Use the `--printPathWarnings` command line flag to check for page collisions, and create a topic on the [forum] if you suspect concurrency problems.

Why isn't Hugo's development server detecting file changes?
: In its default configuration, Hugo's file watcher may not be able detect file changes when:

  - Running Hugo within Windows Subsystem for Linux (WSL/WSL2) with project files on a Windows partition
  - Running Hugo locally with project files on a removable drive
  - Running Hugo locally with project files on a storage server accessed via the NFS, SMB, or CIFS protocols

  In these cases, instead of monitoring native file system events, use the `--poll` command line flag. For example, to poll the project files every 700 milliseconds, use `--poll 700ms`.

Why is my page Store missing a value?
: The [`Store`] method on a `Page` object allows you to create a [scratch pad](g) on the given page to store and manipulate data. Values are often set within a _shortcode_ template, a _partial_ template called by a _shortcode_ template, or by a _render hook_ template. In all three cases, the scratch pad values are not determinate until Hugo renders the page content.

  If you need to access a scratch pad value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop](g) variable:

  ```go-html-template
  {{ $noop := .Content }}
  {{ .Store.Get "mykey" }}
  ```

  You can trigger content rendering with other methods as well. See next FAQ.

Which page methods trigger content rendering?
: The following methods on a `Page` object trigger content rendering: `Content`, `ContentWithoutSummary`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount`.

> [!note]
> For other questions please visit the [forum]. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

[`Paginate`]: /methods/page/paginate/
[`Paginator`]: /methods/page/paginator/
[`Store`]: /methods/page/store
[forum]: https://discourse.gohugo.io
[installation]: /installation/
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
