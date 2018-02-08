---

# A suitable title for this article.
title: Hugo Showcase Template

# Set this to the current date.
date: 2018-02-07

description: "A short description of this page."

# The URL to the site on the internet.
siteURL: https://gohugo.io/

# Link to the site's Hugo source code if public and you can/want to share.
# Remove or leave blank if not needed/wanted.
siteSource: https://github.com/gohugoio/hugoDocs

# Add credit to the article author. Leave blank or remove if not needed/wanted.
byline: "[bep](https://github.com/bep), Hugo Lead"

---

Have a **notable Hugo site[^1]**? We would love to feature it in this **Showcase Section**

We would really appreciate if you could:

1. Fork https://github.com/gohugoio/hugoDocs
1. Create a copy of the [content/showcase/template](https://github.com/gohugoio/hugoDocs/tree/master/content/showcase/template) directory with a suitable name. If you now run `hugo server`, your site should show up in [http://localhost:1313/showcase/](http://localhost:1313/showcase/) and on the front page.
2. Adjust the [files](#files) and write a story about your site
3. Create a new pull request in https://github.com/gohugoio/hugoDocs/pulls

**Note:** The Showcase section uses the latest bells and whistles from Hugo, [resources](/content-management/page-resources/) with [image processing](/content-management/image-processing/), so you need a reasonable up-to-date [Hugo version](https://github.com/gohugoio/hugo/releases).

## Files

The content of the [content/showcase/template](https://github.com/gohugoio/hugoDocs/tree/master/content/showcase/template) directory explained:

index.md
: The main content file. Fill in required front matter metadata and write your story. I does not have to be a novel. It can even be self-promotional, but it should include Hugo in some form.

bio.md
: A short summary of the website. Site credits (who built it) fits nicely here.

featured-template.png
: A reasonably sized screenshot of your website. It can be named anything, but the name must start with "featured". The sample image is `1500x750` (2:1 aspect ratio).



[^1]: We want this to show Hugo in its best light, so this is not for the average Hugo blog. In most cases the answer to "Is my site [notabable](http://www.dictionary.com/browse/notable)?" will be obvious, but if in doubt, create an [issue](https://github.com/gohugoio/hugoDocs/issues) with a link and some words, and we can discuss it. But if you have a site with an interesting Hugo story or a company site where the company itself is notable, you are most welcome.
