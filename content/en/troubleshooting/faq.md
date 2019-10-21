---
title: Frequently Asked Questions
linktitle: FAQ
description: Solutions to some common Hugo problems.
date: 2018-02-10
categories: [troubleshooting]
menu:
  docs:
    parent: "troubleshooting"
keywords: [faqs]
weight: 2
toc: true
aliases: [/faq/]
---

{{% note %}}
**Note:** The answers/solutions presented below are short, and may not be enough to solve your problem. Visit [Hugo Discourse](https://discourse.gohugo.io/) and use the search. It that does not help, start a new topic and ask your questions.
{{% /note %}}

## I can't see my content!

Is your markdown file [in draft mode](https://gohugo.io/content-management/front-matter/#front-matter-variables)? When testing, run `hugo server` with the `-D` or `--buildDrafts` [switch](https://gohugo.io/getting-started/usage/#draft-future-and-expired-content). 

## Can I set configuration variables via OS environment?

Yes you can! See [Configure with Environment Variables](/getting-started/configuration/#configure-with-environment-variables).

## How do I schedule posts?

1. Set `publishDate` in the page [Front Matter](/content-management/front-matter/) to a date in the future.
2. Build and publish at intervals.

How to automate the "publish at intervals" part depends on your situation:

* If you deploy from your own PC/server, you can automate with [Cron](https://en.wikipedia.org/wiki/Cron) or similar.
* If your site is hosted on a service similar to [Netlify](https://www.netlify.com/) you can use a service such as [ifttt](https://ifttt.com/date_and_time) to schedule the updates.

Also see this Twitter thread:

{{< tweet 962380712027590657 >}}

## Can I use the latest Hugo version on Netlify?

Yes you can! Read [this](/hosting-and-deployment/hosting-on-netlify/#configure-hugo-version-in-netlify).

## I get "TOCSS ... this feature is not available in your current Hugo version"

If you process `SCSS` or `SASS` to `CSS` in your Hugo project, you need the Hugo `extended` version, or else you may see this error message:

```bash
error: failed to transform resource: TOCSS: failed to transform "scss/main.scss" (text/x-scss): this feature is not available in your current Hugo version
```

We release two set of binaries for technical reasons. The extended version is not what you get by default for some installation methods. On the [release page](https://github.com/gohugoio/hugo/releases), look for archives with `extended` in the name. To build `hugo-extended`, use `go install --tags extended`

To confirm, run `hugo version` and look for the word `extended`.
