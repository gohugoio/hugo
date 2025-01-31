---
title: Scratch
description: Returns a "scratch pad" to store and manipulate data, scoped to the current shortcode.
categories: []
keywords: []
action:
  related: []
  returnType: maps.Scratch
  signatures: [SHORTCODE.Scratch]
expiryDate: 2025-11-18 #  deprecated 2024-11-18
---

{{% deprecated-in 0.139.0 %}}
Use the [`SHORTCODE.Store`] method instead.

This is a soft deprecation. This method will be removed in a future release, but the removal date has not been established. Although Hugo will not emit a warning if you continue to use this method, you should begin using `SHORTCODE.Store` as soon as possible.

Beginning with v0.139.0 the `SHORTCODE.Scratch` method is aliased to `SHORTCODE.Store`.

[`SHORTCODE.Store`]: /methods/shortcode/store/
{{% /deprecated-in %}}
