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

## Base settings

Configure the following base settings:

{{< code-toggle file=hugo >}}
defaultContentVersion = 'v1.0.0'
defaultContentVersionInSubdir = false
{{< /code-toggle >}}

`defaultContentVersion`
: (`string`) The project's [default version](g). When one or more versions are defined, this value must match one of the defined version names.

`defaultContentVersionInSubdir`
: (`bool`) Whether to publish the default content version to a subdirectory matching the [`defaultContentVersion`](#defaultcontentversion). Default is `false`.

## Version settings

Use the following setting to define how Hugo orders versions.

`weight`
: (`int`) The version [weight](g). When set to a non-zero value, this is the primary sort criteria for this version.

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
