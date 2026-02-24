---
_comment: Do not remove front matter.
---

The returned collection follows a hierarchical sort where each subsequent dimension acts as a tie-breaker for the one above it.

1. [Language](g) is sorted by [weight](g) in ascending order, falling back to lexicographical order if weights are tied or undefined.
1. [Version](g) is then sorted by weight in ascending order, with Hugo defaulting to a descending semantic sort for any ties.
1. [Role](g) is finally sorted by weight in ascending order, using lexicographical order as the final fallback.
