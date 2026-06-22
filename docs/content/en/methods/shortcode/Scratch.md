---
title: Scratch
description: Returns a persistent data structure for storing and manipulating keyed values, scoped to the current shortcode.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: maps.Scratch
    signatures: [SHORTCODE.Scratch]
expiryDate: 2026-11-18 # deprecated 2024-11-18 (soft)
---

{{< deprecated-in 0.139.0 >}}
Use the [`SHORTCODE.Store`](/methods/shortcode/store/) method instead.

This is a soft deprecation. This method will be removed in a future release, but the removal date has not been established. Although Hugo will not emit a warning if you continue to use this method, you should begin using `SHORTCODE.Store` as soon as possible.

Beginning with v0.139.0 the `SHORTCODE.Scratch` method is aliased to `SHORTCODE.Store`.
{{< /deprecated-in >}}
