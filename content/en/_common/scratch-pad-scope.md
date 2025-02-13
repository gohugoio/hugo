---
_comment: Do not remove front matter.
---

## Scope

The method or function used to create a scratch pad determines its scope. For example, use the `Store` method on a `Page` object to create a scratch pad scoped to the page.

Scope|Method or function
:--|:--
page|[`PAGE.Store`]
site|[`SITE.Store`] 
global|[`hugo.Store`]
local|[`collections.NewScratch`]
shortcode|[`SHORTCODE.Store`]

[`page.store`]: /methods/page/store
[`site.store`]: /methods/site/store
[`hugo.store`]: /functions/hugo/store
[`collections.newscratch`]: functions/collections/newscratch
[`shortcode.store`]: /methods/shortcode/store
