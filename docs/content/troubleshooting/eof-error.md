---
title: EOF Error
linktitle: EOF Error
description: If you find yourself seeing an EOF error in the console whenever you create a new content file from Hugo's archetype feature.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [troubleshooting]
menu:
  docs:
    parent: "troubleshooting"
keywords: [eof, end of file, error, faqs]
draft: false
weight:
aliases: [/troubleshooting/strange-eof-error/]
toc: true
---

## Trouble: `hugo new` Aborts with EOF error

> I'm running into an issue where I cannot get archetypes working, when running `hugo new showcase/test.md`, for example, I see an `EOF` error thrown by Hugo.
>
> When I run Hugo with v0.12 via `hugo new -v showcase/test.md`, I see the following output:
>
> ```
> INFO: 2015/01/04 Using config file: /private/tmp/test/config.toml
> INFO: 2015/01/04 attempting to create  showcase/test.md of showcase
> INFO: 2015/01/04 curpath: /private/tmp/test/archetypes/showcase.md
> ERROR: 2015/01/04 EOF
> ```
>
> Is there something that I am blatantly missing?

## Solution: Carriage Returns

The solution is to add a final newline (i.e., `EOL`) to the end of your default.md archetype file of your theme. You can do this by adding a carriage return after the closing `+++` or `---` of your TOML or YAML front matter, respectively.

{{% note "Final EOL Unnecessary in v0.13+" %}}
As of v0.13, Hugo's parser has been enhanced to accommodate archetype files without final EOL thanks to the great work by [@tatsushid](https://github.com/tatsushid).
{{% /note %}}

## Discussion Forum References

* http://discourse.gohugo.io/t/archetypes-not-properly-working-in-0-12/544
* http://discourse.gohugo.io/t/eol-f-in-archetype-files/554

## Related Hugo Issues

* https://github.com/gohugoio/hugo/issues/776
