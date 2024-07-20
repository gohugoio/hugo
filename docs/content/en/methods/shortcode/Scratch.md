---
title: Scratch
description: Returns a "scratch pad" scoped to the shortcode to store and manipulate data. 
categories: []
keywords: []
action:
  related:
    - functions/collections/NewScratch
  returnType: maps.Scratch
  signatures: [SHORTCODE.Scratch]
---

The `Scratch` method within a shortcode creates a [scratch pad] to store and manipulate data. The scratch pad is scoped to the shortcode, and is reset on server rebuilds.

{{% note %}}
With the introduction of the [`newScratch`] function, and the ability to [assign values to template variables] after initialization, the `Scratch` method within a shortcode is obsolete.

[assign values to template variables]: https://go.dev/doc/go1.11#text/template
[`newScratch`]: /functions/collections/newscratch/
{{% /note %}}

[scratch pad]: /getting-started/glossary/#scratch-pad

{{% include "methods/page/_common/scratch-methods.md" %}}
