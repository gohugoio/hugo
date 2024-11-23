---
title: Store
description: Returns a "Store pad" scoped to the shortcode to store and manipulate data. 
categories: []
keywords: []
action:
  related:
    - functions/collections/NewScratch
    - methods/page/Store
    - methods/site/Store
    - functions/hugo/Store
  returnType: maps.Store
  signatures: [SHORTCODE.Store]
---

{{< new-in 0.139.0 >}}

The `Store` method within a shortcode creates a [scratch pad] to store and manipulate data. The scratch pad is scoped to the shortcode.

{{% note %}}
With the introduction of the [`newScratch`] function, and the ability to [assign values to template variables] after initialization, the `Store` method within a shortcode is mostly obsolete.

[assign values to template variables]: https://go.dev/doc/go1.11#text/template
[`newScratch`]: /functions/collections/newScratch/
{{% /note %}}

[Store pad]: /getting-started/glossary/#scratch-pad

{{% include "methods/page/_common/scratch-methods.md" %}}
