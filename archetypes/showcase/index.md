---

title: {{ replace .Name "-" " " | title }}
date: {{ now.Format "2006-01-02" }}

description: "A short description of this page."

# The URL to the site on the internet.
siteURL: https://gohugo.io/

# Link to the site's Hugo source code if public and you can/want to share.
# Remove or leave blank if not needed/wanted.
siteSource: https://github.com/gohugoio/hugoDocs

# Add credit to the article author. Leave blank or remove if not needed/wanted.
byline: "[bep](https://github.com/bep), Hugo Lead"

---

To complete this showcase:

1. Write the story about your site in this file.
2. Add a summary to the `bio.md` file in this folder.
3. Replace the `featured-template.png` with a screenshot of your site. You can rename it, but it must contain the word `featured`.
4. Create a new pull request in https://github.com/gohugoio/hugoDocs/pulls

The content of this bundle explained:

index.md
: The main content file. Fill in required front matter metadata and write your story. I does not have to be a novel. It can even be self-promotional, but it should include Hugo in some form.

bio.md
: A short summary of the website. Site credits (who built it) fits nicely here.

featured.png
: A reasonably sized screenshot of your website. It can be named anything, but the name must start with "featured". The sample image is `1500x750` (2:1 aspect ratio).

