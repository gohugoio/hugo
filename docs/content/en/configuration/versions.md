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

## Settings

Use the following setting to define how Hugo orders versions.

weight
: (`int`) The language [weight](g).

## Sort order

Hugo sorts versions by weight in ascending order, then by their [semantic version][] in descending order. This affects build order and complement selection.

## Example

The following configuration demonstrates how to define multiple versions with specific weights.

{{< code-toggle >}}
[versions."v1.0.0"]
weight = 20
[versions."v2.0.0"]
weight = 10
{{< /code-toggle >}}

[semantic version]: https://semver.org/
