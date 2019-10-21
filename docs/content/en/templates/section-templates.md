---
title: Section Page Templates
linktitle: Section Templates
description: Templates used for section pages are **lists** and therefore have all the variables and methods available to list pages.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [lists,sections,templates]
menu:
  docs:
    parent: "templates"
    weight: 40
weight: 40
sections_weight: 40
draft: false
aliases: [/templates/sections/]
toc: true
---

## Add Content and Front Matter to Section Templates

To effectively leverage section page templates, you should first understand Hugo's [content organization](/content-management/organization/) and, specifically, the purpose of `_index.md` for adding content and front matter to section and other list pages.

## Section Template Lookup Order

See [Template Lookup](/templates/lookup-order/).

## Page Kinds

Every `Page` in Hugo has a `.Kind` attribute.

| Kind           | Description                                                        | Example                                                                       |
|----------------|--------------------------------------------------------------------|-------------------------------------------------------------------------------|
| `home`         | The home page                                                      | `/index.html`                                                                 |
| `page`         | A page showing a _regular page_                                    | `my-post` page (`/posts/my-post/index.html`)                                  |
| `section`      | A page listing _regular pages_ from a given [_section_][sections]  | `posts` section (`/posts/index.html`)                                         |
| `taxonomy`     | A page listing _regular pages_ from a given _taxonomy term_        | page for the term `awesome` from `tags` taxonomy (`/tags/awesome/index.html`) |
| `taxonomyTerm` | A page listing terms from a given _taxonomy_                       | page for the `tags` taxonomy (`/tags/index.html`)                             |

## `.Site.GetPage` with Sections

`Kind` can easily be combined with the [`where` function][where] in your templates to create kind-specific lists of content. This method is ideal for creating lists, but there are times where you may want to fetch just the index page of a single section via the section's path.

The [`.GetPage` function][getpage] looks up an index page of a given `Kind` and `path`.

You can call `.Site.GetPage` with two arguments: `kind` (one of the valid values
of `Kind` from above) and `kind value`.

Examples:

- `{{ .Site.GetPage "section" "posts" }}`
- `{{ .Site.GetPage "page" "search" }}`

## Example: Creating a Default Section Template

{{< code file="layouts/_default/section.html" download="section.html" >}}
{{ define "main" }}
  <main>
      {{ .Content }}
          <ul class="contents">
          {{ range .Paginator.Pages }}
              <li>{{.Title}}
                  <div>
                    {{ partial "summary.html" . }}
                  </div>
              </li>
          {{ end }}
          </ul>
      {{ partial "pagination.html" . }}
  </main>
{{ end }}
{{< /code >}}

### Example: Using `.Site.GetPage`

The `.Site.GetPage` example that follows assumes the following project directory structure:

```
.
└── content
    ├── blog
    │   ├── _index.md # "title: My Hugo Blog" in the front matter
    │   ├── post-1.md
    │   ├── post-2.md
    │   └── post-3.md
    └── events #Note there is no _index.md file in "events"
        ├── event-1.md
        └── event-2.md
```

`.Site.GetPage` will return `nil` if no `_index.md` page is found. Therefore, if `content/blog/_index.md` does not exist, the template will output the section name:

```
<h1>{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}</h1>
```

Since `blog` has a section index page with front matter at `content/blog/_index.md`, the above code will return the following result:

```
<h1>My Hugo Blog</h1>
```

If we try the same code with the `events` section, however, Hugo will default to the section title because there is no `content/events/_index.md` from which to pull content and front matter:

```
<h1>{{ with .Site.GetPage "section" "events" }}{{ .Title }}{{ end }}</h1>
```

Which then returns the following:

```
<h1>Events</h1>
```


[contentorg]: /content-management/organization/
[getpage]: /functions/getpage/
[lists]: /templates/lists/
[lookup]: /templates/lookup-order/
[where]: /functions/where/
[sections]: /content-management/sections/
