---
title: Scratch
description: Creates a "scratch pad" on the given page to store and manipulate data.
categories: []
keywords: []
action:
  related:
    - methods/page/Store
    - functions/collections/NewScratch
  returnType: maps.Scratch
  signatures: [PAGE.Scratch]
aliases: [/extras/scratch/,/doc/scratch/,/functions/scratch]
---

The `Scratch` method on a `Page` object creates a [scratch pad] to store and manipulate data. To create a scratch pad that is not reset on server rebuilds, use the [`Store`] method instead.

To create a locally scoped scratch pad that is not attached to a `Page` object, use the [`newScratch`] function.

[`Store`]: /methods/page/store/
[`newScratch`]: /functions/collections/newscratch/
[scratch pad]: /getting-started/glossary/#scratch-pad

{{% include "methods/page/_common/scratch-methods.md" %}}
