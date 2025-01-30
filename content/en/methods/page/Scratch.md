---
title: Scratch
description: Returns a "scratch pad" on the given page to store and manipulate data.
categories: []
keywords: []
action:
  related: []
  returnType: maps.Scratch
  signatures: [PAGE.Scratch]
expiryDate: 2025-11-18 #  deprecated 2024-11-18
---

{{% deprecated-in 0.138.0 %}}
Use the [`PAGE.Store`] method instead.

This is a soft deprecation. This method will be removed in a future release, but the removal date has not been established. Although Hugo will not emit a warning if you continue to use this method, you should begin using `PAGE.Store` as soon as possible.

Beginning with v0.138.0 the `PAGE.Scratch` method is aliased to `PAGE.Store`.

[`PAGE.Store`]: /methods/page/store/
{{% /deprecated-in %}}
