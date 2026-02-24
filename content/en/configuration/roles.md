---
title: Configure roles
linkTitle: Roles
description: Configure roles.
categories: []
keywords: []
---

{{< new-in 0.153.0 />}}

This is the default configuration:

{{< code-toggle config=roles />}}

## Settings

Use the following setting to define how Hugo orders roles.

weight
: (`int`) The role [weight](g).

## Sort order

Hugo sorts roles by weight in ascending order, then lexicographically in ascending order. This affects build order and complement selection.

## Example

The following configuration demonstrates how to define multiple roles with specific weights.

{{< code-toggle >}}
[roles.guest]
weight = 20
[roles.member]
weight = 10
{{< /code-toggle >}}
