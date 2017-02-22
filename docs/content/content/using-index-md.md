---
aliases:
- /doc/using-index-md/
lastmod: 2017-02-22
date: 2017-02-22
linktitle: Using _index.md
menu:
  main:
    parent: content
prev: /content/example
next: /themes/overview
notoc: true
title: Using _index.md
weight: 70
---
# \_index.md and 'Everything is a Page'

As of version v0.18 Hugo now treats '[everything as a page](http://bepsays.com/en/2016/12/19/hugo-018/)'. This allows you to add content and frontmatter to any page - including List pages like [Sections](/content/sections/), [Taxonomies](/taxonomies/overview/), [Taxonomy Terms pages](/templates/terms/) and even to 'special case' pages like the [Home page](/templates/homepage/).

In order to take advantage of this behaviour you need to do a few things. 

1. Create an \_index.md file that contains the frontmatter and content you would like to apply.

2. Place the \_index.md file in the correct place in the directory structure. 

3. Ensure that the respective template is configured to display `{{ .Content }}` if you wish for the content of the \_index.md file to be rendered on the respective page. 

## How \_index.md pages work

Before continuing it's important to know that this page must reference certain templates to describe how the \_index.md page will be rendered. Hugo has a multitude of possible templates that can be used and placed in various places (think theme templates for instance). For simplicity/brevity the default/top level template location will be used to refer to the entire range of places the template can be placed. 

If this is confusing or you are unfamiliar with Hugo's template hierarchy, visit the various template pages listed below. You may need to find the 'active' template responsible for any particular page on your own site by going through the template hierarchy and matching it to your particular setup/theme you are using. 

- [Home page template](/templates/homepage/)
- [Content List templates](/templates/list/)
- [Single Content templates](/templates/content/)
- [Taxonomy Terms templates](/templates/terms/)

Now that you've got a handle on templates lets recap some Hugo basics to understand how to use an \_index.md file with a List page.

1. Sections and Taxonomies are 'List' pages, NOT single pages.
2. List pages are rendered using the template heirarchy found in the [Content - List Template](http://localhost:1313/templates/list/) docs.
3. The Home page is not technically a list page but has [it's own template](/templates/homepage/) at layouts/index.html.
4. Taxonomy terms pages are 'lists of metadata' not lists of content, so [have their own templates](/templates/terms/). 

Let's put this all together clearly because it's a source of much confusion:

>**\_index.md files used in List pages, Terms pages or the Home page are NOT rendered as single pages or with Single Content templates, they are used to 'inject' frontmatter/content into the respective List pages, Terms pages or the Home page and those pages are rendered using templates listed above.**

Here are a couple of examples to make it clearer...

| \_index.md location               | Page affected             | Rendered by                   |
| -------------------               | ------------              | -----------                   |
| /content/post/\_index.md          | site.com/post/            | /layouts/section/post.html    |
| /content/categories/hugo/index.md | site.com/categories/hugo/ | /layouts/taxonomy/hugo.html   |

## How to display content from \_index.md files

From the information above it should follow that content within an \_index.md file won't be rendered in its own Single Page, it'll be 'injected' into the respective List/Terms/Home page. 

To **_actually display that content_** you need to ensure that the relevant template responsible for rendering the List/Terms/Home page contains (at least) `{{ .Content }}`. 

This is the way to actually display the content within the \_index.md file on the List/Terms/Home page. 

A very simple/naive example of this would be:

```html
{{ partial "header.html" . }}
	<main>
        {{ .Content }}
		{{ $paginator := .Paginate (where .Site.Pages "Type" "post") }}
		{{ range $paginator.Pages }}
			{{ partial "summary.html" . }}
		{{ end }}
		{{ partial "pagination.html" . }}
	</main>
{{ partial "sidebar.html" . }}
{{ partial "footer.html" . }}
```

You can see `{{ .Content }}` just after the `<main>` element. For this particular example, the content of the \_index.md file will show before the main list of summaries.

## Where to organise an \_index.md file

To add content and frontmatter to the home page, a section, a taxonomy or a taxonomy terms listing, add a markdown file with the base name \_index on the relevant place on the file system.

```bash
└── content
    ├── _index.md
    ├── categories
    │   ├── _index.md
    │   └── photo
    │       └── _index.md
    ├── post
    │   ├── _index.md
    │   └── firstpost.md
    └── tags
        ├── _index.md
        └── hugo
            └── _index.md
```

In the above example \_index.md pages have been added to each section/taxonomy. 

An \_index.md file has also been added in the top level 'content' directory. 

### Where to place \_index.md for the Home page

Hugo themes are designed to use the 'content' directory as the root of the website, so adding an \_index.md file here (like has been done in the example above) is how you would add frontmatter/content to the home page. 




