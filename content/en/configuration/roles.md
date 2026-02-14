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

To define "guest" and "member" roles:

{{< code-toggle >}}
[roles.guest]
weight = 20
[roles.member]
weight = 10
{{< /code-toggle >}}

Roles are sorted by their [weight](g) in descending order, then by their name in descending order. This affects build order and complement selection.
