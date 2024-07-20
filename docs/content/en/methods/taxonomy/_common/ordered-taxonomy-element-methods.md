---
# Do not remove front matter.
---

An ordered taxonomy is a slice, where each element is an object that contains the term and a slice of its weighted pages.

Each element of the slice provides these methods:

Count
: (`int`) Returns the number of pages to which the term is assigned.

Page
: (`page.Page`) Returns the term's `Page` object, useful for linking to the term page.

Pages
: (`page.Pages`) Returns a `Pages` object containing the `Page` objects to which the term is assigned, sorted by [taxonomic weight]. To sort or group, use any of the [methods] available to the `Pages` object. For example, sort by the last modification date.

Term
: (`string`) Returns the term name.

WeightedPages
: (`page.WeightedPages`) Returns a slice of weighted pages to which the term is assigned, sorted by [taxonomic weight]. The `Pages` method above is more flexible, allowing you to sort and group.

[methods]: /methods/pages/
[taxonomic weight]: /getting-started/glossary/#taxonomic-weight
