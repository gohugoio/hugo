---
title: Configure versions
linkTitle: Versions
description: Configure versions.
categories: []
keywords: []
---

{{< new-in 0.153.0 />}}

This is the default configuration:

{{< code-toggle config=versions />}}

To define "v1.0.0" and "v2.0.0" versions:

{{< code-toggle >}}
[versions."v1.0.0"]
weight = 0
[versions."v2.0.0"]
weight = 0
{{< /code-toggle >}}

Versions are sorted by their [weight](g) in ascending order, then by their [semantic version] in descending order. This affects build order and complement selection.

[semantic version]: https://semver.org/
