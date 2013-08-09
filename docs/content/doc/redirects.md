---
title: "Redirects"
Pubdate: "2013-07-09"
---

For people migrating existing published content to Hugo theres a good chance
you need a mechanism to handle redirecting old urls.

Luckily, this can be handled easily in a couple of easy steps.

1. Create a special post for the redirect and mark the file as a `redirect`
    file in the front matter.  Here is an example
    `content/redirects/my-awesome-blog-post.md` :

    ```markdown
    ---
    redirect: true
    slug: /my-awesome-blog-post/
    url: /docs/redirects/
    ---
```

2. Set the redirect template `layouts/redirects/single.html`:

    ```html
    <!DOCTYPE html>
    <html>
    <head>
      <link rel="canonical" href="{{ .Url }}"/>
      <meta http-equiv="content-type" content="text/html; charset=utf-8" />
      <meta http-equiv="refresh" content="0;url={{ .Url }}" />
    </head>
    </html>
    ```

Now when you go to `/my-awesome-blog-post/` it will do a meta redirect to
`/docs/redirects/`.