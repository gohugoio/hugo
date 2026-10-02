---
title: cascade
reference: /configuration/cascade/
---

To _cascade_ a value is to apply [_front matter_](g) values from a [_branch_](g) page or the project configuration to descendant pages. Hugo does not cascade a value if the descendant already defines the field, or if a closer ancestor [branch](g) or an earlier element in a cascade array has already set a value for the same field. Use a [_page matcher_](g) to limit cascading to a subset of pages.
