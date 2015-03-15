---
aliases:
- /doc/redirects/
- /doc/alias/
- /doc/aliases/
date: 2013-07-09
menu:
  main:
    parent: extras
next: /extras/builders
prev: /taxonomies/ordering
title: Aliases
weight: 10
---

For people migrating existing published content to Hugo, there's a good chance
you need a mechanism to handle redirecting old URLs.

Luckily, this can be handled easily with aliases in Hugo.

## Example
**content/posts/my-awesome-blog-post.md**

<table class="table">
<thead>
<tr>
<th>TOML</th><th>YAML</th>
</tr>
</thead>
<tbody>
<tr valign="top">
<td><pre><code>+++
aliases = [
    "/posts/my-original-url/",
    "/2010/even-earlier-url.html"
]
+++
</code></pre></td>
<td><pre><code>---
aliases:
    - /posts/my-original-url/
    - /2010/even-earlier-url.html
---
</code></pre></td>
</tr>
</tbody>
</table>

Now when you go to any of the aliases locations, they
will redirect to the page.

## Important Behaviors

1. *Hugo makes no assumptions about aliases. They also don't change based
on your UglyUrls setting. You need to provide absolute path to your webroot and the
complete filename or directory.*

2. *Aliases are rendered prior to any content and will be overwritten by
any content with the same location.*
